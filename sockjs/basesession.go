package sockjs

import (
	"sync"
)

type baseSession struct {
	mu      sync.Mutex
	pool    *pool // owned by Server
	closed_ bool
	in      *queue
}

func newBaseSession(pool *pool) (s *baseSession) {
	s = new(baseSession)
	s.pool = pool
	s.in = newQueue(false)
	return
}

func (s *baseSession) closeBase() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed_ = true
	s.pool.close()
}

func (s *baseSession) closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed_
}
