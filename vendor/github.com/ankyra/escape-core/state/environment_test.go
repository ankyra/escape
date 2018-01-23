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
	"github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Name_Field_Is_Set(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("incomplete_env")
	c.Assert(env.Name, Equals, "incomplete_env")
}

func (s *suite) Test_LookupDeploymentState(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl, err := env.LookupDeploymentState("archive-release")
	c.Assert(err, IsNil)
	c.Assert(depl.Name, Equals, "archive-release")
	c.Assert(depl.Inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(depl.Inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *suite) Test_LookupDeploymentState_doesnt_exist(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	_, err = env.LookupDeploymentState("doesnt-exist")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment 'doesnt-exist' does not exist")
}

func (s *suite) Test_GetOrCreateDeploymentState_no_deps(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("archive-release")
	c.Assert(depl.Name, Equals, "archive-release")
	c.Assert(depl.Inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(depl.Inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *suite) Test_GetOrCreateDeploymentState_doesnt_exist_no_deps_returns_new(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("doesnt-exist")
	c.Assert(depl.Name, Equals, "doesnt-exist")
	c.Assert(depl.Inputs, HasLen, 0)
}

func (s *suite) Test_GetProviders(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("provider")
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.SetProvides([]string{"test-provider"})
	depl.CommitVersion("deploy", metadata)
	providers := env.GetProviders()
	c.Assert(providers, HasLen, 1)
	c.Assert(providers["test-provider"], DeepEquals, []string{"provider"})
}

func (s *suite) Test_GetProvidersOfType(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("provider")
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.SetProvides([]string{"test-provider"})
	depl.CommitVersion("deploy", metadata)
	providers := env.GetProvidersOfType("test-provider")
	c.Assert(providers, HasLen, 1)
	c.Assert(providers, DeepEquals, []string{"provider"})

	providers = env.GetProvidersOfType("no-test-provider")
	c.Assert(providers, HasLen, 0)
}

func (s *suite) Test_ResolveDeploymentPath(c *C) {
	proj, _ := NewProjectState("project")
	env := proj.GetEnvironmentStateOrMakeNew("env")

	_, err := env.ResolveDeploymentPath("deploy", "test")
	c.Assert(err, DeepEquals, DeploymentDoesNotExistError("test"))
	_, err = env.ResolveDeploymentPath("build", "test")
	c.Assert(err, DeepEquals, DeploymentDoesNotExistError("test"))

	depl := env.GetOrCreateDeploymentState("test")
	returnedDepl, err := env.ResolveDeploymentPath("deploy", "test")
	c.Assert(err, IsNil)
	c.Assert(returnedDepl, DeepEquals, depl)

	deplDep := depl.GetDeploymentOrMakeNew("deploy", "test-dependency")
	returnedDepl, err = env.ResolveDeploymentPath("deploy", "test:test-dependency")
	c.Assert(err, IsNil)
	c.Assert(returnedDepl, DeepEquals, deplDep)
	_, err = env.ResolveDeploymentPath("build", "test:test-dependency")
	c.Assert(err, DeepEquals, DeploymentPathResolveError("build", "test:test-dependency", "test-dependency"))

	deplDep2 := deplDep.GetDeploymentOrMakeNew("deploy", "test-dependency2")
	returnedDepl, err = env.ResolveDeploymentPath("deploy", "test:test-dependency:test-dependency2")
	c.Assert(err, IsNil)
	c.Assert(returnedDepl, DeepEquals, deplDep2)
}

func (s *suite) Test_ResolveDeploymentPath_with_build_stage(c *C) {
	proj, _ := NewProjectState("project")
	env := proj.GetEnvironmentStateOrMakeNew("env")

	depl := env.GetOrCreateDeploymentState("test")
	returnedDepl, err := env.ResolveDeploymentPath("build", "test")
	c.Assert(err, IsNil)
	c.Assert(returnedDepl, DeepEquals, depl)

	deplDep := depl.GetDeploymentOrMakeNew("build", "test-dependency")
	returnedDepl, err = env.ResolveDeploymentPath("build", "test:test-dependency")
	c.Assert(err, IsNil)
	c.Assert(returnedDepl, DeepEquals, deplDep)
	_, err = env.ResolveDeploymentPath("deploy", "test:test-dependency")
	c.Assert(err, DeepEquals, DeploymentPathResolveError("deploy", "test:test-dependency", "test-dependency"))
}
