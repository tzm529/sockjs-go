package sockjs

import (
	. "launchpad.net/gocheck"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) TestDataFrame(c *C) {
	c.Assert(aframe([][]byte{
		{'a', 'b', 'c'},
		{'d', 'e', 'f'},
		{'g', 'h', 'i'},
	}...), DeepEquals, []byte(`a["abc","def","ghi"]`))
}

func (s *UtilsSuite) TestCloseFrame(c *C) {
	c.Assert(cframe(3210, "multifail"),
		DeepEquals,
		[]byte(`c[3210,"multifail"]`))
}

func (s *UtilsSuite) TestVerifyAddr(c *C) {
	c.Assert(verifyAddr("foo:123", "foo:123"), Equals, true)
	c.Assert(verifyAddr("foo:123", "foo:456"), Equals, true)
	c.Assert(verifyAddr("foo:123", "bar:123"), Equals, false)
	c.Assert(verifyAddr("foo:123", "bar:456"), Equals, false)
}

func (s *UtilsSuite) BenchmarkDataFrame(c *C) {
	for i := 0; i < c.N; i++ {
		aframe([][]byte{
			{'a', 'b', 'c'},
			{'d', 'e', 'f'},
			{'g', 'h', 'i'},
		}...)
	}
}
