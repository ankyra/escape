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
	"github.com/ankyra/escape/model/escape_plan"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Extensions(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error %s", dep.ReleaseId)
	}
	c.Assert(compileExtensions(ctx), IsNil)
	c.Assert(ctx.VariableCtx["_/dependency"].GetQualifiedReleaseId(), Equals, "_/dependency-v1.0")
	c.Assert(ctx.Metadata.GetExtensions(), DeepEquals, []string{"_/dependency-v1.0"})
}

func (s *suite) Test_Compile_Extensions_adds_dependencies_to_plan(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"dependency-v1.0"}

	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		if dep.ReleaseId == "_/dependency-v1.0" {
			m := core.NewReleaseMetadata("dependency", "1.0")
			m.SetDependencies([]string{"recursive-dep-v1.0 as dep", "another-dep-v1.0", "another-dep-v1.0"})
			return m, nil
		}
		return nil, fmt.Errorf("Resolve error %s", dep.ReleaseId)
	}
	c.Assert(compileExtensions(ctx), IsNil)
	deps, err := ctx.Plan.GetDependencies()
	c.Assert(err, IsNil)
	c.Assert(deps, HasLen, 2)
	cfg1 := core.NewDependencyConfig("recursive-dep-v1.0 as dep")
	cfg2 := core.NewDependencyConfig("another-dep-v1.0")
	c.Assert(cfg1.Validate(nil), IsNil)
	c.Assert(cfg2.Validate(nil), IsNil)

	// These fields will be set when the dependency configs of the escape
	// plan are validated (see dependencies.go). Not ideal.
	cfg1.Project = ""
	cfg1.Name = ""
	cfg1.Version = ""
	cfg2.Project = ""
	cfg2.Name = ""
	cfg2.Version = ""
	c.Assert(deps[0], DeepEquals, cfg1)
	c.Assert(deps[1], DeepEquals, cfg2)
}

func (s *suite) Test_Compile_Extensions_nil(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = nil
	ctx := NewCompilerContext(plan, nil)
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
		ctx := NewCompilerContext(plan, nil)
		c.Assert(compileExtensions(ctx), Not(IsNil))
	}
}

func (s *suite) Test_Compile_Extensions_fails_if_version_cant_be_resolved(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Extends = []string{"test-v1"}
	ctx := NewCompilerContext(plan, nil)
	ctx.DependencyFetcher = func(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
		return nil, fmt.Errorf("Resolve error")
	}
	c.Assert(compileExtensions(ctx), Not(IsNil))
}
