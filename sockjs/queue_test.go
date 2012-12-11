package sockjs

import (
	. "launchpad.net/gocheck"
)

type QueueSuite struct{
	q *queue
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) SetUpTest(c *C) {
	s.q = newQueue()
}

func (s *QueueSuite) TearDownTest(c *C) {
	s.q.close()
}

func (s *QueueSuite) TestQueue(c *C) {
	s.q.push([]byte{'a'})
	s.q.push([]byte{'b'})
	s.q.push([]byte{'c'})

	v, err := s.q.pull()
	c.Check(v, DeepEquals, []byte{'a'})
	c.Check(err, IsNil)

	v, err = s.q.pull()
	c.Check(v, DeepEquals, []byte{'b'})
	c.Check(err, IsNil)

	v, err = s.q.pull()
	c.Check(v, DeepEquals, []byte{'c'})
	c.Check(err, IsNil)
}

func (s *QueueSuite) TestQueueClose(c *C) {
	s.q.push([]byte{'a'})
	s.q.close()

	err := s.q.push([]byte{'b'})
	c.Check(err, Equals, ErrSessionClosed)

	_, err = s.q.pull()
	c.Check(err, Equals, ErrSessionClosed)	
}

