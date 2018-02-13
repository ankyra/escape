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
	"errors"
	"fmt"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state/validate"
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
	env, err := p.GetEnvironmentStateOrMakeNew("dev")
	c.Assert(err, IsNil)

	depl, err = env.GetOrCreateDeploymentState("archive-release")
	c.Assert(err, IsNil)
	fullDepl, err = env.GetOrCreateDeploymentState("archive-full")
	c.Assert(err, IsNil)

	dep, err := env.GetOrCreateDeploymentState("archive-release-with-deps")
	c.Assert(err, IsNil)
	deplWithDeps, err = dep.GetDeploymentOrMakeNew(DeployStage, "archive-release")
	c.Assert(err, IsNil)
	dep2, err := dep.GetDeploymentOrMakeNew(BuildStage, "build-release")
	c.Assert(err, IsNil)
	buildRootStage, err = dep2.GetDeploymentOrMakeNew(DeployStage, "build-root-release")
	c.Assert(err, IsNil)

	dep, err = env.GetOrCreateDeploymentState("archive-release-deployed-deps")
	c.Assert(err, IsNil)
	deployedDepsDepl, err = dep.GetDeploymentOrMakeNew(BuildStage, "archive-release")
	c.Assert(err, IsNil)
}

func (s *suite) Test_Deployment_NewDeploymentState(c *C) {
	d, err := NewDeploymentState(nil, "name", "project/application")
	c.Assert(err, IsNil)
	c.Assert(d.Name, Equals, "name")
	c.Assert(d.Release, Equals, "project/application")
	c.Assert(d.Stages, Not(IsNil))
	c.Assert(d.Inputs, Not(IsNil))
	c.Assert(d.environment, IsNil)
}

func (s *suite) Test_Deployment_validateAndFix_fixes_nils(c *C) {
	d, err := NewDeploymentState(nil, "name", "project/application")
	c.Assert(err, IsNil)
	d.Stages = nil
	d.Inputs = nil
	c.Assert(d.validateAndFix("name", nil), IsNil)
	c.Assert(d.Stages, Not(IsNil))
	c.Assert(d.Inputs, Not(IsNil))
}

func (s *suite) Test_Deployment_validateAndFix_fails_on_invalid_name(c *C) {
	for _, test := range validate.InvalidDeploymentNames {
		d, err := NewDeploymentState(nil, "name", "project/application")
		c.Assert(err, IsNil)
		c.Assert(d.validateAndFix(test, nil), DeepEquals, validate.InvalidDeploymentNameError(test))
	}
}

func (s *suite) Test_Deployment_validateAndFix_valid_names(c *C) {
	for _, test := range validate.ValidDeploymentNames {
		d, err := NewDeploymentState(nil, "name", "project/application")
		c.Assert(err, IsNil)
		c.Assert(d.validateAndFix(test, nil), IsNil)
		c.Assert(d.Name, Equals, test)
	}
}

func (s *suite) Test_GetRootDeploymentName(c *C) {
	c.Assert(depl.GetRootDeploymentName(), Equals, "archive-release")
	c.Assert(fullDepl.GetRootDeploymentName(), Equals, "archive-full")
	c.Assert(deplWithDeps.GetRootDeploymentName(), Equals, "archive-release-with-deps")
}

func (s *suite) Test_GetRootDeploymentStage(c *C) {
	c.Assert(depl.GetRootDeploymentStage(), Equals, "")
	c.Assert(fullDepl.GetRootDeploymentStage(), Equals, "")
	c.Assert(deplWithDeps.GetRootDeploymentStage(), Equals, DeployStage)
	c.Assert(buildRootStage.GetRootDeploymentStage(), Equals, BuildStage)
}

func (s *suite) Test_GetDeploymentPath(c *C) {
	c.Assert(depl.GetDeploymentPath(), Equals, "archive-release")
	c.Assert(fullDepl.GetDeploymentPath(), Equals, "archive-full")
	c.Assert(deplWithDeps.GetDeploymentPath(), Equals, "archive-release-with-deps:archive-release")
}

func (s *suite) Test_GetDeploymentOrMakeNew(c *C) {
	depDepl := deployedDepsDepl
	c.Assert(depDepl.Name, Equals, "archive-release")
	c.Assert(depDepl.parentStage.Name, Equals, BuildStage)

	depDepl2, err := depDepl.GetDeploymentOrMakeNew(DeployStage, "deploy-dep-name")
	c.Assert(err, IsNil)
	c.Assert(depDepl2.Name, Equals, "deploy-dep-name")
	c.Assert(depDepl2.parentStage.Name, Equals, DeployStage)

	depDepl3, err := depDepl2.GetDeploymentOrMakeNew(DeployStage, "deploy-dep-name")
	c.Assert(err, IsNil)
	c.Assert(depDepl3.Name, Equals, "deploy-dep-name")
	c.Assert(depDepl3.parentStage.Name, Equals, DeployStage)
}

func (s *suite) Test_GetDeploymentOrMakeNew_fails_on_invalid_deployment_names(c *C) {
	for _, test := range validate.InvalidDeploymentNames {
		_, err := deployedDepsDepl.GetDeploymentOrMakeNew(DeployStage, test)
		c.Assert(err, DeepEquals, validate.InvalidDeploymentNameError(test))
	}
}

func (s *suite) Test_GetPreStepInputs_for_dependency_uses_parent_build_stage(c *C) {
	inputs := deployedDepsDepl.GetPreStepInputs(DeployStage)
	c.Assert(inputs["variable"], Equals, "build_variable")
}

func (s *suite) Test_GetPreStepInputs_for_nested_dependency_uses_parent_build_stage(c *C) {
	parentDepl, err := deployedDepsDepl.GetDeploymentOrMakeNew(DeployStage, "nested1")
	c.Assert(err, IsNil)
	nestedDepl, err := parentDepl.GetDeploymentOrMakeNew(DeployStage, "nested2")
	c.Assert(err, IsNil)
	inputs := nestedDepl.GetPreStepInputs(DeployStage)
	c.Assert(inputs["variable"], Equals, "build_variable")
}

func (s *suite) Test_GetEnvironmentState(c *C) {
	env := depl.GetEnvironmentState()
	c.Assert(env.Name, Equals, "dev")
}
func (s *suite) Test_CommitVersion(c *C) {
	c.Assert(depl.GetVersion(BuildStage), Equals, "")
	c.Assert(depl.GetVersion(DeployStage), Equals, "")
	depl.CommitVersion(BuildStage, core.NewReleaseMetadata("test", "1"))
	depl.CommitVersion(DeployStage, core.NewReleaseMetadata("test", "10"))
	c.Assert(depl.GetVersion(BuildStage), Equals, "1")
	c.Assert(depl.GetVersion(DeployStage), Equals, "10")
}

func (s *suite) Test_CommitVersion_sets_provides_field(c *C) {
	metadata := core.NewReleaseMetadata("test", "1")
	metadata.SetProvides([]string{"test-provider"})
	depl.CommitVersion(DeployStage, metadata)
	c.Assert(depl.GetStageOrCreateNew(DeployStage).Provides, DeepEquals, []string{"test-provider"})
}

func (s *suite) Test_GetBuildInputs(c *C) {
	inputs := depl.GetPreStepInputs(DeployStage)
	c.Assert(inputs["input_variable"], DeepEquals, "depl_override")
	c.Assert(inputs["list_input"], DeepEquals, []interface{}{"depl_override"})
	c.Assert(inputs["env_level_variable"], DeepEquals, "env")
	c.Assert(inputs["depl_level_variable"], DeepEquals, "depl")
	c.Assert(inputs["user_level"], DeepEquals, "user")
}

func (s *suite) Test_GetProviders_nil_providers(c *C) {
	depl.GetStageOrCreateNew(DeployStage).Providers = nil
	providers := depl.GetProviders(DeployStage)
	c.Assert(providers, HasLen, 0)
}

func (s *suite) Test_GetProviders_no_providers(c *C) {
	providers := depl.GetProviders(DeployStage)
	c.Assert(providers, HasLen, 0)
}

func (s *suite) Test_GetProviders_includes_parent_providers(c *C) {
	providers := deplWithDeps.GetProviders(DeployStage)
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist")
}

func (s *suite) Test_GetProviders_includes_parent_build_providers_for_dep(c *C) {
	p, err := NewProjectStateFromFile("prj", "testdata/project.json", nil)
	c.Assert(err, IsNil)
	env, err := p.GetEnvironmentStateOrMakeNew("dev")
	c.Assert(err, IsNil)
	dep, err := env.GetOrCreateDeploymentState("archive-release-with-deps")
	c.Assert(err, IsNil)
	deplWithDeps, err = dep.GetDeploymentOrMakeNew(BuildStage, "archive-release")
	c.Assert(err, IsNil)
	providers := deplWithDeps.GetProviders(DeployStage)
	c.Assert(providers, HasLen, 3)
	c.Assert(providers["kubernetes"], Equals, "archive-release")
	c.Assert(providers["gcp"], Equals, "archive-release-build")
	c.Assert(providers["doesnt-exist"], Equals, "doesnt-exist-build")
}

func (s *suite) Test_ConfigureProviders_uses_extra_providers(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	providers := map[string]string{
		"provider1": "otherdepl",
	}
	err := deplWithDeps.ConfigureProviders(metadata, DeployStage, providers)
	c.Assert(err, IsNil)
	returnedProviders := deplWithDeps.GetProviders(DeployStage)
	c.Assert(returnedProviders["provider1"], Equals, "otherdepl")
}

func (s *suite) Test_ConfigureProviders_uses_renamed_extra_providers(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	cfg, _ := core.NewConsumerConfigFromString("provider1 as p1")
	metadata.Consumes = []*core.ConsumerConfig{cfg}
	providers := map[string]string{
		"p1": "otherdepl",
	}
	err := deplWithDeps.ConfigureProviders(metadata, DeployStage, providers)
	c.Assert(err, IsNil)
	returnedProviders := deplWithDeps.GetProviders(DeployStage)
	c.Assert(returnedProviders["p1"], Equals, "otherdepl")
}

func (s *suite) Test_ConfigureProviders_fails_if_renamed_provider_not_found(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	cfg, _ := core.NewConsumerConfigFromString("provider1 as p1")
	metadata.Consumes = []*core.ConsumerConfig{cfg}
	providers := map[string]string{}
	err := deplWithDeps.ConfigureProviders(metadata, DeployStage, providers)
	c.Assert(err, DeepEquals, errors.New("Missing provider 'p1' of type 'provider1'. This can be configured using the -p / --extra-provider flag."))
}

func (s *suite) Test_ConfigureProviders_fails_if_provider_missing(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	err := deplWithDeps.ConfigureProviders(metadata, DeployStage, nil)
	c.Assert(err, DeepEquals, fmt.Errorf("Missing provider of type 'provider1'. This can be configured using the -p / --extra-provider flag."))
}

func (s *suite) Test_ConfigureProviders_succeeds_if_provider_already_configured(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	deplWithDeps.SetProvider(DeployStage, "provider1", "otherdepl")
	err := deplWithDeps.ConfigureProviders(metadata, DeployStage, nil)
	c.Assert(err, IsNil)
	returnedProviders := deplWithDeps.GetProviders(DeployStage)
	c.Assert(returnedProviders["provider1"], Equals, "otherdepl")
}

func (s *suite) Test_Deployment_ValidateNames_fails_if_invalid_name(c *C) {
	d, err := NewDeploymentState(nil, "name", "project/application")
	c.Assert(err, IsNil)
	for _, name := range validate.InvalidDeploymentNames {
		d.Name = name
		c.Assert(d.ValidateNames(), DeepEquals, validate.InvalidDeploymentNameError(name))
	}
}

func (s *suite) Test_Deployment_ValidateNames_fails_if_invalid_sub_deployment_name(c *C) {
	for _, name := range validate.InvalidDeploymentNames {
		d, _ := NewDeploymentState(nil, "name", "project/application")
		st := d.GetStageOrCreateNew("deploy")
		brokenDepl, _ := NewDeploymentState(nil, name, "project/application")
		st.Deployments[name] = brokenDepl
		c.Assert(d.ValidateNames(), DeepEquals, validate.InvalidDeploymentNameError(name))
	}
}

func (s *suite) Test_Deployment_Summarize(c *C) {
	d := deployedDepsDepl.Summarize()
	c.Assert(d.Name, Equals, deployedDepsDepl.Name)
	c.Assert(d.Release, Equals, deployedDepsDepl.Release)
	c.Assert(d.Inputs, HasLen, 0)
	c.Assert(d.Stages, HasLen, 1)
	c.Assert(d.Stages["deploy"].Deployments, HasLen, 0)

}
