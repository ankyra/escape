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

	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/model/inventory"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Version_resolve_version(c *C) {
	versions := map[string]string{
		"@":     "0",
		"1.@":   "1.0",
		"1.0.@": "1.0.0",
		"auto":  "0",
	}
	var capturedProject, capturedName *string
	inventory := inventory.NewMockInventory()
	inventory.NextVersion = func(project, name, versionPrefix string) (string, error) {
		capturedProject = &project
		capturedName = &name
		return versionPrefix + "0", nil
	}
	for version, expected := range versions {
		plan := escape_plan.NewEscapePlan()
		plan.Name = "my-build"
		plan.Version = version
		ctx := NewCompilerContext(plan, inventory)
		ctx.Metadata.Name = "my-build"
		ctx.Metadata.Project = "cheeky-project"
		c.Assert(compileVersion(ctx), IsNil)
		c.Assert(ctx.Metadata.Version, Equals, expected)
		c.Assert(*capturedProject, Equals, "cheeky-project")
		c.Assert(*capturedName, Equals, "my-build")
	}
}

func (s *suite) Test_Compile_Version_no_resolve_needed(c *C) {
	versions := []string{
		"1",
		"1.0",
		"1.0.0",
	}
	for _, version := range versions {
		plan := escape_plan.NewEscapePlan()
		plan.Version = version
		ctx := NewCompilerContext(plan, nil)
		c.Assert(compileVersion(ctx), IsNil)
		c.Assert(ctx.Metadata.Version, Equals, version)
	}
}

func (s *suite) Test_Compile_Version_fails_if_resolve_fails(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Version = "1.@"
	inventory := inventory.NewMockInventory()
	inventory.NextVersion = func(project, name, versionPrefix string) (string, error) {
		return "", fmt.Errorf("Resolve error")
	}
	ctx := NewCompilerContext(plan, inventory)
	c.Assert(compileVersion(ctx).Error(), Equals, "Resolve error")
}

func (s *suite) Test_Compile_Version_fails_if_version_expression_cant_be_parsed(c *C) {
	versions := []string{
		"$.$.$.",
		"$waosdijasoidjasoidj",
		"0asd",
		"v1.0.0.102",
		"",
	}
	for _, version := range versions {
		plan := escape_plan.NewEscapePlan()
		plan.Version = version
		ctx := NewCompilerContext(plan, nil)
		c.Assert(compileVersion(ctx).Error(), Not(IsNil))
	}
}
