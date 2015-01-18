package config

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

// wire up gocheck
func Test(t *testing.T) { TestingT(t) }

type ConfigTestSuite struct{}

var _ = Suite(&ConfigTestSuite{})

func (s *ConfigTestSuite) TearDownTest(c *C) {
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("TEST_VAR_0")
	os.Unsetenv("TEST_VAR_1")
	os.Unsetenv("TEST_VAR_2")
	os.Unsetenv("TEST_VAR_3")
}

func (s *ConfigTestSuite) TestDefaultValue(c *C) {
	os.Setenv("TEST_VAR", "")
	c.Assert(os.Getenv("TEST_VAR"), Equals, "")
	c.Assert(defaultValue(os.Getenv("TEST_VAR"), "groovy"), Equals, "groovy")
}

func (s *ConfigTestSuite) TestConfigSliceNoData(c *C) {
	c.Assert(configSlice("TEST_VAR"), DeepEquals, []string(nil))
}

func (s *ConfigTestSuite) TestConfigSliceOneElement(c *C) {
	os.Setenv("TEST_VAR_0", "Melchior")
	c.Assert(configSlice("TEST_VAR"), DeepEquals, []string{"Melchior"})
}

func (s *ConfigTestSuite) TestConfigSliceThreeElements(c *C) {
	os.Setenv("TEST_VAR_0", "Melchior")
	os.Setenv("TEST_VAR_1", "Belthasar")
	os.Setenv("TEST_VAR_2", "Gaspar")
	c.Assert(configSlice("TEST_VAR"), DeepEquals, []string{"Melchior", "Belthasar", "Gaspar"})
}

func (s *ConfigTestSuite) TestConfigSliceThreeElementsWithGap(c *C) {
	os.Setenv("TEST_VAR_0", "Melchior")
	os.Setenv("TEST_VAR_1", "Belthasar")
	os.Setenv("TEST_VAR_3", "Gaspar")
	c.Assert(configSlice("TEST_VAR"), DeepEquals, []string{"Melchior", "Belthasar"})
}
