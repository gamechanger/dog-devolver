package devolve

import (
	"testing"

	. "gopkg.in/check.v1"
)

// wire up gocheck
func Test(t *testing.T) { TestingT(t) }

type DevolveTestSuite struct{}

var _ = Suite(&DevolveTestSuite{})

func (s *DevolveTestSuite) TestDevolveWithEmptyStringErrors(c *C) {
	in := ""
	result, err := Devolve(in)
	c.Assert(result, Equals, "")
	c.Assert(err, NotNil)
}

func (s *DevolveTestSuite) TestDevolveGarbageInGarbageOut(c *C) {
	in := "I AM A POTATO"
	result, err := Devolve(in)
	c.Assert(result, Equals, "")
	c.Assert(err, NotNil)
}

func (s *DevolveTestSuite) TestDevolveWithStatsDCounter(c *C) {
	in := "my.metric:1|c"
	result, err := Devolve(in)
	c.Assert(result, Equals, "my.metric:1|c")
	c.Assert(err, IsNil)
}

func (s *DevolveTestSuite) TestDevolveWithDogStatsDCounterPartial(c *C) {
	in := "my.metric:1|c|@0.1"
	result, err := Devolve(in)
	c.Assert(result, Equals, "my.metric:1|c|@0.1")
	c.Assert(err, IsNil)
}

func (s *DevolveTestSuite) TestDevolveWithDogStatsDCounterFull(c *C) {
	in := "my.metric:1|c|@0.1|#env:production,1420"
	result, err := Devolve(in)
	c.Assert(result, Equals, "my.metric:1|c|@0.1")
	c.Assert(err, IsNil)
}

func (s *DevolveTestSuite) TestDevolveWithDogStatsDEvent(c *C) {
	in := "_e{10,50}:Cataclysm!:Literally everything exploded and you died, sorry."
	result, err := Devolve(in)
	c.Assert(result, Equals, "")
	c.Assert(err, NotNil)
}
