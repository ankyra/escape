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
)

type deplSuite struct{}

var _ = Suite(&deplSuite{})

var depl *DeploymentState
var deplWithDeps *DeploymentState
var fullDepl *DeploymentState
var deployedDepsDepl *DeploymentState

func (s *deplSuite) SetUpTest(c *C) {
	var err error
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl = env.GetOrCreateDeploymentState("archive-release")
	fullDepl = env.GetOrCreateDeploymentState("archive-full")
	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeployment("deploy", "archive-release")

	dep = env.GetOrCreateDeploymentState("archive-release-deployed-deps")
	deployedDepsDepl = dep.GetDeployment("build", "archive-release")
}

func (s *deplSuite) Test_GetDeployment(c *C) {
	depDepl := deployedDepsDepl
	c.Assert(depDepl.GetName(), Equals, "archive-release")
	c.Assert(depDepl.parentStage.Name, Equals, "build")

	depDepl2 := depDepl.GetDeployment("deploy", "deploy-dep-name")
	c.Assert(depDepl2.GetName(), Equals, "deploy-dep-name")
	c.Assert(depDepl2.parentStage.Name, Equals, "deploy")

	depDepl3 := depDepl2.GetDeployment("deploy", "deploy-dep-name")
	c.Assert(depDepl3.GetName(), Equals, "deploy-dep-name")
	c.Assert(depDepl3.parentStage.Name, Equals, "deploy")
}

func (s *deplSuite) Test_GetPreStepInputs_for_dependency_uses_parent_build_stage(c *C) {
	inputs := deployedDepsDepl.GetPreStepInputs("deploy")
	c.Assert(inputs["variable"], Equals, "build_variable")
}

func (s *deplSuite) Test_GetPreStepInputs_for_nested_dependency_uses_parent_build_stage(c *C) {
	nestedDepl := deployedDepsDepl.GetDeployment("deploy", "nested1").GetDeployment("deploy", "nested2")
	inputs := nestedDepl.GetPreStepInputs("deploy")
	c.Assert(inputs["variable"], Equals, "build_variable")
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
	inputs := depl.GetPreStepInputs("deploy")
	c.Assert(inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
	c.Assert(inputs["env_level_variable"], DeepEquals, "env")
	c.Assert(inputs["depl_level_variable"], DeepEquals, "depl")
	c.Assert(inputs["user_level"], DeepEquals, "user")
}

func (s *deplSuite) Test_GetProviders_nil_providers(c *C) {
	depl.getStage("deploy").Providers = nil
	providers := depl.GetProviders("deploy")
	c.Assert(providers, HasLen, 0)
}

func (s *deplSuite) Test_GetProviders_no_providers(c *C) {
	providers := depl.GetProviders("deploy")
	c.Assert(providers, HasLen, 0)
}

func (s *deplSuite) Test_GetProviders_includes_parent_providers(c *C) {
	providers := deplWithDeps.GetProviders("deploy")
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist")
}

func (s *deplSuite) Test_GetProviders_includes_parent_build_providers_for_dep(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeployment("build", "archive-release")
	providers := deplWithDeps.GetProviders("deploy")
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release-build")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist-build")
}
