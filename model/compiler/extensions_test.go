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
	core "github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Extensions(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		if releaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error %s", releaseId)
	}
	c.Assert(compileExtensions(ctx), IsNil)
	c.Assert(ctx.VariableCtx["dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["dependency"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.GetExtensions(), DeepEquals, []string{"_/dependency-v1.0"})
}

func (s *suite) Test_Compile_Extensions_adds_dependencies_to_plan(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		if releaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.SetDependencies([]string{"recursive-dep-latest as dep", "another-dep-v1.0", "another-dep-v1.0"})
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error %s", releaseId)
	}
	c.Assert(compileExtensions(ctx), IsNil)
	deps, err := ctx.Plan.GetDependencies()
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, []*core.DependencyConfig{core.NewDependencyConfig("recursive-dep-latest as dep"), core.NewDependencyConfig("another-dep-v1.0")})
}

func (s *suite) Test_Compile_Extensions_adds_variable_context(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		if releaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.VariableCtx = map[string]string{
				"oh": "project/recursive-dependency-v1.0",
			}
			return m, nil
		} else if releaseId == "project/recursive-dependency-v1.0" {
			m := core.NewReleaseMetadata("recursive-dependency", "1.0")
			m.Project = "project"
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error %s", releaseId)
	}
	c.Assert(compileExtensions(ctx), IsNil)
	c.Assert(ctx.VariableCtx["dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.VariableCtx["oh"].GetQualifiedReleaseId(), Equals, "project/recursive-dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["dependency"], Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.VariableCtx["oh"], Equals, "project/recursive-dependency-v1.0")
	c.Assert(ctx.Metadata.GetExtensions(), DeepEquals, []string{"_/dependency-v1.0"})
}

func (s *suite) Test_Compile_Extensions_nil(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = nil
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileExtensions(ctx), IsNil)
	c.Assert(ctx.Metadata.GetExtensions(), DeepEquals, []string{})
}

func (s *suite) Test_Compile_Extensions_fails_if_invalid_format(c *C) {
	cases := []string{
		"adoijasodijasd oiajs doiajs doiajsd",
		"123",
		"",
		"1.0",
		"$latest",
	}
	for _, test := range cases {
		plan := escape_plan.NewEscapePlan()
		plan.Extends = []string{test}
		ctx := NewCompilerContext(plan, nil, "_")
		c.Assert(compileExtensions(ctx), Not(IsNil))
	}
}

func (s *suite) Test_Compile_Extensions_fails_if_version_cant_be_resolved(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"test-v1"}
	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileExtensions(ctx), Not(IsNil))
}

func (s *suite) Test_Compile_Extensions_fails_if_variable_context_cant_be_parsed(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}
	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		if releaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.VariableCtx = map[string]string{
				"oh": "oasdoasidja ospdij apsdojk apsodk apsodk",
			}
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileExtensions(ctx), Not(IsNil))
}

func (s *suite) Test_Compile_Extensions_fails_if_variable_context_cant_be_resolved(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil, "_")
	ctx.DependencyFetcher = func(releaseId string) (*core.ReleaseMetadata, error) {
		if releaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.VariableCtx = map[string]string{
				"oh": "_/cant-be-resovled-v1",
			}
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileExtensions(ctx), Not(IsNil))
}
