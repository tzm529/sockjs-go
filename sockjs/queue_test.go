package sockjs

import (
	. "launchpad.net/gocheck"
)

type QueueSuite struct {
	q *queue
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestPull(c *C) {
	s.q = newQueue(false)
	defer s.q.close()

	s.q.push([]byte{'a'}, []byte{'b'}, []byte{'c'})

	v, err := s.q.pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'a'})

	v, err = s.q.pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'b'})

	v, err = s.q.pull()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'c'})
}

func (s *QueueSuite) TestPullAll(c *C) {
	s.q = newQueue(false)
	defer s.q.close()

	s.q.push([]byte{'a'}, []byte{'b'}, []byte{'c'})

	v, err := s.q.pullAll()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, [][]byte{{'a'}, {'b'}, {'c'}})
}

func (s *QueueSuite) TestClosedPullError(c *C) {
	s.q = newQueue(false)
	defer s.q.close()

	s.q.push([]byte{'a'})
	s.q.close()

	_, err := s.q.pull()
	c.Assert(err, Equals, errQueueClosed)
}

func (s *QueueSuite) TestClosedPushPanic(c *C) {
	s.q = newQueue(false)
	defer s.q.close()

	s.q.push([]byte{'a'})
	s.q.close()

	c.Assert(func() { s.q.push([]byte{'b'}) }, Panics, errQueueClosed)
}

func (s *QueueSuite) TestWaitPullError(c *C) {
	s.q = newQueue(true)
	defer s.q.close()

	f := func() {
		_, err := s.q.pull()
		if !(err == errQueueClosed || err == errQueueWait) {
			c.Fatal("wrong error value")
		}
	}

	go f()
	go f()
}
