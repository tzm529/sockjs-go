package sockjs

import (
	"net/http"
)

func handleXhrPolling(h *Handler,w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "application/javascript; charset=UTF-8")
	disableCache(header)
	preflight(header, r)

	sessionFactory := func() Session {
		s := newPollingBaseSession(h.pool)
		s.out = newQueue(true)
		return Session(s)
	}

	si, exists := h.pool.getOrCreate(sessid, sessionFactory)
	if !exists {
		// initiate connection
		_ ,err := w.Write([]byte("o\n"))
		if err != nil {
			h.pool.remove(sessid)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go h.hfunc(si)
		return
	}

	s, ok := si.(*pollingBaseSession)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if s.closed() {
		w.Write([]byte("c[3000,\"Go away!\"]\n"))
		return
	}

	m, err := s.out.pullAll()
	if err != nil {
		if err == errQueueWait {
			w.Write([]byte("c[2010,\"Another connection still open\"]\n"))
		}
		return
	}
	if m == nil {
		http.NotFound(w, r)
		return
	}
	w.Write(frame("a", "\n", m...))
}
