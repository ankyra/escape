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
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/util"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Basics(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Name = "testor"
	plan.Description = "  trim me\nplease\n"
	plan.Provides = []string{"provider1", "provider2"}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileBasicFields(ctx), IsNil)
	c.Assert(ctx.Metadata.Name, Equals, "testor")
	c.Assert(ctx.Metadata.Description, Equals, "trim me\nplease")
	c.Assert(ctx.Metadata.Project, Equals, "_")
	c.Assert(ctx.Metadata.GetProvides(), DeepEquals, []string{"provider1", "provider2"})
	c.Assert(ctx.Metadata.BuiltWithEscapeVersion, Equals, util.EscapeVersion)
}

func (s *suite) Test_Compile_Basics_fails_if_name_is_not_set(c *C) {
	plan := escape_plan.NewEscapePlan()
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileBasicFields(ctx).Error(), Equals, "Missing build name. Add a 'name' field to your Escape plan")
}

func (s *suite) Test_Compile_Basics_set_default_project(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Name = "testor"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileBasicFields(ctx), IsNil)
	c.Assert(ctx.Metadata.Name, Equals, "testor")
	c.Assert(ctx.Metadata.Project, Equals, "_")
}

func (s *suite) Test_Compile_Basics_parse_project(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Name = "project/testor"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileBasicFields(ctx), IsNil)
	c.Assert(ctx.Metadata.Name, Equals, "testor")
	c.Assert(ctx.Metadata.Project, Equals, "project")
}

func (s *suite) Test_Compile_Basics_parse_project_fails(c *C) {
	testCases := []string{
		"project/",
		"/",
	}
	for _, test := range testCases {
		plan := escape_plan.NewEscapePlan()
		plan.Name = test
		ctx := NewCompilerContext(plan, nil)
		c.Assert(compileBasicFields(ctx), Not(IsNil))
	}
}
