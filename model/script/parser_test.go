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

func (p *parserSuite) Test_parseString(c *C) {
	result := parseString(`"test"  `)
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "  ")
	c.Assert(IsStringAtom(result.Result), Equals, true)
	c.Assert(ExpectStringAtom(result.Result), Equals, "test")
}

func (p *parserSuite) Test_parseString_escaping(c *C) {
	result := parseString(`"test\"test\n\t\\"`)
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")
	c.Assert(IsStringAtom(result.Result), Equals, true)
	c.Assert(ExpectStringAtom(result.Result), Equals, "test\"test\n\t\\")
}

func (p *parserSuite) Test_parseString_escaping_on_unknown_characters(c *C) {
	result := parseString(`"test\atest"`)
	c.Assert(result.Error, Not(IsNil))
}

func (p *parserSuite) Test_ParseScript_escaped_string(c *C) {
	result, err := ParseScript("$$escaped")
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "$$escaped")
}

func (p *parserSuite) Test_ParseScript_env_lookup(c *C) {
	result, err := ParseScript("$gcp")
	c.Assert(err, IsNil)
	atom := hasStringArgument(c, result, "gcp")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseScript_env_lookup_application_full(c *C) {
	result, err := ParseScript("$gcp.outputs.key_ident")
	c.Assert(err, IsNil)
	atom := hasStringArgument(c, result, "key_ident")
	atom = hasStringArgument(c, atom.To, "outputs")
	atom = hasStringArgument(c, atom.To, "gcp")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func hasStringArgument(c *C, s Script, key string) *apply {
	c.Assert(IsApplyAtom(s), Equals, true)
	atom := ExpectApplyAtom(s)
	c.Assert(atom.Arguments, HasLen, 1)
	c.Assert(ExpectStringAtom(atom.Arguments[0]), Equals, key)
	return atom
}

func (p *parserSuite) Test_ParseScript_env_func_call(c *C) {
	result, err := ParseScript("$__id( $test, $test2 )")
	c.Assert(err, IsNil)
	c.Assert(IsApplyAtom(result), Equals, true)

	atom := ExpectApplyAtom(result)
	c.Assert(IsApplyAtom(atom.To), Equals, true)
	c.Assert(atom.Arguments, HasLen, 2)

	arg := hasStringArgument(c, atom.Arguments[0], "test")
	arg = hasStringArgument(c, arg.To, "$")
	c.Assert(IsFunctionAtom(arg.To), Equals, true)

	arg = hasStringArgument(c, atom.Arguments[1], "test2")
	arg = hasStringArgument(c, arg.To, "$")
	c.Assert(IsFunctionAtom(arg.To), Equals, true)

	atom = hasStringArgument(c, atom.To, "__id")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseScript_env_lookup_fails_on_illegal_ident(c *C) {
	_, err := ParseScript("$_gcp")
	c.Assert(err, Not(IsNil))
	_, err = ParseScript("$12gcp")
	c.Assert(err, Not(IsNil))
	_, err = ParseScript("$.gcp")
	c.Assert(err, Not(IsNil))
	_, err = ParseScript("$-gcp")
	c.Assert(err, Not(IsNil))
}

func (p *parserSuite) Test_ParseScript_method_call(c *C) {
	result, err := ParseScript("$test.concat($test2)")
	c.Assert(err, IsNil)
	c.Assert(IsApplyAtom(result), Equals, true)

	atom := ExpectApplyAtom(result)
	c.Assert(IsApplyAtom(atom.To), Equals, true)
	c.Assert(atom.Arguments, HasLen, 2)

	arg := hasStringArgument(c, atom.Arguments[0], "test")
	arg = hasStringArgument(c, arg.To, "$")
	c.Assert(IsFunctionAtom(arg.To), Equals, true)

	arg = hasStringArgument(c, atom.Arguments[1], "test2")
	arg = hasStringArgument(c, arg.To, "$")
	c.Assert(IsFunctionAtom(arg.To), Equals, true)

	atom = hasStringArgument(c, atom.To, "__concat")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_Parse_And_Eval_Env_Lookup_with_function_calls(c *C) {
	inputsDict := LiftDict(map[string]Script{
		"version": LiftString("1.0"),
		"extra":   LiftString("alpha"),
	})
	gcpDict := LiftDict(map[string]Script{
		"inputs": inputsDict,
	})
	globalsDict := map[string]Script{
		"gcp": gcpDict,
	}
	env := NewScriptEnvironmentWithGlobals(globalsDict)

	script, err := ParseScript(`$__concat($gcp.inputs.version, "-", $gcp.inputs.extra)`)
	c.Assert(err, IsNil)

	result, err := EvalToGoValue(script, env)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "1.0-alpha")
}

func (p *parserSuite) Test_Parse_And_Eval_Env_Lookup_with_method_calls(c *C) {
	inputsDict := LiftDict(map[string]Script{
		"version": LiftString("1.0"),
		"extra":   LiftString("alpha"),
	})
	gcpDict := LiftDict(map[string]Script{
		"inputs": inputsDict,
	})
	globalsDict := map[string]Script{
		"gcp": gcpDict,
	}
	env := NewScriptEnvironmentWithGlobals(globalsDict)

	script, err := ParseScript(`$gcp.inputs.version.concat("-", $gcp.inputs.extra)`)
	c.Assert(err, IsNil)

	result, err := EvalToGoValue(script, env)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "1.0-alpha")
}
