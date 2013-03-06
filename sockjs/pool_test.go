package sockjs

import (
	. "launchpad.net/gocheck"
)

type SessionPoolSuite struct {
	p *SessionPool
}

var _ = Suite(&SessionPoolSuite{})

func (s *SessionPoolSuite) SetUpTest(c *C) {
	s.p = NewSessionPool()
}

func (s *SessionPoolSuite) TestAdd(c *C) {
	session := new(legacySession)
	s.p.Add(session)
	_, exists := s.p.pool[session]
	c.Check(exists, Equals, true)
}

func (s *SessionPoolSuite) TestRemove(c *C) {
	session := new(legacySession)
	s.p.Add(session)
	s.p.Remove(session)
	_, exists := s.p.pool[session]
	c.Check(exists, Equals, false)
}

type LegacyPoolSuite struct {
	p *legacyPool
}

var _ = Suite(&LegacyPoolSuite{})

func (s *LegacyPoolSuite) SetUpTest(c *C) {
	s.p = newLegacyPool()
}

func (s *LegacyPoolSuite) TestGet(c *C) {
	session := new(legacySession)
	s.p.pool["foo"] = session
	c.Check(s.p.get("foo"), Equals, session)
}

func (s *LegacyPoolSuite) TestGetOrCreate(c *C) {
	session, exists := s.p.getOrCreate("foo")
	c.Assert(session, DeepEquals, session)
	c.Assert(exists, Equals, false)

	r, exists := s.p.getOrCreate("foo")
	c.Assert(r, DeepEquals, session)
	c.Assert(exists, Equals, true)
}

func (s *LegacyPoolSuite) TestRemove(c *C) {
	session := new(legacySession)
	s.p.pool["foo"] = session
	s.p.remove("foo")
	_, exists := s.p.pool["foo"]
	c.Check(exists, Equals, false)
}
