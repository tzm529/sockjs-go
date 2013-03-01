package sockjs

import (
	"sync"
	"time"
	"container/list"
)

// Session describes a sockjs session.
type Session interface {
	// Receive blocks until a message can be returned from session receive buffer or nil, 
	// if the session is closed.
	Receive() (m []byte)

	// Send appends the given message to session send buffer.
	// Panics, if the session is closed.
	Send(m []byte)

	// Close closes the session.
	// Pending sends will be discarded unless the client receives them within 
	// Config.DisconnectDelay.
	Close(code int, reason string)

	// End is a convenience method for closing with the default code and reason, 
	// Close(3000, "Go away!").
	End()

	// Info returns a RequestInfo object containing information copied from the last received 
	// request.
	Info() RequestInfo

	// Protocol returns the underlying protocol of the session.
	Protocol() Protocol
}

// Session for legacy protocols.
type legacySession struct {
	// read-only
	proto protocol
	config *Config
	sessid string
	pool *pool

	sendBuffer chan []byte
	sendFrame chan []byte
	hbTicker *time.Ticker
	dcTicker *time.Ticker

	rio sync.Mutex
	rbufEmpty *sync.Cond
	rbuf *list.List

	mu sync.RWMutex
	closed_ bool
	sendBufClosed bool
	closeCode int
	closeReason string
	info         *RequestInfo
	reserved     bool
	recvStamp time.Time
}

//* Public methods

func (s *legacySession) Receive() []byte {
	s.rio.Lock()
	defer s.rio.Unlock()

	for s.rbuf.Len() == 0 {
		if s.closed() { return nil }
		s.rbufEmpty.Wait()
	}

	return s.rbuf.Remove(s.rbuf.Front()).([]byte)
}

func (s *legacySession) Send(m []byte) {
	s.sendBuffer <- m
}

func (s *legacySession) End() {
	s.Close(3000, "Go away!")
}

func (s *legacySession) Close(code int, reason string) {
	s.close(code, reason)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sendBufClosed { return }
	s.sendBufClosed = true
	close(s.sendBuffer)
}

func (s *legacySession) Protocol() Protocol {
	return s.proto.protocol()
}

func (s *legacySession) Info() RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.info
}

//* Private methods

func (s *legacySession) init(config *Config, proto protocol, sessid string, pool *pool) {
	s.config = config
	s.proto = proto
	s.sessid = sessid
	s.pool = pool
	s.rbufEmpty = sync.NewCond(&s.rio)
	s.rbuf = list.New()
	s.sendBuffer = make(chan []byte)
	s.sendFrame = make(chan []byte)
	s.hbTicker = time.NewTicker(config.HeartbeatDelay)
	s.dcTicker = time.NewTicker(config.DisconnectDelay)
	go s.backend()
	go s.sendBuffer_(s.sendBuffer, s.sendFrame)
}

func (s *legacySession) backend() {
loop:
	for {
		select {
		case <-s.hbTicker.C:
			s.sendFrame <- []byte{'h'}

		case <-s.dcTicker.C:
			if s.timeouted() { 
				// close in case it wasn't closed already
				s.end()
				break loop
			}
		}
	}

	s.dcTicker.Stop()
	s.pool.remove(s.sessid)
}

func (s *legacySession) rbufAppend(m []byte) {
	s.rio.Lock()
	defer s.rio.Unlock()

	s.rbuf.PushBack(m)
	s.rbufEmpty.Signal()
}

func (s *legacySession) sendBuffer_(in <-chan []byte, out chan<- []byte) {
 	var pending [][]byte
	defer close(out)
	
loop:
	for {
		// keep pending non-empty
		if len(pending) == 0 {
			v, ok := <-in
			if !ok {
				break
			}
			pending = append(pending, v)
		}

		select {
		case v, ok := <-in:
			if !ok {
				break loop
			}
			pending = append(pending, v)

		case out <- aframe(pending...):
			pending = nil
		}
	}
	
	// Try sending the remaining values, but don't wait more than Config.DisconnectDelay.
	if len(pending) > 0 {
		select {
		case out <- aframe(pending...):
		case <-time.After(s.config.DisconnectDelay):
		}
	}
}

func (s *legacySession) closeFrame() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cframe(s.closeCode, s.closeReason)
}

// Closed returns true, if the session is closed.
func (s *legacySession) closed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed_
}

func (s *legacySession) end() {
	s.close(3000, "Go away!")
}

func (s *legacySession) close(code int, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed_ { return }
	s.closed_ = true
	s.closeCode = code
	s.closeReason = reason
	s.hbTicker.Stop()
	s.rbufEmpty.Broadcast()
}

func (s *legacySession) setInfo(info *RequestInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.info = info
}

// Reserve marks the session reserved so that other connections know not receive from it.
// False is returned, if the reservation fails, otherwise true.
func (s *legacySession) reserve() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reserved {
		return false
	}
	s.reserved = true
	return true
}

// Free marks the session free for receiving for other connections.
func (s *legacySession) free() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reserved = false
}

// VerifyAddr returns true, if the given remote address matches the one used in the last request,
// otherwise false.
func (s *legacySession) verifyAddr(addr string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return verifyAddr(s.info.RemoteAddr, addr)
}

func (s *legacySession) timeouted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if time.Since(s.recvStamp) > time.Duration(s.config.DisconnectDelay)*time.Second {
		return true
	}
	return false
}

func (s *legacySession) updateRecvStamp() {
	s.mu.Lock()
	defer s.mu.Unlock()	
	s.recvStamp = time.Now()
}
