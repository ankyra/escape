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

package compiler

import (
	"fmt"
	"github.com/ankyra/escape-client/model/escape_plan"
	"github.com/ankyra/escape-client/model/registry"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Dependencies(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []string{
		"dependency-v1.0 as dep",
	}
	reg := registry.NewMockRegistry()
	reg.ReleaseMetadata = func(project, name, version string) (*core.ReleaseMetadata, error) {
		if project == "_" && name == "dependency" && version == "v1.0" {
			m := core.NewReleaseMetadata(name, "1.0")
			m.Project = "_"
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	ctx := NewCompilerContext(plan, reg, "_")
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.GetDependencies(), DeepEquals, []string{"_/dependency-v1.0"})
	c.Assert(ctx.Metadata.VariableCtx["dependency"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["dep"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dep"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
}

func (s *suite) Test_Compile_Dependencies_adds_inputs_without_defaults(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []string{
		"dependency-v1.0 as dep",
	}
	reg := registry.NewMockRegistry()
	v1, _ := variables.NewVariableFromString("no-default", "string")
	v2, _ := variables.NewVariableFromString("with-default", "string")
	v2.Default = "test"
	reg.ReleaseMetadata = func(project, name, version string) (*core.ReleaseMetadata, error) {
		if project == "_" && name == "dependency" && version == "v1.0" {
			m := core.NewReleaseMetadata(name, "1.0")
			m.Project = "_"
			m.AddInputVariable(v1)
			m.AddInputVariable(v2)
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	ctx := NewCompilerContext(plan, reg, "_")
	c.Assert(compileDependencies(ctx), IsNil)
	inputs := ctx.Metadata.GetInputs()
	c.Assert(inputs, HasLen, 1)
	c.Assert(inputs[0], DeepEquals, v1)
}

func (s *suite) Test_Compile_Dependencies_adds_consumers(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []string{
		"dependency-v1.0 as dep",
	}
	reg := registry.NewMockRegistry()
	reg.ReleaseMetadata = func(project, name, version string) (*core.ReleaseMetadata, error) {
		if project == "_" && name == "dependency" && version == "v1.0" {
			m := core.NewReleaseMetadata(name, "1.0")
			m.Project = "_"
			m.AddConsumes("test-consumer-1")
			m.AddConsumes("test-consumer-1")
			m.AddConsumes("test-consumer-2")
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	ctx := NewCompilerContext(plan, reg, "_")
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.GetConsumes(), DeepEquals, []string{"test-consumer-1", "test-consumer-2"})
}

func (s *suite) Test_Compile_Dependencies_nil(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.GetDependencies(), DeepEquals, []string{})
}

func (s *suite) Test_Compile_Dependencies_fails_if_invalid_format(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = []string{
		"$invalid_dependency$",
	}
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileDependencies(ctx).Error(), Equals, "Invalid release format: $invalid_dependency$")
}

func (s *suite) Test_Compile_Dependencies_fails_if_resolve_version_fails(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = []string{
		"dependency-latest",
	}
	reg := registry.NewMockRegistry()
	reg.ReleaseMetadata = func(project, name, version string) (*core.ReleaseMetadata, error) {
		return nil, fmt.Errorf("Resolve error")
	}
	ctx := NewCompilerContext(plan, reg, "_")
	c.Assert(compileDependencies(ctx).Error(), Equals, "Resolve error")
}
