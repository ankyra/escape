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

	"github.com/ankyra/escape-core/state/validate"

	. "gopkg.in/check.v1"
)

type suite struct{}

var _ = Suite(&suite{})

func Test(t *testing.T) { TestingT(t) }

func (s *suite) Test_Project_NewProjectState(c *C) {
	p, err := NewProjectState("prj")
	c.Assert(err, IsNil)
	c.Assert(p.Name, Equals, "prj")
	c.Assert(p.Environments, Not(IsNil))
	c.Assert(p.Backend, IsNil)
}

func (s *suite) Test_Project_ValidateAndFix_fixes_nil(c *C) {
	p, err := NewProjectState("prj")
	c.Assert(err, IsNil)
	p.Environments = nil
	c.Assert(p.ValidateAndFix(), IsNil)
	c.Assert(p.Environments, Not(IsNil))
}

func (s *suite) Test_Project_ValidateAndFix_fails_on_invalid_project_names(c *C) {
	for _, test := range validate.InvalidStateProjectNames {
		p, err := NewProjectState(test)
		c.Assert(err, DeepEquals, validate.InvalidProjectNameError(test))
		p.Name = test
		c.Assert(p.ValidateAndFix(), DeepEquals, validate.InvalidProjectNameError(test))
	}
}

func (s *suite) Test_Project_ValidateAndFix_valid_names(c *C) {
	for _, test := range validate.ValidStateProjectNames {
		p, err := NewProjectState(test)
		c.Assert(err, IsNil)
		p.Name = test
		c.Assert(p.ValidateAndFix(), IsNil)
		c.Assert(p.Name, Equals, test)
	}
}

func (s *suite) Test_Project_FromJson(c *C) {
	json := `
    {
        "name": "hello",
        "environments": {}
    }`
	p, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, IsNil)
	c.Assert(p.Name, Equals, "hello")
}

func (s *suite) Test_Project_FromJson_fails_if_name_is_missing(c *C) {
	json := `
    {
        "environments": {}
    }`
	_, err := NewProjectStateFromJsonString(json, nil)
	c.Assert(err, Not(IsNil))
}

func (s *suite) Test_Project_From_File(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	c.Assert(p.Name, Equals, "project_name")
	env, err := p.GetEnvironmentStateOrMakeNew("dev")
	c.Assert(err, IsNil)
	c.Assert(env.Inputs["input_variable"], DeepEquals, "env_override")
}

func (s *suite) Test_Project_From_File_That_Doesnt_Exist_Returns_Empty_State(c *C) {
	_, err := NewProjectStateFromFile("prj", "asodifjaowijefowaiejfoawijefoiasjdfoiasdf.state", nil)
	c.Assert(err, IsNil)
}

func (s *suite) Test_Project_From_File_With_Empty_File_Fails(c *C) {
	p, err := NewProjectStateFromFile("prj", "", nil)
	c.Assert(p, IsNil)
	c.Assert(err.Error(), Equals, "Configuration file path is required.")
}

func (s *suite) Test_Project_GetEnvironmentStateOrMakeNew(c *C) {
	p, err := NewProjectState("prj")
	c.Assert(err, IsNil)
	state1, err := p.GetEnvironmentStateOrMakeNew("test-env")
	c.Assert(err, IsNil)
	state2, err := p.GetEnvironmentStateOrMakeNew("test-env")
	c.Assert(err, IsNil)
	c.Assert(state1, Not(IsNil))
	c.Assert(state2, Not(IsNil))
	c.Assert(state1, DeepEquals, state2)
}

func (s *suite) Test_Project_GetEnvironmentStateOrMakeNew_fails_on_invalid_env_name(c *C) {
	p, err := NewProjectState("prj")
	c.Assert(err, IsNil)
	_, err = p.GetEnvironmentStateOrMakeNew("$")
	c.Assert(err, DeepEquals, validate.InvalidEnvironmentNameError("$"))
}
