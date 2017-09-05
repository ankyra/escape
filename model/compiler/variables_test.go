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
	"github.com/ankyra/escape-client/model/escape_plan"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Inputs(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Inputs = []interface{}{"input1"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileInputs(ctx), IsNil)
	c.Assert(ctx.Metadata.GetInputs("deploy"), HasLen, 1)
	c.Assert(ctx.Metadata.GetInputs("deploy")[0].Id, Equals, "input1")
	c.Assert(ctx.Metadata.GetInputs("build"), HasLen, 1)
	c.Assert(ctx.Metadata.GetInputs("build")[0].Id, Equals, "input1")
}

func (s *suite) Test_Compile_Build_Inputs(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.BuildInputs = []interface{}{"input1"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileInputs(ctx), IsNil)
	c.Assert(ctx.Metadata.GetInputs("deploy"), HasLen, 0)
	c.Assert(ctx.Metadata.GetInputs("build"), HasLen, 1)
	c.Assert(ctx.Metadata.GetInputs("build")[0].Id, Equals, "input1")
}

func (s *suite) Test_Compile_Deploy_Inputs(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.DeployInputs = []interface{}{"input1"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileInputs(ctx), IsNil)
	c.Assert(ctx.Metadata.GetInputs("deploy"), HasLen, 1)
	c.Assert(ctx.Metadata.GetInputs("deploy")[0].Id, Equals, "input1")
	c.Assert(ctx.Metadata.GetInputs("build"), HasLen, 0)
}
