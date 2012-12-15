package sockjs

import (
	. "launchpad.net/gocheck"
)

type QueueSuite struct{
	q *queue
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestPull(c *C) {
	s.q = newQueue(false)
	defer s.q.Close()

	s.q.Push([]byte{'a'})
	s.q.Push([]byte{'b'})
	s.q.Push([]byte{'c'})

	v, err := s.q.Pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'a'})

	v, err = s.q.Pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'b'})

	v, err = s.q.Pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'c'})
}

func (s *QueueSuite) TestPullAll(c *C) {
	s.q = newQueue(false)
	defer s.q.Close()

	s.q.Push([]byte{'a'})
	s.q.Push([]byte{'b'})
	s.q.Push([]byte{'c'})
	
	v, err := s.q.PullAll()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, [][]byte{
		[]byte{'a'},
		[]byte{'b'},
		[]byte{'c'}})
}

func (s *QueueSuite) TestClosedPullError(c *C) {
	s.q = newQueue(false)
	defer s.q.Close()

	s.q.Push([]byte{'a'})
	s.q.Close()

	_, err := s.q.Pull()
	c.Assert(err, Equals, errQueueClosed)
}

func (s *QueueSuite) TestClosedPushPanic(c *C) {
	s.q = newQueue(false)
	defer s.q.Close()

	s.q.Push([]byte{'a'})
	s.q.Close()

	c.Assert(func() { s.q.Push([]byte{'b'}) }, Panics, errQueueClosed)
}

func (s *QueueSuite) TestWaitPullError(c *C) {
	s.q = newQueue(true)
	defer s.q.Close()

	f := func() {
		_, err := s.q.Pull()
		if !(err == errQueueClosed || err == errQueueWait) {
			c.Fatal("wrong error value")
		}
	}

	go f()
	go f()
}

