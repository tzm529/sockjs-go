package sockjs

import(
	"code.google.com/p/go.net/websocket"
	"errors"
"container/list"
"sync"
)

type sessionKind uint8

const (
	sessionKindWebsocket sessionKind = iota
	sessionKindRawWebsocket
)

type Session struct {
	kind sessionKind
	closed bool
	queue list.List // message queue
	mu sync.Mutex // lock for message queue
	ws *websocket.Conn
}

// pull a message from the message queue
func (s *Session) pull() *string {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.queue.Front()
	if e == nil { return nil }
	m, _ := s.queue.Remove(e).(string)
	return &m
}

// push a message to the message queue
func (s *Session) push(m string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue.PushBack(m)
}

func (s *Session) Receive() (m string, err error) {
	pm := s.pull()
	if pm != nil {
		// receive from queue
		return *pm, nil
	} else {
		// receive from connection
		if s.closed {
			return "", errors.New("connection closed")
		}
		
		switch s.kind {
		case sessionKindWebsocket:
			m, err = receiveWebsocket(s)
		case sessionKindRawWebsocket:
			m, err = receiveRawWebsocket(s)
		}
	}
	return 
}

func (s *Session) Send(m string) (err error) {
	if s.closed {
		return errors.New("connection closed")
	}

	switch s.kind {
	case sessionKindWebsocket:
		err = sendWebsocket(s, m)
	case sessionKindRawWebsocket:
		err = sendRawWebsocket(s, m)
	}
	return 
}

func (s *Session) Close() (err error) {
	switch s.kind {
	case sessionKindWebsocket:
		err = closeWebsocket(s)
	case sessionKindRawWebsocket:
		err = closeRawWebsocket(s)
	}
	s.closed = true
	return
}
