package sockjs

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type ServeMuxSuite struct{}

var _ = Suite(&ServeMuxSuite{})

// plain http.Handler is not comparable
type nopHandler int

func (n nopHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func (s *ServeMuxSuite) TestMatch(c *C) {
	conf := NewConfig()
	alt := nopHandler(0)
	long := newHandler("/prefix/long", nil, &conf)
	short := newHandler("/prefix", nil, &conf)

	mux := NewServeMux(alt)
	mux.m["/prefix/long"] = long
	mux.m["/prefix"] = short

	c.Check(mux.match("/prefix/long"), Equals, long)
	c.Check(mux.match("/prefix/long/foobar?zot=5&baz=2"), Equals, long)
	c.Check(mux.match("/prefix"), Equals, short)
	c.Check(mux.match("/prefix/foobar?zot=5&baz=2"), Equals, short)
	c.Check(mux.match("/notfound"), Equals, alt)
}
