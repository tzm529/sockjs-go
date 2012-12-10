package sockjs

import(
	"io"
	"errors"
)

type connKind uint8

const (
	connKindWebsocket connKind = iota
	connKindRawWebsocket
)

type Conn struct {
	kind connKind
	wc io.WriteCloser
	closed bool
}

func (c *Conn) Send(s string) (err error) {
	if c.closed {
		return errors.New("connection closed")
	}

	switch c.kind {
	case connKindWebsocket:
		sendWebsocket(c.wc, s)
	case connKindRawWebsocket:
		sendRawWebsocket(c.wc, s)
	}
	return 
}

func (c *Conn) Close() (err error) {
	switch c.kind {
	case connKindWebsocket:
		err = closeWebsocket(c.wc)
	case connKindRawWebsocket:
		err = closeRawWebsocket(c.wc)
	}
	c.closed = true
	return
}
