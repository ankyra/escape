/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package escape_plan

import (
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"testing"
)

type planSuite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&planSuite{})

func (s *planSuite) Test_LoadConfig_fails_if_not_exists(c *C) {
	unit := NewEscapePlan()
	err := unit.LoadConfig("testdata/doesnt_exist.yml")
	c.Assert(err.Error(), Equals, "Escape plan 'testdata/doesnt_exist.yml' was not found. Use 'escape plan init' to create it")
}

func (s *planSuite) Test_LoadConfig_fails_if_read_fails(c *C) {
	unit := NewEscapePlan()
	os.Chmod("testdata/cant_read.yml", 0)
	err := unit.LoadConfig("testdata/cant_read.yml")
	c.Assert(err.Error(), Equals, "Couldn't read Escape plan 'testdata/cant_read.yml': permission denied")
	os.Chmod("testdata/cant_read.yml", 0666)
}

func (s *planSuite) Test_LoadConfig_fails_if_invalid_yaml(c *C) {
	unit := NewEscapePlan()
	err := unit.LoadConfig("testdata/cant_parse.yml")
	errorString := err.Error()
	expectError := "Couldn't parse Escape plan 'testdata/cant_parse.yml': yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `$invalid` into escape_plan.EscapePlan"
	c.Assert(errorString, Equals, expectError)
}

func (s *planSuite) Test_Init(c *C) {
	unit := NewEscapePlan()
	unit.Init("test")
	c.Assert(unit.GetName(), Equals, "test")
	c.Assert(unit.GetVersion(), Equals, "@")
}

func (s *planSuite) Test_GetReleaseId(c *C) {
	unit := NewEscapePlan()
	err := unit.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	c.Assert(unit.GetReleaseId(), Equals, "test-v0.1.@")
}

func (s *planSuite) Test_GetVersionlessReleaseId(c *C) {
	unit := NewEscapePlan()
	err := unit.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	c.Assert(unit.GetVersionlessReleaseId(), Equals, "test")
}

func (s *planSuite) Test_PrettyPrint(c *C) {
	unit := NewEscapePlan()
	err := unit.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	pretty := unit.ToYaml()
	expected, err := ioutil.ReadFile("testdata/fixture_pretty_print.yml")
	c.Assert(err, IsNil)
	c.Assert(string(pretty), Equals, string(expected))
}
