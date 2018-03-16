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

package core

import (
	"github.com/ankyra/escape-core/script"
	. "gopkg.in/check.v1"
)

type execSuite struct{} // lol

var _ = Suite(&execSuite{})

func (s *execSuite) Test_ExecStage_ValidateAndFix_parses_deprecated_Script(c *C) {
	cases := [][]interface{}{
		[]interface{}{"myscript.sh", []string{"sh", "-c", "./myscript.sh .escape/outputs.json"}},
		[]interface{}{"myscript.sh test", []string{"sh", "-c", "./myscript.sh test .escape/outputs.json"}},
		[]interface{}{"./myscript.sh test", []string{"sh", "-c", "./myscript.sh test .escape/outputs.json"}},
		[]interface{}{"/test/myscript.sh test", []string{"sh", "-c", "/test/myscript.sh test .escape/outputs.json"}},
	}
	for _, test := range cases {
		unit := ExecStage{
			RelativeScript: test[0].(string),
		}
		c.Assert(unit.ValidateAndFix(), IsNil)
		cmd, err := unit.GetAsCommand()
		c.Assert(err, IsNil)
		c.Assert(cmd, DeepEquals, test[1])
	}
}

func (s *execSuite) Test_ExecStage_ValidateAndFix_errors_when_both_cmd_and_inline_are_given(c *C) {
	unit := ExecStage{
		Inline: "test",
		Cmd:    "test",
	}
	c.Assert(unit.ValidateAndFix(), Not(IsNil))
}

func (s *execSuite) Test_ExecStage_from_dict(c *C) {
	unit, err := NewExecStageFromDict(map[interface{}]interface{}{
		"script": "test.sh",
		"cmd":    "docker",
		"args":   []interface{}{"clean"},
		"inline": "inline",
	})
	c.Assert(err, IsNil)
	c.Assert(unit.RelativeScript, Equals, "test.sh")
	c.Assert(unit.Inline, Equals, "inline")
	c.Assert(unit.Cmd, Equals, "docker")
	c.Assert(unit.Args, DeepEquals, []string{"clean"})
}

func (s *execSuite) Test_ExecStage_Eval_no_script_used(c *C) {
	globals := map[string]script.Script{
		"$": script.LiftDict(map[string]script.Script{
			"test": script.LiftString("testing"),
		}),
	}
	unit, err := NewExecStageFromDict(map[interface{}]interface{}{
		"script": "test.sh",
		"inline": "echo hallo\necho hallo",
	})
	c.Assert(err, IsNil)
	env := script.NewScriptEnvironmentFromMap(globals)
	newUnit, err := unit.Eval(env)
	c.Assert(err, IsNil)
	c.Assert(newUnit.RelativeScript, Equals, "test.sh")
	c.Assert(newUnit.Inline, Equals, "echo hallo\necho hallo")
}

func (s *execSuite) Test_ExecStage_Eval_all_fields(c *C) {
	globals := map[string]script.Script{
		"$": script.LiftDict(map[string]script.Script{
			"test": script.LiftString("testing"),
		}),
	}
	unit, err := NewExecStageFromDict(map[interface{}]interface{}{
		"script": "$test",
		"cmd":    "$test",
		"args":   []interface{}{"$test", "123", "$test"},
		"inline": "$test",
	})
	c.Assert(err, IsNil)
	env := script.NewScriptEnvironmentFromMap(globals)
	newUnit, err := unit.Eval(env)
	c.Assert(err, IsNil)
	c.Assert(newUnit.RelativeScript, Equals, "testing")
	c.Assert(newUnit.Cmd, Equals, "testing")
	c.Assert(newUnit.Inline, Equals, "$test") // Inline uses shell already
	c.Assert(newUnit.Args, DeepEquals, []string{"testing", "123", "testing"})
}

func (s *execSuite) Test_ExecStage_String(c *C) {
	unit := &ExecStage{RelativeScript: "script.sh"}
	c.Assert(unit.String(), Equals, "script.sh")
	unit = &ExecStage{Cmd: "script.sh"}
	c.Assert(unit.String(), Equals, "script.sh ")
	unit = &ExecStage{Cmd: "script.sh", Args: []string{}}
	c.Assert(unit.String(), Equals, "script.sh ")
	unit = &ExecStage{Cmd: "script.sh", Args: []string{"test"}}
	c.Assert(unit.String(), Equals, "script.sh test")
	unit = &ExecStage{Cmd: "script.sh", Args: []string{"test", "test2"}}
	c.Assert(unit.String(), Equals, "script.sh test test2")
	unit = &ExecStage{}
	c.Assert(unit.String(), Equals, "<inline script starting with ''>")
	unit = &ExecStage{Inline: "script.sh"}
	c.Assert(unit.String(), Equals, "<inline script starting with 'script.sh'>")
	unit = &ExecStage{Inline: "script.sh\nscriptasdasdasdasd"}
	c.Assert(unit.String(), Equals, "<inline script starting with 'script.sh'>")
}
