package sockjs

import (
	"net/http"
	"sync"
)

type pollingSession struct {
	mu sync.Mutex
	pool *pool
	closed_ bool
	in_  *queue
	out_ *queue
}

func pollingSessionFactory(pool *pool) session {
	s := new(pollingSession)
	s.pool = pool
	s.in_ = newQueue(false)
	s.out_ = newQueue(true)
	return s
}

func (s *pollingSession) Receive() ([]byte, error) {
	return s.in_.pull()
}

func (s *pollingSession) Send(m []byte) error {
	s.out_.push(m)
	return nil
}

func (s *pollingSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed_ = true
	s.in_.close()
	s.out_.close()
	return nil
}

func (s *pollingSession) in() *queue { return s.in_ }
func (s *pollingSession) out() *queue { return s.out_ }

func (s *pollingSession) closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed_
}

func handleXhrPolling(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "application/javascript; charset=UTF-8")
	disableCache(header)
	preflight(header, r)

	s, exists := h.pool.getOrCreate(sessid, pollingSessionFactory)
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

	m, err := s.out().pullAll()
	if err != nil {
		if err == errQueueWait {
			w.Write(cframe("\n", 2010, "Another connection still open"))
		} else {
			w.Write(cframe("\n", 3000, "Go away!"))
		}
		return
	}
	w.Write(aframe("\n", m...))
}
