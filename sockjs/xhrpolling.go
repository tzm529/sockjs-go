package sockjs

import (
	"net/http"
)

func handleXhrPolling(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "application/javascript; charset=UTF-8")
	disableCache(header)
	preflight(header, r)

	s, exists := h.pool.getOrCreate(sessid)
	if !exists {
		// initiate connection
		_, err := w.Write([]byte("o\n"))
		if err != nil {
			h.pool.remove(sessid)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go h.hfunc(s)
		return
	}

	fail := s.reserve()
	if fail {
		w.Write(cframe("\n", 2010, "Another connection still open"))
		return
	}
	defer s.free()

	m, err := s.out.pullAll()
	if err != nil {
		w.Write(cframe("\n", 3000, "Go away!"))
		return
	}
	w.Write(aframe("\n", m...))
}
