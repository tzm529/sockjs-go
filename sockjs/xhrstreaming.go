package sockjs

import (
	"net/http"
	"net/http/httputil"
)

var prelude []byte = make([]byte, 2049)

func init() {
	for i := 0; i < 2048; i++ {
		prelude[i] = byte('h')
	}
	prelude[2048] = byte('\n')
}

func handleXhrStreaming(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "application/javascript; charset=UTF-8")
	header.Add("Transfer-Encoding", "chunked")
	disableCache(header)
	preflight(header, r)

	conn, bufrw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	chunkedw := httputil.NewChunkedWriter(bufrw)
	defer func() {
		chunkedw.Close()
		bufrw.Write([]byte("\r\n")) // close for chunked data
		bufrw.Flush()
	}()
		
	bufrw.Write([]byte("HTTP/1.1 200 OK\n"))
	header.Write(bufrw)
	bufrw.Write([]byte("\n"))
	chunkedw.Write(prelude)
	bufrw.Flush()

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		_, err := chunkedw.Write([]byte("o\n"))
		if err != nil { goto fail }
		err = bufrw.Flush()
		if err == nil { goto ok }
	fail: 
		h.pool.remove(sessid)
		return
	ok:
		go h.hfunc(s)
	}
	defer h.pool.remove(sessid)

	fail := s.reserve()
	if fail {
		chunkedw.Write(cframe("\n", 2010, "Another connection still open"))
		bufrw.Flush()
		return
	}
	defer s.free()

	for sent := 0; sent < h.config.ResponseLimit; {
		m, err := s.out.pullAll()
		if err != nil {
			chunkedw.Write(cframe("\n", 3000, "Go away!"))
			bufrw.Flush()
			return
		}
		n, err := chunkedw.Write(aframe("\n", m...))
		if err != nil { return }
		err = bufrw.Flush()
		if err != nil { return }
		sent += n
	}
}
