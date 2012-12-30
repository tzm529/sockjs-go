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
}

type streamingProtocol interface {
	protocol
	writePrelude(io.Writer) error
}

func pollingHandler(h *Handler,
	w http.ResponseWriter,
	r *http.Request,
	sessid string,
	p protocol) {
	var err error
	header := w.Header()
	header.Add("Content-Type", p.contentType())
	disableCache(header)
	preflight(header, r)

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = p.writeOpen(w); err != nil {
			h.pool.remove(sessid)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.init(r, h.prefix, p, h.config.Headers)
		go h.hfunc(s)
		return
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
		s.Close()
		p.writeClose(w, 2010, "Another connection still open")
		return
	}
	defer s.free()

	m, err := s.out.pullAll()
	if err != nil {
		p.writeClose(w, 3000, "Go away!")
		return
	}
	p.writeData(w, m...)
}

// little helper to simplify streaming handler code
type streamWriter struct {
	bufrw *bufio.ReadWriter
	wc io.WriteCloser
}

func newStreamWriter(bufrw *bufio.ReadWriter) (sw *streamWriter) {
	sw = new(streamWriter)
	sw.bufrw = bufrw
	sw.wc = httputil.NewChunkedWriter(bufrw)
	return
}

func (sw *streamWriter) Write(b []byte) (n int, err error) {
	n, err = sw.wc.Write(b)
	if err != nil { return 0, err }
	err = sw.bufrw.Flush()
	if err != nil { return 0, err }
	return
}

func (sw *streamWriter) Close() error {
	sw.wc.Close()
	sw.bufrw.Write([]byte("\r\n")) // close for chunked data
	return sw.bufrw.Flush()
}

func streamingHandler(h *Handler,
	w http.ResponseWriter,
	r *http.Request,
	sessid string,
	p streamingProtocol) {
	header := w.Header()
	header.Add("Content-Type", p.contentType())
	disableCache(header)
	preflight(header, r)
	w.WriteHeader(http.StatusOK)

	conn, bufrw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	sw := newStreamWriter(bufrw)
	defer sw.Close()

	if err = p.writePrelude(sw); err != nil {
		return
	}

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = p.writeOpen(sw); err != nil {
			h.pool.remove(sessid)
			return
		}

		s.init(r, h.prefix, p, h.config.Headers)
		go h.hfunc(s)
	}

	if h.config.VerifyAddr && !s.verifyAddr(r.RemoteAddr) {
		p.writeClose(sw, 2500, "Remote address mismatch")
		return
	}

	if s.interrupted() {
		p.writeClose(sw, 1002, "Connection interrupted")
		return
	}

	fail := s.reserve()
	if fail {
		s.interrupt()
		s.Close()
		p.writeClose(sw, 2010, "Another connection still open")
		return
	}
	defer s.free()

	for sent := 0; sent < h.config.ResponseLimit; {
		m, err := s.out.pullAll()
		if err != nil {
			p.writeClose(sw, 3000, "Go away!")
			return
		}

		n, err := p.writeData(sw, m...)
		if err != nil {
			return
		}
		sent += n
	}
}
