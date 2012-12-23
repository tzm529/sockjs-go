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

func (s *PoolSuite) TearDownTest(c *C) {
	s.p.close()
}

func (s *PoolSuite) TestPoolGet(c *C) {
	session := newSession(s.p)
	s.p.pool["foo"] = session
	c.Check(s.p.get("foo"), Equals, session)
}

func (s *PoolSuite) TestPoolGetOrCreate(c *C) {
	session_ := newSession(s.p)

	r, exists := s.p.getOrCreate("foo")
	c.Assert(r, DeepEquals, session_)
	c.Assert(exists, Equals, false)

	r, exists = s.p.getOrCreate("foo")
	c.Assert(r, DeepEquals, session_)
	c.Assert(exists, Equals, true)
}

func (s *PoolSuite) TestPoolRemove(c *C) {
	session := newSession(s.p)
	s.p.pool["foo"] = session
	s.p.remove("foo")
	_, exists := s.p.pool["foo"]
	c.Check(exists, Equals, false)
}
