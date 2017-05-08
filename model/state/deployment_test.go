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

type deplSuite struct{}

var _ = Suite(&deplSuite{})

var depl *deploymentState
var deplWithDeps *deploymentState
var fullDepl *deploymentState

func (s *deplSuite) SetUpTest(c *C) {
	p, err := NewProjectStateFromFile("testdata/project.json")
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl_, err := env.GetDeploymentState([]string{"archive-release"})
	c.Assert(err, IsNil)
	depl = depl_.(*deploymentState)

	deplWithDeps_, err := env.GetDeploymentState([]string{"archive-release-with-deps", "archive-release"})
	c.Assert(err, IsNil)
	deplWithDeps = deplWithDeps_.(*deploymentState)

	fullDepl_, err := env.GetDeploymentState([]string{"archive-full"})
	c.Assert(err, IsNil)
	fullDepl = fullDepl_.(*deploymentState)
}

func (s *deplSuite) Test_GetEnvironmentState(c *C) {
	env := depl.GetEnvironmentState()
	c.Assert(env.GetName(), Equals, "dev")
}
func (s *deplSuite) Test_SetVersion(c *C) {
	c.Assert(depl.GetVersion("build"), Equals, "")
	c.Assert(depl.GetVersion("deploy"), Equals, "")
	depl.SetVersion("build", "1")
	depl.SetVersion("deploy", "10")
	c.Assert(depl.GetVersion("build"), Equals, "1")
	c.Assert(depl.GetVersion("deploy"), Equals, "10")
}

func (s *deplSuite) Test_GetBuildInputs(c *C) {
	inputs := *depl.GetPreStepInputs("deploy")
	c.Assert(inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
	c.Assert(inputs["project_level_variable"], DeepEquals, "project")
	c.Assert(inputs["env_level_variable"], DeepEquals, "env")
	c.Assert(inputs["depl_level_variable"], DeepEquals, "depl")
	c.Assert(inputs["user_level"], DeepEquals, "user")
}

func (s *deplSuite) Test_GetProviders_no_providers(c *C) {
	providers := depl.GetProviders()
	c.Assert(providers, HasLen, 0)
}

func (s *deplSuite) Test_GetProviders_includes_parent_providers(c *C) {
	providers := deplWithDeps.GetProviders()
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist")
}

func (s *deplSuite) Test_ResolveConsumer(c *C) {
	depl, err := deplWithDeps.ResolveConsumer("kubernetes")
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
	depl, err = deplWithDeps.ResolveConsumer("gcp")
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
}

func (s *deplSuite) Test_ResolveConsumer_doesnt_exist(c *C) {
	_, err := deplWithDeps.ResolveConsumer("whatever")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "No provider of type 'whatever' was configured in the deployment state.")
}

func (s *deplSuite) Test_ResolveConsumer_provider_doesnt_exist(c *C) {
	_, err := deplWithDeps.ResolveConsumer("doesnt-exist")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment 'doesnt-exist' does not exist")
}
