package sockjs

import (
	. "launchpad.net/gocheck"
)

type QueueSuite struct {
	q *queue
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestPull(c *C) {
	s.q = newQueue()

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

	// closed error
	s.q.push([]byte{'a'})
	s.q.close()
	v, err = s.q.pull()
	c.Assert(v, IsNil)
	c.Assert(err, Equals, errQueueClosed)
}

func (s *QueueSuite) TestPullAll(c *C) {
	s.q = newQueue()

	s.q.push([]byte{'a'}, []byte{'b'}, []byte{'c'})

	v, err := s.q.pullAll()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, [][]byte{{'a'}, {'b'}, {'c'}})

	// closed error
	s.q.push([]byte{'a'})
	s.q.close()
	v, err = s.q.pullAll()
	c.Assert(v, IsNil)
	c.Assert(err, Equals, errQueueClosed)
}

func (s *QueueSuite) TestPullNow(c *C) {
	s.q = newQueue()

	v, err := s.q.pullNow()
	c.Assert(v, IsNil)
	c.Assert(err, IsNil)

	s.q.push([]byte{'a'}, []byte{'b'}, []byte{'c'})

	v, err = s.q.pullNow()
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte{'a'})

	// closed error
	s.q.push([]byte{'a'})
	s.q.close()
	v, err = s.q.pullNow()
	c.Assert(v, IsNil)
	c.Assert(err, Equals, errQueueClosed)
}


func (s *QueueSuite) TestClosedPullError(c *C) {
	s.q = newQueue()
	defer s.q.close()

}

func (s *QueueSuite) TestClosedPushPanic(c *C) {
	s.q = newQueue()
	defer s.q.close()

	s.q.push([]byte{'a'})
	s.q.close()

	c.Assert(func() { s.q.push([]byte{'b'}) }, Panics, errQueueClosed)
}
