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

func (s *QueueSuite) TestQueue(c *C) {
	s.q.push([]byte{'a'})
	s.q.push([]byte{'b'})
	s.q.push([]byte{'c'})
	c.Check(s.q.pull(), DeepEquals, []byte{'a'})
	c.Check(s.q.pull(), DeepEquals, []byte{'b'})
	c.Check(s.q.pull(), DeepEquals, []byte{'c'})
}

