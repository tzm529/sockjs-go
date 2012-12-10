package sockjs

import(
	"code.google.com/p/go.net/websocket"
	"errors"
)

type sessionKind uint8

const (
	sessionKindWebsocket sessionKind = iota
	sessionKindRawWebsocket
)

type Session struct {
	kind sessionKind
	ws *websocket.Conn
	closed bool
}

func (s *Session) Receive() (m string, err error) {
	if s.closed {
		return "", errors.New("connection closed")
	}

	switch s.kind {
	case sessionKindWebsocket:
		m, err = receiveWebsocket(s.ws)
	case sessionKindRawWebsocket:
		m, err = receiveRawWebsocket(s.ws)
	}
	return 
}

func (s *Session) Send(m string) (err error) {
	if s.closed {
		return errors.New("connection closed")
	}

	switch s.kind {
	case sessionKindWebsocket:
		err = sendWebsocket(s.ws, m)
	case sessionKindRawWebsocket:
		err = sendRawWebsocket(s.ws, m)
	}
	return 
}

func (s *Session) Close() (err error) {
	switch s.kind {
	case sessionKindWebsocket:
		err = closeWebsocket(s.ws)
	case sessionKindRawWebsocket:
		err = closeRawWebsocket(s.ws)
	}
	s.closed = true
	return
}
