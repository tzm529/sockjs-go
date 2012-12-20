package sockjs

import (
	"encoding/json"
	"io"
	"net/http"
)

func handleXhrSend(h *Handler,w http.ResponseWriter, r *http.Request, sessid string) {
	s := h.pool.get(sessid)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	x, ok := s.(*pollingBaseSession)
	if !ok {
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
		x.in.push([]byte(v))
	}
	header := w.Header()
	header.Add("Content-Type", "text/plain; charset=UTF-8")
	preflight(header, r)
	disableCache(header)
	w.WriteHeader(http.StatusNoContent)
}

func handleXhrOptions(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Access-Control-Allow-Methods", "OPTIONS, POST")
	preflight(h, r)
	enableCache(h)
	w.WriteHeader(http.StatusNoContent)
}
