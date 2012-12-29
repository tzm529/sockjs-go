package sockjs

import (
	"errors"
	"net/http"
	"sync"
)

var ErrSessionClosed error = errors.New("session closed")
var ErrSessionTimeout error = errors.New("session timeout")

type Session interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
	Info() RequestInfo
	Protocol() Protocol
}

// structure for polling sessions
type session struct {
	proto protocol
	in    *queue
	out   *queue

	mu           sync.Mutex
	closed_      bool
	interrupted_ bool
	reserved     bool
	info         *RequestInfo
}

func (s *session) init(r *http.Request,
	prefix string,
	protocol protocol,
	headers []string) {
	s.in = newQueue()
	s.out = newQueue()
	s.info = newRequestInfo(r, prefix, headers)
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

func (s *session) Info() RequestInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return *s.info
}

func (s *session) Protocol() Protocol {
	return s.proto.protocol()
}

func (s *session) setInfo(info *RequestInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.info = info
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

// Interrupted returns true, if the session was interrupted.
func (s *session) interrupted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.interrupted_
}

// Interrupt marks the session as interrupted.
func (s *session) interrupt() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interrupted_ = true
}

// VerifyAddr returns true, if the given remote address matches the one used in the last request,
// otherwise false.
func (s *session) verifyAddr(addr string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return verifyAddr(s.info.RemoteAddr, addr)
}
