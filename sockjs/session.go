package sockjs

import (
	"errors"
	"sync"
)

var ErrSessionClosed error = errors.New("session closed")
var ErrSessionTimeout error = errors.New("session timeout")

type Session interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
}

// structure for non-websocket sessions
type session struct {
	mu       sync.Mutex
	proto    protocol
	pool     *pool // owned by Server
	in       *queue
	out      *queue
	closed_  bool
	reserved bool
}

func newSession(pool *pool) *session {
	s := new(session)
	s.pool = pool
	s.in = newQueue()
	s.out = newQueue()
	return s
}

func (s *session) Receive() ([]byte, error) {
	return s.in.pull()
}

func (s *session) Send(m []byte) error {
	s.out.push(m)
	return nil
}

func (s *session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed_ = true
	s.in.close()
	s.out.close()
	return nil
}

// Reserve marks the session reserved so that other connections know not read from it.
// True is returned, if the reservation fails, otherwise false.
func (s *session) reserve() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reserved {
		return true
	}
	s.reserved = true
	return false
}

// Free marks the session free for read for other connections.
func (s *session) free() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reserved = false
}

// Closed returns true, if the session is closed, otherwise false.
func (s *session) closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed_
}
