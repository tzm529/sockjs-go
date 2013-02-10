package sockjs

import (
	"errors"
	"net/http"
	"sync"
	"time"
//	"container/list"
)

var ErrSessionClosed error = errors.New("session closed")
var ErrSessionTimeout error = errors.New("session timeout")

type readEnvelope struct {
	m []byte
	err error
}

type writeEnvelope struct {
	m []byte
	rv chan error
}

type bufEnvelope struct {
	buf *list.List
	done chan struct{}
}

type Session interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
	Info() RequestInfo
	Protocol() Protocol
}

// structure for polling sessions
type session struct {
	// read-only
	proto protocol

	in chan []byte
	rio sync.Mutex
	rbuf *list


	closer chan struct{}
	inBuf chan chan *bufEnvelope
	inOut chan []byte
	out chan *writeEnvelope
	wc chan io.WriteCloser
	hbTicker *time.Ticker
	dcTicker *time.Ticker

	mu sync.RWMutex
	interrupted_ bool
	reserved     bool
	info         *RequestInfo
}

func (s *session) init(r *http.Request, proto protocol, config *config) {
	s.proto = proto
	s.hbTicker = time.NewTicker(config.HeartbeatDelay)
	s.dcTicker = time.NewTicker(config.DisconnectionDelay)
	go backend()
}

func (s *session) backend() {
	rbuf := list.New()
	wbuf := list.New()
	defer close(s.closer)
	defer close(s.inBuf)
	defer close(s.inOut)
	defer s.hbTicker.Stop()
	defer s.dcTicker.Stop()

	for {
		if rbuf.Len() > 0 {
			// try sending a messages to Read()
			front := rbuf.Front()
			select {
			default:
			case s.inOut <- front.Value.([]byte):
				rbuf.Remove(front)
			}
		}


		rbufDone := make(chan struct{})
		select {
		case <-s.closer:
			return

			// loan rbuf to a reader
		case s.inBuf <- &bufEnvelope{rbuf, rbufDone):
			<-rbufDone

		case <-s.dcTicker:
			// TODO
		}
	}
}

func (s *session) Receive() ([]byte, error) {
	m, closed := <-s.inOut
	if closed { return nil, ErrSessionClosed }
	return m, nil
}

func (s *session) Send(m []byte) (err error) {
	rv := make(chan error)
	select {
	case <-s.closer:
		return ErrSessionClosed
	case s.out <- &connWrite(m, rv):
	}
	return <-rv
}

func (s *session) Close() error {
	select{
	default:
	case <-s.closer:
	}
}

func (s *session) Info() RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
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

// Interrupted returns true, if the session was interrupted.
func (s *session) interrupted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.interrupted_
}

// Interrupt marks the session interrupted.
func (s *session) interrupt() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interrupted_ = true
}

// VerifyAddr returns true, if the given remote address matches the one used in the last request,
// otherwise false.
func (s *session) verifyAddr(addr string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return verifyAddr(s.info.RemoteAddr, addr)
}
