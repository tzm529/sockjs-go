package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"container/list"
	"sync"
)

type protocol uint8

const (
	protocolWebsocket protocol = iota
	protocolRawWebsocket
	protocolXhrPolling
)

type Session struct {
	proto   protocol
	queue  list.List  // message queue
	mu     sync.Mutex // lock for message queue
	ws     *websocket.Conn
}

func (s *Session) Receive() (m []byte, err error) {
	pm := s.pull()
	if pm != nil {
		// receive from queue
		return pm, nil
	} else if s.proto == protocolWebsocket || s.proto  == protocolRawWebsocket {
		// receive from connection
		switch s.proto {
		case protocolWebsocket:
			m, err = receiveWebsocket(s)
		case protocolRawWebsocket:
			m, err = receiveRawWebsocket(s)
		default:
			panic("unknown protocol")
		}
	}
	return
}

func (s *Session) Send(m []byte) (err error) {
	switch s.proto {
	case protocolWebsocket:
		err = sendWebsocket(s, m)
	case protocolRawWebsocket:
		err = sendRawWebsocket(s, m)
	case protocolXhrPolling:
		err = sendXhrPolling(s, m)
	default:
		panic("unknown protocol")
	}
	return
}

func (s *Session) Close() (err error) {
	switch s.proto {
	case protocolWebsocket:
		err = closeWebsocket(s)
	case protocolRawWebsocket:
		err = closeRawWebsocket(s)
	case protocolXhrPolling:
		err = closeXhrPolling(s)
	default:
		panic("unknown protocol")
	}
	return
}

// pull a message from the message queue
func (s *Session) pull() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.queue.Front()
	if e == nil {
		return nil
	}
	m, _ := s.queue.Remove(e).([]byte)
	return m
}

// push a message to the message queue
func (s *Session) push(m []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue.PushBack(m)
}
