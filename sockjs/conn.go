package sockjs

import(
	"io"
	"errors"
)

type Conn struct {
	wc io.WriteCloser
	closed bool
}

func newConn(wc io.WriteCloser) (c *Conn) {
	c = new(Conn)
	c.wc = wc
	return
}

func (c *Conn) send(s []byte) (err error) {
	if c.closed {
		return errors.New("connection closed")
	}

	_, err = c.wc.Write(s)
	return
}

func (c *Conn) Send(s ...string) (err error) {
	return c.send(aframe(s...))
}

func (c *Conn) Close() (err error) {
	c.send([]byte(`c[3000,"Go away!"]`))
	err = c.wc.Close()
	c.closed = true
	return
}
