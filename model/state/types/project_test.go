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

package types

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

type projectSuite struct{}

var _ = Suite(&projectSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *projectSuite) Test_FromJson(c *C) {
	json := `
    {
        "name": "hello",
        "environments": {}
    }`
	p, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, IsNil)
	c.Assert(p.GetName(), Equals, "hello")
}

func (s *projectSuite) Test_FromJson_fails_if_name_is_missing(c *C) {
	json := `
    {
        "environments": {}
    }`
	_, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, Not(IsNil))
}

func (s *projectSuite) Test_From_File(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	c.Assert(p.GetName(), Equals, "project_name")
	env := p.GetEnvironmentStateOrMakeNew("dev")
	c.Assert(env.getInputs()["input_variable"], DeepEquals, "env_override")
}

func (s *projectSuite) Test_From_File_That_Doesnt_Exist_Returns_Empty_State(c *C) {
	_, err := NewProjectStateFromFile("prj", "asodifjaowijefowaiejfoawijefoiasjdfoiasdf.state", nil)
	c.Assert(err, IsNil)
}

func (s *projectSuite) Test_Save_Non_Existing_File(c *C) {
	os.RemoveAll("testdata/doesnt_exist.state")
	p, err := NewProjectStateFromFile("prj", "testdata/doesnt_exist.state", nil)
	c.Assert(err, IsNil)
	c.Assert(p.GetName(), Not(Equals), "overwrite")
	p.SetName("overwrite")
	err = p.Save()
	c.Assert(err, IsNil)
	p, err = NewProjectStateFromFile("prj", "testdata/doesnt_exist.state", nil)
	c.Assert(err, IsNil)
	c.Assert(p.GetName(), Equals, "overwrite")
	os.RemoveAll("testdata/doesnt_exist.state")
}

func (s *projectSuite) Test_From_File_With_Empty_File_Fails(c *C) {
	p, err := NewProjectStateFromFile("prj", "", nil)
	c.Assert(p, IsNil)
	c.Assert(err.Error(), Equals, "Configuration file path is required.")
}

func (s *projectSuite) Test_GetEnvironmentStateOrMakeNew(c *C) {
	p, err := newProjectState("prj")
	c.Assert(err, IsNil)
	state1 := p.GetEnvironmentStateOrMakeNew("test-env")
	state2 := p.GetEnvironmentStateOrMakeNew("test-env")
	c.Assert(state1, Not(IsNil))
	c.Assert(state2, Not(IsNil))
	c.Assert(state1, DeepEquals, state2)
}
