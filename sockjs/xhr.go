package sockjs

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
)

func xhrSendHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "text/plain; charset=UTF-8")
	preflight(header, r)
	disableCache(header)

	s := h.pool.get(sessid)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	var messages []string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&messages); err != nil {
		if err == io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Payload expected."))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Broken JSON encoding."))
		return
	}
	for _, v := range messages {
		s.in.push([]byte(v))
	}
	w.WriteHeader(http.StatusNoContent)
}

func xhrOptionsHandler(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Access-Control-Allow-Methods", "OPTIONS, POST")
	preflight(h, r)
	enableCache(h)
	w.WriteHeader(http.StatusNoContent)
}

type protocol interface{
	contentType() string
	writeOpen(io.Writer) error
	writeData(io.Writer, ...[]byte) (int, error)
	writeClose(io.Writer, int, string)
}

type streamingProtocol interface {
	protocol
	writePrelude(io.Writer) error
}

func baseHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string, proto protocol) {
	var err error
	header := w.Header()
	header.Add("Content-Type", proto.contentType())
	disableCache(header)
	preflight(header, r)

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = proto.writeOpen(w); err != nil {
			h.pool.remove(sessid)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go h.hfunc(s)
		return
	}

	fail := s.reserve()
	if fail {
		proto.writeClose(w, 2010, "Another connection still open")
		return
	}
	defer s.free()

	m, err := s.out.pullAll()
	if err != nil {
		proto.writeClose(w, 3000, "Go away!")
		return
	}
	proto.writeData(w, m...)
}


func baseStreamingHandler(h *Handler, 
	w http.ResponseWriter, 
	r *http.Request, 
	sessid string, 
	proto streamingProtocol) {
	header := w.Header()
	header.Add("Content-Type", proto.contentType())
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
	
	if _, err = bufrw.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")); err != nil { return }
	if err = bufrw.Flush(); err != nil { return }
	if err = proto.writePrelude(chunkedw); err != nil { return }
	if err = bufrw.Flush(); err != nil { return }

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		if err = proto.writeOpen(chunkedw); err != nil { goto fail }
		if err = bufrw.Flush(); err != nil { goto success }
		goto success
	fail:
		h.pool.remove(sessid)
		return
	success:
		go h.hfunc(s)	
	}

	fail := s.reserve()
	if fail {
		proto.writeClose(chunkedw, 2010, "Another connection still open")
		bufrw.Flush()
		return
	}
	defer s.free()

	for sent := 0; sent < h.config.ResponseLimit; {
		m, err := s.out.pullAll()
		if err != nil {
			proto.writeClose(chunkedw, 3000, "Go away!")
			bufrw.Flush()
			return
		}
		
		n, err := proto.writeData(chunkedw, m...)
		if err != nil { return }
		if err = bufrw.Flush(); err != nil { return }
		sent += n
	}
}
