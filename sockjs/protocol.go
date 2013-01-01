package sockjs

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httputil"
)

type protocol interface {
	contentType() string
	writeOpen(io.Writer) error
	writeData(io.Writer, ...[]byte) (int, error)
	writeClose(io.Writer, int, string)
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

func sessionHeader(w http.ResponseWriter, r *http.Request, p protocol) {
	header := w.Header()
	header.Add("Content-Type", p.contentType())
	noCache(header)
	xhrCors(header, r)
}

func protocolHandler(h *Handler,
	rw http.ResponseWriter,
	r *http.Request,
	sessid string,
	p protocol) {
	var err error
	var w io.Writer
	pw := p.streaming()

	header := rw.Header()
	header.Add("Content-Type", p.contentType())
	sid(h, rw, r)
	noCache(header)

	if pw != nil {
		rw.WriteHeader(http.StatusOK)

		conn, bufrw, err := rw.(http.Hijacker).Hijack()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		wc := newStreamWriter(bufrw)
		defer wc.Close()
		w = wc

		if err = pw.writePrelude(w); err != nil {
			return
		}
	} else {
		w = rw
	}

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = p.writeOpen(w); err != nil {
			h.pool.remove(sessid)
			return
		}
		s.init(r, h.prefix, p, h.config.Headers)
		go h.hfunc(s)
		if pw == nil {
			return
		}
	}

	if h.config.VerifyAddr && !s.verifyAddr(r.RemoteAddr) {
		p.writeClose(w, 2500, "Remote address mismatch")
		return
	}

	if s.interrupted() {
		p.writeClose(w, 1002, "Connection interrupted")
		return
	}

	fail := s.reserve()
	if fail {
		s.interrupt()
		//s.Close()
		p.writeClose(w, 2010, "Another connection still open")
		return
	}
	defer s.free()

	if pw == nil {
		m, err := s.out.pullAll()
		if err != nil {
			p.writeClose(w, 3000, "Go away!")
			return
		}
		p.writeData(w, m...)
	} else {
		for sent := 0; sent < h.config.ResponseLimit; {
			m, err := s.out.pullAll()
			if err != nil {
				p.writeClose(w, 3000, "Go away!")
				return
			}

			n, err := p.writeData(w, m...)
			if err != nil {
				return
			}
			sent += n
		}
	}
}
