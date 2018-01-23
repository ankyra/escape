/*
Copyright 2017, 2018 Ankyra

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

package state

import (
	"testing"

	. "gopkg.in/check.v1"
)

type suite struct{}

var _ = Suite(&suite{})

func Test(t *testing.T) { TestingT(t) }

func (s *suite) Test_FromJson(c *C) {
	json := `
    {
        "name": "hello",
        "environments": {}
    }`
	p, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, IsNil)
	c.Assert(p.Name, Equals, "hello")
}

func (s *suite) Test_FromJson_fails_if_name_is_missing(c *C) {
	json := `
    {
        "environments": {}
    }`
	_, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, Not(IsNil))
}

func (s *suite) Test_From_File(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	c.Assert(p.Name, Equals, "project_name")
	env := p.GetEnvironmentStateOrMakeNew("dev")
	c.Assert(env.Inputs["input_variable"], DeepEquals, "env_override")
}

func (s *suite) Test_From_File_That_Doesnt_Exist_Returns_Empty_State(c *C) {
	_, err := NewProjectStateFromFile("prj", "asodifjaowijefowaiejfoawijefoiasjdfoiasdf.state", nil)
	c.Assert(err, IsNil)
}

func (s *suite) Test_From_File_With_Empty_File_Fails(c *C) {
	p, err := NewProjectStateFromFile("prj", "", nil)
	c.Assert(p, IsNil)
	c.Assert(err.Error(), Equals, "Configuration file path is required.")
}

func (s *suite) Test_GetEnvironmentStateOrMakeNew(c *C) {
	p, err := NewProjectState("prj")
	c.Assert(err, IsNil)
	state1 := p.GetEnvironmentStateOrMakeNew("test-env")
	state2 := p.GetEnvironmentStateOrMakeNew("test-env")
	c.Assert(state1, Not(IsNil))
	c.Assert(state2, Not(IsNil))
	c.Assert(state1, DeepEquals, state2)
}
