package sockjs

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

type nopHandler int

func (n nopHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func (s *ServerSuite) TestServerMatch(c *C) {
	alt := nopHandler(0)
	long := nopHandler(1)
	short := nopHandler(2)

	server := NewServer(alt)
	defer server.Close()
	server.m["/prefix/long"] = long
	server.m["/prefix"] = short

	c.Check(server.match("/prefix/long"), Equals, long)
	c.Check(server.match("/prefix/long/foobar?zot=5&baz=2"), Equals, long)
	c.Check(server.match("/prefix"), Equals, short)
	c.Check(server.match("/prefix/foobar?zot=5&baz=2"), Equals, short)
	c.Check(server.match("/notfound"), Equals, alt)
}
