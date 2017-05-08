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

package script

import (
	. "gopkg.in/check.v1"
)

type parserSuite struct{}

var _ = Suite(&parserSuite{})

// "$gcp.inputs.test"
func (p *parserSuite) Test_Parse_And_Eval_Env_Lookup(c *C) {
	inputsDict := LiftDict(map[string]Script{
		"version": LiftString("1.0"),
	})
	gcpDict := LiftDict(map[string]Script{
		"inputs": inputsDict,
	})
	globalsDict := LiftDict(map[string]Script{
		"gcp": gcpDict,
	})
	env := NewScriptEnvironment()
	(*env)["$"] = globalsDict

	script, err := ParseScript("$gcp.inputs.version")
	c.Assert(err, IsNil)

	result, err := EvalToGoValue(script, env)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "1.0")
}
