package sockjs

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type sessionXhrPolling struct { 
	*pool // owned by Handler
	sync.Mutex
	closed bool
	in, out *queue
}

func (s *sessionXhrPolling) Receive() ([]byte, error) {
	return s.in.Pull()
}

func (s *sessionXhrPolling) Send(m []byte) error {
	s.out.Push(m)
	return nil
}

func (s *sessionXhrPolling) Close() error {
	s.Lock()
	s.closed = true	
	s.Unlock()
	s.in.Close()
	s.out.Close()
	return nil
}

func handleXhrPolling(h *Handler,w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "application/javascript; charset=UTF-8")
	disableCache(header)
	preflight(header, r)

	sessionFactory := func() Session {
		s := new(sessionXhrPolling)
		s.pool = h.pool
		s.in = newQueue(false)
		s.out = newQueue(true)
		return Session(s)
	}

	s, exists := h.pool.GetOrCreate(sessid, sessionFactory)
	if !exists {
		// initiate connection
		_ ,err := w.Write([]byte("o\n"))
		if err != nil {
			h.pool.Remove(sessid)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go h.hfunc(s)
		return
	}

	x, ok := s.(*sessionXhrPolling)
	if !ok {
		http.NotFound(w, r)
		return
	}

	x.Lock()
	closed := x.closed
	x.Unlock()
	if closed {
		w.Write([]byte("c[3000,\"Go away!\"]\n"))
		return
	}

	m, err := x.out.PullAll()
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

func handleXhrSend(h *Handler,w http.ResponseWriter, r *http.Request, sessid string) {
	s := h.pool.Get(sessid)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	x, ok := s.(*sessionXhrPolling)
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
		x.in.Push([]byte(v))
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
