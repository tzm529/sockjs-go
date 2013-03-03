package sockjs

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httputil"
)

type protocol interface {
	contentType() string
	write(io.Writer, []byte) (int, error)
	protocol() Protocol

	// returns a preludeWriter or nil, if the protocol is not streaming.
	streaming() preludeWriter
}

type preludeWriter interface {
	writePrelude(io.Writer) error
}

//* helpers to reduce duplicate code in polling and streaming handlers

type streamWriter struct {
	bufrw *bufio.ReadWriter
	wc    io.WriteCloser
}

func newStreamWriter(bufrw *bufio.ReadWriter) io.WriteCloser {
	sw := new(streamWriter)
	sw.bufrw = bufrw
	sw.wc = httputil.NewChunkedWriter(bufrw)
	return sw
}

func (sw *streamWriter) Write(b []byte) (n int, err error) {
	n, err = sw.wc.Write(b)
	if err != nil {
		return 0, err
	}
	err = sw.bufrw.Flush()
	if err != nil {
		return 0, err
	}
	return
}

func (sw *streamWriter) Close() error {
	sw.wc.Close()
	sw.bufrw.Write([]byte("\r\n")) // close for chunked data
	return sw.bufrw.Flush()
}

func legacyHandler(h *handler,
	rw http.ResponseWriter,
	r *http.Request,
	sessid string,
	p protocol) {
	var err error
	var w io.Writer
	var ok bool
	var infoset bool
	pw := p.streaming()

	header := rw.Header()
	header.Add("Content-Type", p.contentType())
	sid(h, rw, r)
	noCache(header)

	if pw != nil {
		rw.WriteHeader(http.StatusOK)

		// use chunked format for http/1.1.
		// http/1.0 does not support it.
		if r.ProtoMinor == 1 {
			conn, bufrw, err := rw.(http.Hijacker).Hijack()
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			wc := newStreamWriter(bufrw)
			defer wc.Close()
			w = wc
		} else {
			w = rw
		}

		if err = pw.writePrelude(w); err != nil {
			return
		}
	} else {
		w = rw
	}

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if _, err = p.write(w, []byte{'o'}); err != nil {
			h.pool.remove(sessid)
			goto disconnect
		}
		s.init(h.config, p, sessid, h.pool)
		s.setInfo(newRequestInfo(r, h.prefix, h.config.Headers))
		infoset = true
		go h.hfunc(s)
		if pw == nil {
			return
		}
	}

	if h.config.VerifyAddr && !s.verifyAddr(r.RemoteAddr) {
		// not sure what a proper code should be here
		p.write(w, cframe(3001, "Remote address mismatch"))
		logPrintf(h.config.Logger, "%s: request with remote address mismatch from: %s\n",
			s, r.RemoteAddr)
		return
	}

	if s.closed() {
		p.write(w, s.closeFrame())
		return
	}

	ok = s.reserve()
	if !ok {
		// one time close message
		p.write(w, cframe(2010, "Another connection still open"))

		s.close(1002, "Connection interrupted")
		return
	}
	defer s.free()

	if !infoset {
		s.setInfo(newRequestInfo(r, h.prefix, h.config.Headers))
	}
	s.updateRecvStamp()

	if pw == nil {
		//* polling
		m, ok := <-s.sendFrame
		if !ok {
			goto disconnect
		}
		_, err = p.write(w, m)
		if err != nil {
			logPrintf(h.config.Logger, "%s: send error: %s\n", s, err)
			goto disconnect
		}
	} else {
		//* streaming
		var n int
		for sent := 0; sent < h.config.ResponseLimit; {
			m, ok := <-s.sendFrame
			if !ok {
				goto disconnect
			}
			n, err = p.write(w, m)
			if err != nil {
				logPrintf(h.config.Logger, "%s: send error: %s\n", s, err)
				goto disconnect
			}
			sent += n
		}
	}
	return

disconnect:
	// close the session in case it hasn't been closed already
	s.end()
	p.write(w, s.closeFrame())
}
