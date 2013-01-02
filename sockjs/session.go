package sockjs

import (
	"errors"
	"net/http"
	"sync"
	"time"
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
	proto protocol // read-only
	in    *queue
	out   *queue

	mu           sync.RWMutex
	closed_      bool
	interrupted_ bool
	reserved     bool
	timeouted_ bool
	info         *RequestInfo
	lastMsgTime_  time.Time
}

func (s *session) init(r *http.Request,
	prefix string,
	proto protocol,
	headers []string) {
	s.proto = proto
	s.in = newQueue()
	s.out = newQueue()
	s.info = newRequestInfo(r, prefix, headers)
}

func (s *session) Receive() (m []byte, err error) {
	m, err = s.in.pull()

	switch {
	case err == nil:
		s.setLastMsgTime(time.Now())
	case s.timeouted():
		err = ErrSessionTimeout
	default: // errQueueClosed
		err = ErrSessionClosed
	}
	return
}

func (s *session) Send(m []byte) error {
	s.out.push(m)
	return nil
}

func (s *session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed_ = true
	s.cleanup()
	return nil
}

func (s *session) Info() RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.info
}

func (s *session) Protocol() Protocol {
	return s.proto.protocol()
}

func (s *session) timeouted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.timeouted_
}

func (s *session) timeout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timeouted_ = true
}

func (s *session) lastMsgTime() time.Time{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastMsgTime_
}

func (s *session) setLastMsgTime(lastMsgTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastMsgTime_ = lastMsgTime
}

func (s *session) cleanup() {
	s.in.close()
	s.out.close()
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed_
}

// Interrupted returns true, if the session was interrupted.
func (s *session) interrupted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.interrupted_
}

// Interrupt marks the session closed and interrupted.
func (s *session) interrupt() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed_ = true
	s.interrupted_ = true
}

// VerifyAddr returns true, if the given remote address matches the one used in the last request,
// otherwise false.
func (s *session) verifyAddr(addr string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return verifyAddr(s.info.RemoteAddr, addr)
}
