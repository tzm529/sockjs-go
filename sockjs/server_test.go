package sockjs

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

// plain http.Handler is not comparable
type nopHandler int

func (n nopHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func (s *ServerSuite) TestServerMatch(c *C) {
	conf := NewConfig()
	alt := nopHandler(0)
	long := newHandler(nil, "/prefix/long", nil, &conf)
	short := newHandler(nil, "/prefix", nil, &conf)

	server := NewServer(alt)
	server.m["/prefix/long"] = long
	server.m["/prefix"] = short

	c.Check(server.match("/prefix/long"), Equals, long)
	c.Check(server.match("/prefix/long/foobar?zot=5&baz=2"), Equals, long)
	c.Check(server.match("/prefix"), Equals, short)
	c.Check(server.match("/prefix/foobar?zot=5&baz=2"), Equals, short)
	c.Check(server.match("/notfound"), Equals, alt)
}
