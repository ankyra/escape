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
	"github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
)

type envSuite struct{}

var _ = Suite(&envSuite{})

func (s *envSuite) Test_Name_Field_Is_Set(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("incomplete_env")
	c.Assert(env.GetName(), Equals, "incomplete_env")
}

func (s *envSuite) Test_LookupDeploymentState(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl, err := env.LookupDeploymentState("archive-release")
	c.Assert(err, IsNil)
	c.Assert(depl.GetName(), Equals, "archive-release")
	c.Assert(depl.Inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(depl.Inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *envSuite) Test_LookupDeploymentState_doesnt_exist(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	_, err = env.LookupDeploymentState("doesnt-exist")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment 'doesnt-exist' does not exist")
}

func (s *envSuite) Test_GetOrCreateDeploymentState_no_deps(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("archive-release")
	c.Assert(depl.GetName(), Equals, "archive-release")
	c.Assert(depl.Inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(depl.Inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
}

func (s *envSuite) Test_GetOrCreateDeploymentState_doesnt_exist_no_deps_returns_new(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("doesnt-exist")
	c.Assert(depl.GetName(), Equals, "doesnt-exist")
	c.Assert(depl.Inputs, HasLen, 0)
}

func (s *envSuite) Test_GetProviders(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("provider")
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.Provides = []string{"test-provider"}
	depl.CommitVersion("deploy", metadata)
	providers := env.GetProviders()
	c.Assert(providers, HasLen, 1)
	c.Assert(providers["test-provider"], DeepEquals, []string{"provider"})
}

func (s *envSuite) Test_GetProvidersOfType(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env := p.GetEnvironmentStateOrMakeNew("dev")
	depl := env.GetOrCreateDeploymentState("provider")
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.Provides = []string{"test-provider"}
	depl.CommitVersion("deploy", metadata)
	providers := env.GetProvidersOfType("test-provider")
	c.Assert(providers, HasLen, 1)
	c.Assert(providers, DeepEquals, []string{"provider"})

	providers = env.GetProvidersOfType("no-test-provider")
	c.Assert(providers, HasLen, 0)
}
