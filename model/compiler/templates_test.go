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
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Templates(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Templates = []interface{}{"templates.yml.tpl"}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileTemplates(ctx), IsNil)
	c.Assert(ctx.Metadata.GetTemplates("deploy"), HasLen, 1)
	c.Assert(ctx.Metadata.GetTemplates("deploy")[0].File, Equals, "templates.yml.tpl")
	c.Assert(ctx.Metadata.GetTemplates("deploy")[0].Target, Equals, "templates.yml")
	c.Assert(ctx.Metadata.GetTemplates("build"), HasLen, 1)
	c.Assert(ctx.Metadata.GetTemplates("build")[0].File, Equals, "templates.yml.tpl")
	c.Assert(ctx.Metadata.GetTemplates("build")[0].Target, Equals, "templates.yml")
}

func (s *suite) Test_Compile_Build_Templates(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.BuildTemplates = []interface{}{"templates.yml.tpl"}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileTemplates(ctx), IsNil)
	c.Assert(ctx.Metadata.GetTemplates("deploy"), HasLen, 0)
	c.Assert(ctx.Metadata.GetTemplates("build"), HasLen, 1)
	c.Assert(ctx.Metadata.GetTemplates("build")[0].File, Equals, "templates.yml.tpl")
	c.Assert(ctx.Metadata.GetTemplates("build")[0].Target, Equals, "templates.yml")
}

func (s *suite) Test_Compile_Deploy_Templates(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.DeployTemplates = []interface{}{"templates.yml.tpl"}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileTemplates(ctx), IsNil)
	c.Assert(ctx.Metadata.GetTemplates("deploy"), HasLen, 1)
	c.Assert(ctx.Metadata.GetTemplates("deploy")[0].File, Equals, "templates.yml.tpl")
	c.Assert(ctx.Metadata.GetTemplates("deploy")[0].Target, Equals, "templates.yml")
	c.Assert(ctx.Metadata.GetTemplates("build"), HasLen, 0)
}
