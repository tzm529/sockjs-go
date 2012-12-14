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

func (s *QueueSuite) TestPull(c *C) {
	s.q.push([]byte{'a'})
	s.q.push([]byte{'b'})
	s.q.push([]byte{'c'})
	c.Assert(s.q.pull(), DeepEquals, []byte{'a'})
	c.Assert(s.q.pull(), DeepEquals, []byte{'b'})
	c.Assert(s.q.pull(), DeepEquals, []byte{'c'})
}

func (s *QueueSuite) TestPullAll(c *C) {
	s.q.push([]byte{'a'})
	s.q.push([]byte{'b'})
	s.q.push([]byte{'c'})
	c.Assert(s.q.pullAll(), DeepEquals, [][]byte{
		[]byte{'a'},
		[]byte{'b'},
		[]byte{'c'}})
}

func (s *QueueSuite) TestClose(c *C) {
	s.q.push([]byte{'a'})
	s.q.close()

	c.Assert(func() { s.q.push([]byte{'b'}) }, Panics, errQueueClosed)
	c.Check(s.q.pull(), IsNil)
}

