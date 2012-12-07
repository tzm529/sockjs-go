package sockjs

import (
	. "launchpad.net/gocheck"
	"testing"
	"net/http"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) {
	TestingT(t)
}

type ServeMuxSuite struct{}

func init() {
	Suite(&ServeMuxSuite{})
}

type nopHandler int

func (n nopHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func (s *ServeMuxSuite) TestServeMuxMatch(c *C) {
	alt := nopHandler(0)
	long := nopHandler(1)
	short := nopHandler(2)

	mux := NewServeMux(alt)
	mux.m["/prefix/long"] = long
	mux.m["/prefix"] = short

	c.Check(mux.match("/prefix/long"), Equals, long)
	c.Check(mux.match("/prefix/long/foobar?zot=5&baz=2"), Equals, long)
	c.Check(mux.match("/prefix"), Equals, short)
	c.Check(mux.match("/prefix/foobar?zot=5&baz=2"), Equals, short)
	c.Check(mux.match("/notfound"), Equals, alt)
}


