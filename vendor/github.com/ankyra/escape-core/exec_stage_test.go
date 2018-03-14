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
	. "gopkg.in/check.v1"
)

type execSuite struct{} // lol

var _ = Suite(&execSuite{})

func (s *execSuite) Test_ExecStage_ValidateAndFix_parses_deprecated_Script(c *C) {
	cases := [][]interface{}{
		[]interface{}{"myscript.sh", []string{"/bin/sh", "-c", "./myscript.sh .escape/outputs.json"}},
		[]interface{}{"myscript.sh test", []string{"/bin/sh", "-c", "./myscript.sh test .escape/outputs.json"}},
		[]interface{}{"deps/_/escape/escape", []string{"/bin/sh", "-c", "./deps/_/escape/escape .escape/outputs.json"}},
	}
	for _, test := range cases {
		unit := ExecStage{
			RelativeScript: test[0].(string),
		}
		c.Assert(unit.ValidateAndFix(), IsNil)
		c.Assert(unit.GetAsCommand(), DeepEquals, test[1])
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
