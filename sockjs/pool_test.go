package sockjs

import (
	. "launchpad.net/gocheck"
)

type PoolSuite struct {
	p *pool
}

var _ = Suite(&PoolSuite{})

func (s *PoolSuite) SetUpTest(c *C) {
	s.p = newPool()
}

func (s *PoolSuite) TestGet(c *C) {
	session := new(legacySession)
	s.p.pool["foo"] = session
	c.Check(s.p.get("foo"), Equals, session)
}

func (s *PoolSuite) TestGetOrCreate(c *C) {
	session, exists := s.p.getOrCreate("foo")
	c.Assert(session, DeepEquals, session)
	c.Assert(exists, Equals, false)

	r, exists := s.p.getOrCreate("foo")
	c.Assert(r, DeepEquals, session)
	c.Assert(exists, Equals, true)
}

func (s *PoolSuite) TestRemove(c *C) {
	session := new(legacySession)
	s.p.pool["foo"] = session
	s.p.remove("foo")
	_, exists := s.p.pool["foo"]
	c.Check(exists, Equals, false)
}
