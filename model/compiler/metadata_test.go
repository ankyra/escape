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
	core "github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Metadata(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Metadata = map[string]string{
		"test":  "normal field",
		"test2": "$$escaped field",
		"test3": "$dep.version",
	}
	ctx := NewCompilerContext(plan, nil)
	ctx.VariableCtx = map[string]*core.ReleaseMetadata{
		"dep": core.NewReleaseMetadata("test", "1.0"),
	}
	c.Assert(compileMetadata(ctx), IsNil)
	c.Assert(ctx.Metadata.Metadata["test"], Equals, "normal field")
	c.Assert(ctx.Metadata.Metadata["test2"], Equals, "$$escaped field")
	c.Assert(ctx.Metadata.Metadata["test3"], Equals, "1.0")
}

func (s *suite) Test_Compile_Metadata_nil(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Metadata = nil
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileMetadata(ctx), IsNil)
	c.Assert(ctx.Metadata.Metadata, DeepEquals, map[string]string{})
}

func (s *suite) Test_Compile_Metadata_fails_if_field_cant_be_evaluated(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Metadata = map[string]string{
		"test": "$.$.$.$uhoh",
	}
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileMetadata(ctx), Not(IsNil))
	c.Assert(ctx.Metadata.Metadata["test"], Equals, "")
}
