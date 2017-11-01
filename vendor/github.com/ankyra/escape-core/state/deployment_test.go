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
	"github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
)

var depl *DeploymentState
var deplWithDeps *DeploymentState
var fullDepl *DeploymentState
var buildRootStage *DeploymentState
var deployedDepsDepl *DeploymentState

func (s *suite) SetUpTest(c *C) {
	var err error
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")

	depl = env.GetOrCreateDeploymentState("archive-release")
	fullDepl = env.GetOrCreateDeploymentState("archive-full")

	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeploymentOrMakeNew("deploy", "archive-release")
	dep2 := dep.GetDeploymentOrMakeNew("build", "build-release")
	buildRootStage = dep2.GetDeploymentOrMakeNew("deploy", "build-root-release")

	dep = env.GetOrCreateDeploymentState("archive-release-deployed-deps")
	deployedDepsDepl = dep.GetDeploymentOrMakeNew("build", "archive-release")
}

func (s *suite) Test_GetRootDeploymentName(c *C) {
	c.Assert(depl.GetRootDeploymentName(), Equals, "archive-release")
	c.Assert(fullDepl.GetRootDeploymentName(), Equals, "archive-full")
	c.Assert(deplWithDeps.GetRootDeploymentName(), Equals, "archive-release-with-deps")
}

func (s *suite) Test_GetRootDeploymentStage(c *C) {
	c.Assert(depl.GetRootDeploymentStage(), Equals, "")
	c.Assert(fullDepl.GetRootDeploymentStage(), Equals, "")
	c.Assert(deplWithDeps.GetRootDeploymentStage(), Equals, "deploy")
	c.Assert(buildRootStage.GetRootDeploymentStage(), Equals, "build")
}

func (s *suite) Test_GetDependencyPath(c *C) {
	c.Assert(depl.GetDependencyPath(), Equals, "archive-release")
	c.Assert(fullDepl.GetDependencyPath(), Equals, "archive-full")
	c.Assert(deplWithDeps.GetDependencyPath(), Equals, "archive-release-with-deps:archive-release")
}

func (s *suite) Test_GetDeploymentOrMakeNew(c *C) {
	depDepl := deployedDepsDepl
	c.Assert(depDepl.Name, Equals, "archive-release")
	c.Assert(depDepl.parentStage.Name, Equals, "build")

	depDepl2 := depDepl.GetDeploymentOrMakeNew("deploy", "deploy-dep-name")
	c.Assert(depDepl2.Name, Equals, "deploy-dep-name")
	c.Assert(depDepl2.parentStage.Name, Equals, "deploy")

	depDepl3 := depDepl2.GetDeploymentOrMakeNew("deploy", "deploy-dep-name")
	c.Assert(depDepl3.Name, Equals, "deploy-dep-name")
	c.Assert(depDepl3.parentStage.Name, Equals, "deploy")
}

func (s *suite) Test_GetPreStepInputs_for_dependency_uses_parent_build_stage(c *C) {
	inputs := deployedDepsDepl.GetPreStepInputs("deploy")
	c.Assert(inputs["variable"], Equals, "build_variable")
}

func (s *suite) Test_GetPreStepInputs_for_nested_dependency_uses_parent_build_stage(c *C) {
	nestedDepl := deployedDepsDepl.GetDeploymentOrMakeNew("deploy", "nested1").GetDeploymentOrMakeNew("deploy", "nested2")
	inputs := nestedDepl.GetPreStepInputs("deploy")
	c.Assert(inputs["variable"], Equals, "build_variable")
}

func (s *suite) Test_GetEnvironmentState(c *C) {
	env := depl.GetEnvironmentState()
	c.Assert(env.Name, Equals, "dev")
}
func (s *suite) Test_CommitVersion(c *C) {
	c.Assert(depl.GetVersion("build"), Equals, "")
	c.Assert(depl.GetVersion("deploy"), Equals, "")
	depl.CommitVersion("build", core.NewReleaseMetadata("test", "1"))
	depl.CommitVersion("deploy", core.NewReleaseMetadata("test", "10"))
	c.Assert(depl.GetVersion("build"), Equals, "1")
	c.Assert(depl.GetVersion("deploy"), Equals, "10")
}

func (s *suite) Test_CommitVersion_sets_provides_field(c *C) {
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.SetProvides([]string{"test-provider"})
	depl.CommitVersion("deploy", metadata)
	c.Assert(depl.GetStageOrCreateNew("deploy").Provides, DeepEquals, []string{"test-provider"})
}

func (s *suite) Test_GetBuildInputs(c *C) {
	inputs := depl.GetPreStepInputs("deploy")
	c.Assert(inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
	c.Assert(inputs["env_level_variable"], DeepEquals, "env")
	c.Assert(inputs["depl_level_variable"], DeepEquals, "depl")
	c.Assert(inputs["user_level"], DeepEquals, "user")
}

func (s *suite) Test_GetProviders_nil_providers(c *C) {
	depl.GetStageOrCreateNew("deploy").Providers = nil
	providers := depl.GetProviders("deploy")
	c.Assert(providers, HasLen, 0)
}

func (s *suite) Test_GetProviders_no_providers(c *C) {
	providers := depl.GetProviders("deploy")
	c.Assert(providers, HasLen, 0)
}

func (s *suite) Test_GetProviders_includes_parent_providers(c *C) {
	providers := deplWithDeps.GetProviders("deploy")
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist")
}

func (s *suite) Test_GetProviders_includes_parent_build_providers_for_dep(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeploymentOrMakeNew("build", "archive-release")
	providers := deplWithDeps.GetProviders("deploy")
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release-build")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist-build")
}
