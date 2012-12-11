package sockjs

import (
	. "launchpad.net/gocheck"
)

type PoolSuite struct{
	p *pool
}

var _ = Suite(&PoolSuite{})

func (s *PoolSuite) SetUpTest(c *C) {
	s.p = newPool()
}

func (s *PoolSuite) TestPoolGet(c *C) {
	session := protoRawWebsocket{nil}
	s.p.pool["foo"] = session
	c.Check(s.p.get("foo"), Equals, session)
}

func (s *PoolSuite) TestPoolSet(c *C) {
	session := protoRawWebsocket{nil}
	s.p.set("foo", session)
	c.Check(s.p.pool["foo"], Equals, session)
}

func (s *PoolSuite) TestPoolRemove(c *C) {
	session := protoRawWebsocket{nil}
	s.p.pool["foo"] = session
	s.p.remove("foo")
	_, exists := s.p.pool["foo"]
	c.Check(exists, Equals, false)
}

