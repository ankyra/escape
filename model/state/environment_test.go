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

package state

import (
	. "gopkg.in/check.v1"
)

type envSuite struct{}

var _ = Suite(&envSuite{})

func (s *envSuite) Test_Name_Field_Is_Set(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("incomplete_env")
	c.Assert(env.GetName(), Equals, "incomplete_env")
}

func (s *envSuite) Test_LookupDeploymentState(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.LookupDeploymentState("archive-release")
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
	c.Assert((*depl.Inputs)["input_variable"], DeepEquals, "depl_override")
	c.Assert((*depl.Inputs)["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *envSuite) Test_LookupDeploymentState_doesnt_exist(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	_, err = env.LookupDeploymentState("doesnt-exist")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment 'doesnt-exist' does not exist")
}

func (s *envSuite) Test_GetDeploymentState_missing_args(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	_, err = env.GetDeploymentState([]string{})
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing name to resolve deployment state. This is a bug in Escape.")

	_, err = env.GetDeploymentState(nil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing name to resolve deployment state. This is a bug in Escape.")
}

func (s *envSuite) Test_GetDeploymentState_no_deps(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"archive-release"})
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
	c.Assert((*depl.Inputs)["input_variable"], DeepEquals, "depl_override")
	c.Assert((*depl.Inputs)["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *envSuite) Test_GetDeploymentState_doesnt_exist_no_deps_returns_new(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"doesnt-exist"})
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "doesnt-exist")
	c.Assert((*depl.Inputs), HasLen, 0)
}

func (s *envSuite) Test_GetDeploymentState_with_deps(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"archive-release-with-deps", "archive-release"})
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
	c.Assert((*depl.Inputs)["input_variable"], DeepEquals, "depl_override2")
	c.Assert((*depl.Inputs)["list_input"], DeepEquals, []interface{}{"depl_override2"})
}

func (s *envSuite) Test_GetDeploymentState_doesnt_exist_with_deps_returns_new(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"archive-release-with-deps", "doesnt-exist"})
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "doesnt-exist")
	c.Assert((*depl.Inputs), HasLen, 0)
}

func (s *envSuite) Test_GetDeploymentState_doesnt_exist_with_non_existing_roots_returns_new(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"doesnt-exist1", "doesnt-exist2"})
	depl := depl_.(*deploymentState)
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "doesnt-exist2")
	c.Assert((*depl.Inputs), HasLen, 0)
}
