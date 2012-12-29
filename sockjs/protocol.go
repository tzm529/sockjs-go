package sockjs

import (
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

	fail := s.reserve()
	if fail {
		p.writeClose(w, 2010, "Another connection still open")
		return
	}
	defer s.free()

	if h.config.VerifyAddr && !s.verifyAddr(r.RemoteAddr) {
		p.writeClose(w, 2500, "Remote address mismatch")
		return
	}

	m, err := s.out.pullAll()
	if err != nil {
		p.writeClose(w, 3000, "Go away!")
		return
	}
	p.writeData(w, m...)
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

	chunkedw := httputil.NewChunkedWriter(bufrw)
	defer func() {
		chunkedw.Close()
		bufrw.Write([]byte("\r\n")) // close for chunked data
		bufrw.Flush()
	}()

	if err = p.writePrelude(chunkedw); err != nil {
		return
	}
	if err = bufrw.Flush(); err != nil {
		return
	}

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = p.writeOpen(chunkedw); err != nil {
			goto fail
		}
		if err = bufrw.Flush(); err != nil {
			goto fail
		}
		goto success
	fail:
		h.pool.remove(sessid)
		return
	success:
		s.init(r, h.prefix, p, h.config.Headers)
		go h.hfunc(s)
	}

	fail := s.reserve()
	if fail {
		p.writeClose(chunkedw, 2010, "Another connection still open")
		bufrw.Flush()
		return
	}
	defer s.free()

	if h.config.VerifyAddr && !s.verifyAddr(r.RemoteAddr) {
		p.writeClose(w, 2500, "Remote address mismatch")
		return
	}

	for sent := 0; sent < h.config.ResponseLimit; {
		m, err := s.out.pullAll()
		if err != nil {
			p.writeClose(chunkedw, 3000, "Go away!")
			bufrw.Flush()
			return
		}

		n, err := p.writeData(chunkedw, m...)
		if err != nil {
			return
		}
		if err = bufrw.Flush(); err != nil {
			return
		}
		sent += n
	}
}
