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

package compiler

import (
	"fmt"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/variables"
	"github.com/ankyra/escape/model/escape_plan"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Dependencies(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []interface{}{
		"dependency-latest as dep",
	}
	lookupResult := core.NewReleaseMetadata("dependency", "1.0")
	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			return lookupResult, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	ctx.ReleaseQuery = func(dep *core.Dependency) (*core.ReleaseMetadata, error) {
		if dep.GetQualifiedReleaseId() == "_/dependency-latest" {
			return lookupResult, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.Depends[0], DeepEquals, core.NewDependencyConfig("_/dependency-v1.0"))
	c.Assert(ctx.Metadata.VariableCtx["dependency"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["dep"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dep"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
}

func (s *suite) Test_Compile_Dependencies_with_mapping(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []interface{}{
		map[interface{}]interface{}{
			"release_id": "dependency-latest as dep",
			"mapping": map[interface{}]interface{}{
				"input_variable": "test",
			},
		},
	}
	lookupResult := core.NewReleaseMetadata("dependency", "1.0")
	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			return lookupResult, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	ctx.ReleaseQuery = func(dep *core.Dependency) (*core.ReleaseMetadata, error) {
		if dep.GetQualifiedReleaseId() == "_/dependency-latest" {
			return lookupResult, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.Depends[0].ReleaseId, Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.Depends[0].BuildMapping["input_variable"], Equals, "test")
	c.Assert(ctx.Metadata.Depends[0].DeployMapping["input_variable"], Equals, "test")
	c.Assert(ctx.Metadata.VariableCtx["dependency"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["dep"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["dep"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
}

func (s *suite) Test_Compile_Dependencies_adds_inputs_without_defaults(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []interface{}{
		"dependency-v1.0",
	}
	v1, _ := variables.NewVariableFromString("no-default", "string")
	v2, _ := variables.NewVariableFromString("with-default", "string")
	v2.Default = "test"
	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.AddInputVariable(v1)
			m.AddInputVariable(v2)
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileDependencies(ctx), IsNil)
	inputs := ctx.Metadata.Inputs
	c.Assert(inputs, HasLen, 1)
	c.Assert(inputs[0], DeepEquals, v1)
}

func (s *suite) Test_Compile_Dependencies_adds_consumers(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	plan.Depends = []interface{}{
		"dependency-v1.0",
	}
	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.AddConsumes(core.NewConsumerConfig("test-consumer-1"))
			m.AddConsumes(core.NewConsumerConfig("test-consumer-1"))
			consumer := core.NewConsumerConfig("test-consumer-2")
			consumer.Scopes = []string{"deploy"}
			m.AddConsumes(consumer)
			consumer = core.NewConsumerConfig("test-consumer-3")
			consumer.Scopes = []string{"build"}
			m.AddConsumes(consumer)
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.GetConsumes("deploy"), DeepEquals, []string{"test-consumer-1", "test-consumer-2"})
	c.Assert(ctx.Metadata.GetConsumes("build"), DeepEquals, []string{"test-consumer-1", "test-consumer-3"})
}

func (s *suite) Test_Compile_Dependencies_nil(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = nil
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileDependencies(ctx), IsNil)
	c.Assert(ctx.Metadata.Depends, HasLen, 0)
}

func (s *suite) Test_Compile_Dependencies_fails_if_invalid_format(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = []interface{}{
		"$invalid_dependency$",
	}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileDependencies(ctx).Error(), Equals, "Invalid release format: $invalid_dependency$")
}

func (s *suite) Test_Compile_Dependencies_fails_if_resolve_version_fails(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = []interface{}{
		"dependency-latest",
	}
	ctx := NewCompilerContext(plan, nil)
	ctx.ReleaseQuery = func(dep *core.Dependency) (*core.ReleaseMetadata, error) {
		return nil, fmt.Errorf("Resolve error")
	}
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileDependencies(ctx).Error(), Equals, "Resolve error")
}

func (s *suite) Test_Compile_Dependencies_fails_if_unexpected_type(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Depends = []interface{}{
		5,
	}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileDependencies(ctx).Error(), Equals, "Invalid dependency format '5' (expecting dict or string, got 'int')")
}
