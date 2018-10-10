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
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/escape_plan"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_Compile_Scripts_commands(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "docker build -t test ."
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), IsNil)
	c.Assert(ctx.Metadata.GetExecStage("deploy").Cmd, Equals, "docker")
	c.Assert(ctx.Metadata.GetExecStage("deploy").Args, DeepEquals, []string{"build", "-t", "test", "."})
}

func (s *suite) Test_Compile_Scripts_relative_script(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "testdata/script.sh"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), IsNil)
	c.Assert(ctx.Metadata.GetExecStage("deploy").RelativeScript, Equals, "testdata/script.sh")
}

func (s *suite) Test_Compile_Scripts_relative_script2(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "./testdata/script.sh"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), IsNil)
	c.Assert(ctx.Metadata.GetExecStage("deploy").RelativeScript, Equals, "./testdata/script.sh")
}

func (s *suite) Test_Compile_Scripts_relative_script_fails_if_script_doesnt_exist(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "./testdata/script_doesnt_exist.sh"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), DeepEquals, ScriptFieldError("deploy", core.ScriptDoesNotExistError("./testdata/script_doesnt_exist.sh", "./testdata/script_doesnt_exist.sh")))
}

func (s *suite) Test_Compile_Scripts_relative_script_fails_if_outside_of_basedir(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "../runners/build/testdata/test.sh"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), DeepEquals, ScriptFieldError("deploy", core.RelativeScriptOutsideOfBaseDirError("../runners/build/testdata/test.sh")))
}

func (s *suite) Test_Compile_Scripts_absolute_path(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Deploy = "/bin/ls -al"
	ctx := NewCompilerContext(plan, nil)
	c.Assert(compileScripts(ctx), IsNil)
	c.Assert(ctx.Metadata.GetExecStage("deploy").Cmd, Equals, "/bin/ls")
	c.Assert(ctx.Metadata.GetExecStage("deploy").Args, DeepEquals, []string{"-al"})
}
