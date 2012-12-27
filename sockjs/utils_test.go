package sockjs

import (
	. "launchpad.net/gocheck"
)

type UtilsSuite struct {}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) TestDataFrame(c *C) {
	c.Assert(aframe("foo", "bar", [][]byte{
		{'a','b','c'},
		{'d','e','f'},
		{'g','h','i'},
	}...), DeepEquals, []byte(`fooa["abc","def","ghi"]bar`))
}

func (s *UtilsSuite) TestCloseFrame(c *C) {
	c.Assert(cframe("foo", 3210, "multifail", "bar"), 
		DeepEquals, 
		[]byte(`fooc[3210,"multifail"]bar`))
}

func (s *UtilsSuite) BenchmarkDataFrame(c *C) { 
	for i := 0; i < c.N; i++ { 
        aframe("foo", "bar", [][]byte{
			{'a','b','c'},
			{'d','e','f'},
			{'g','h','i'},
		}...)
	}
} 