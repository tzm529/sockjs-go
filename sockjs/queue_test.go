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
	c.Assert(s.q.pull(), DeepEquals, []byte{'a'})
	c.Assert(s.q.pull(), DeepEquals, []byte{'b'})
	c.Assert(s.q.pull(), DeepEquals, []byte{'c'})
}

func (s *QueueSuite) TestQueueClose(c *C) {
	s.q.push([]byte{'a'})
	s.q.close()

	c.Assert(func() { s.q.push([]byte{'b'}) }, Panics, errQueueClosed)
	c.Check(s.q.pull(), IsNil)
}

