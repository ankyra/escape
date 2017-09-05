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

func (s *suite) Test_Compile_Consumes(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Consumes = []interface{}{"consumer1", "consumer2"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileConsumers(ctx), IsNil)
	c.Assert(ctx.Metadata.GetConsumes("deploy"), DeepEquals, []string{"consumer1", "consumer2"})
	c.Assert(ctx.Metadata.GetConsumes("build"), DeepEquals, []string{"consumer1", "consumer2"})
}

func (s *suite) Test_Compile_Build_Consumes(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.BuildConsumes = []interface{}{"consumer1", "consumer2"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileConsumers(ctx), IsNil)
	c.Assert(ctx.Metadata.GetConsumes("deploy"), DeepEquals, []string{})
	c.Assert(ctx.Metadata.GetConsumes("build"), DeepEquals, []string{"consumer1", "consumer2"})
}

func (s *suite) Test_Compile_Deploy_Consumes(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.DeployConsumes = []interface{}{"consumer1", "consumer2"}
	ctx := NewCompilerContext(plan, nil, "my-project")
	c.Assert(compileConsumers(ctx), IsNil)
	c.Assert(ctx.Metadata.GetConsumes("deploy"), DeepEquals, []string{"consumer1", "consumer2"})
	c.Assert(ctx.Metadata.GetConsumes("build"), DeepEquals, []string{})
}
