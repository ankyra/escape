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
		"version": LiftString("1.2.3"),
		"extra":   LiftString("alpha"),
	})
	gcpDict := LiftDict(map[string]Script{
		"inputs": inputsDict,
	})
	globalsDict := map[string]Script{
		"gcp": gcpDict,
		"lst": LiftList([]Script{LiftString("first item"), LiftString("second item")}),
	}
	env := NewScriptEnvironmentWithGlobals(globalsDict)

	cases := map[string]string{
		`$gcp.inputs.version.concat("-", $gcp.inputs.extra)`:     `1.2.3-alpha`,
		`$__concat($gcp.inputs.version, "-", $gcp.inputs.extra)`: `1.2.3-alpha`,
		`$lst[0]`:                                                               `first item`,
		`$lst.join(", ")`:                                                       `first item, second item`,
		`$lst[0:].join(", ")`:                                                   `first item, second item`,
		`$lst[:2].join(", ")`:                                                   `first item, second item`,
		`$lst[0:2].join(", ")`:                                                  `first item, second item`,
		`$lst[0:1].join(", ")`:                                                  `first item`,
		`$lst[:1].join(", ")`:                                                   `first item`,
		`$lst[:-1].join(", ")`:                                                  `first item`,
		`$func(listVar, joinStr) { $listVar.join($joinStr) }($lst, " ## ")`:     `first item ## second item`,
		`$func(listVar, joinStr) { $listVar.join($joinStr) }($lst[1:], " ## ")`: `second item`,
		`$gcp.inputs.version.track_version()`:                                   `1.2.3.@`,
		`$gcp.inputs.version.track_major_version()`:                             `1.@`,
		`$gcp.inputs.version.track_minor_version()`:                             `1.2.@`,
		`$gcp.inputs.version.track_patch_version()`:                             `1.2.3.@`,
	}
	for testCase, expected := range cases {
		script, err := ParseScript(testCase)
		c.Assert(err, IsNil, Commentf("Couldn't parse '%s'", testCase))

		result, err := EvalToGoValue(script, env)
		c.Assert(err, IsNil)
		c.Assert(result, Equals, expected, Commentf("Error in '%s'", testCase))
	}
}

func (p *parserSuite) Test_Parse_And_Eval_Env_Lookup_failing_cases(c *C) {
	inputsDict := LiftDict(map[string]Script{
		"version": LiftString("1.0"),
		"extra":   LiftString("alpha"),
	})
	gcpDict := LiftDict(map[string]Script{
		"inputs": inputsDict,
	})
	globalsDict := map[string]Script{
		"gcp": gcpDict,
		"lst": LiftList([]Script{LiftString("first item")}),
	}
	env := NewScriptEnvironmentWithGlobals(globalsDict)

	cases := []string{
		`$lst[0].replace()`,
		`$lst[-1]`,
		`$lst[1]`,
	}
	for _, testCase := range cases {
		script, err := ParseScript(testCase)
		c.Assert(err, IsNil, Commentf("Couldn't parse '%s'", testCase))

		_, err = EvalToGoValue(script, env)
		c.Assert(err, Not(IsNil), Commentf("Should have failed '%s'", testCase))
	}
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

func (p *parserSuite) Test_ParseListIndex(c *C) {
	to := LiftString("whatever")
	result := parseListIndex(to, "[12]")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")

	apply := result.Result
	c.Assert(IsApplyAtom(apply), Equals, true)

	atom := ExpectApplyAtom(apply)
	c.Assert(IsApplyAtom(atom.To), Equals, true)
	c.Assert(atom.Arguments, HasLen, 2)

	c.Assert(ExpectStringAtom(atom.Arguments[0]), Equals, "whatever")
	c.Assert(ExpectIntegerAtom(atom.Arguments[1]), Equals, 12)
	atom = hasStringArgument(c, atom.To, "__list_index")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseListSlice(c *C) {
	to := LiftString("whatever")
	result := parseListIndex(to, "[12:13]")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")

	apply := result.Result
	c.Assert(IsApplyAtom(apply), Equals, true)

	atom := ExpectApplyAtom(apply)
	c.Assert(IsApplyAtom(atom.To), Equals, true)

	c.Assert(atom.Arguments, HasLen, 3)
	c.Assert(ExpectStringAtom(atom.Arguments[0]), Equals, "whatever")
	c.Assert(ExpectIntegerAtom(atom.Arguments[1]), Equals, 12)
	c.Assert(ExpectIntegerAtom(atom.Arguments[2]), Equals, 13)

	atom = hasStringArgument(c, atom.To, "__list_slice")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseListSlice_to_end(c *C) {
	to := LiftString("whatever")
	result := parseListIndex(to, "[12:]")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")

	apply := result.Result
	c.Assert(IsApplyAtom(apply), Equals, true)

	atom := ExpectApplyAtom(apply)
	c.Assert(IsApplyAtom(atom.To), Equals, true)

	c.Assert(atom.Arguments, HasLen, 2)
	c.Assert(ExpectStringAtom(atom.Arguments[0]), Equals, "whatever")
	c.Assert(ExpectIntegerAtom(atom.Arguments[1]), Equals, 12)

	atom = hasStringArgument(c, atom.To, "__list_slice")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseListSlice_from_beginning(c *C) {
	to := LiftString("whatever")
	result := parseListIndex(to, "[:2]")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")

	apply := result.Result
	c.Assert(IsApplyAtom(apply), Equals, true)

	atom := ExpectApplyAtom(apply)
	c.Assert(IsApplyAtom(atom.To), Equals, true)

	c.Assert(atom.Arguments, HasLen, 3)
	c.Assert(ExpectStringAtom(atom.Arguments[0]), Equals, "whatever")
	c.Assert(ExpectIntegerAtom(atom.Arguments[1]), Equals, 0)
	c.Assert(ExpectIntegerAtom(atom.Arguments[2]), Equals, 2)

	atom = hasStringArgument(c, atom.To, "__list_slice")
	atom = hasStringArgument(c, atom.To, "$")
	c.Assert(IsFunctionAtom(atom.To), Equals, true)
}

func (p *parserSuite) Test_ParseExpression_int(c *C) {
	result := parseExpression("12")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")
	c.Assert(ExpectIntegerAtom(result.Result), Equals, 12)
}

func (p *parserSuite) Test_ParseExpression_negative_int(c *C) {
	result := parseExpression("-12")
	c.Assert(result.Error, IsNil)
	c.Assert(result.Rest, Equals, "")
	c.Assert(ExpectIntegerAtom(result.Result), Equals, -12)
}

func (p *parserSuite) Test_ParseExpression_negative(c *C) {
	result := parseExpression("-")
	c.Assert(result.Error, Not(IsNil))
}

func (p *parserSuite) Test_ParseExpression_table(c *C) {
	cases := []string{
		`12`,
		`-12`,
		`"string"`,
		"$test",
		"$test.test",
		"$test.test[12]",
		"$test.test(12)",
		"$test.test(12, 12, 123)",
		`$__test()`,
		`$test.test()`,
		`$test.test("test")`,
		`$test.test(12, "test")`,
		`$test.test(12, "test", $recurse.test(12, "test"))`,
		`$__test()[0]`,
		`$test.test[12].test`,
		`$test.test()[12]`,
		`$test.test(12, "test", $test.test[1].contact("hallo", 12))[12]`,
	}
	for _, testCase := range cases {
		result := parseExpression(testCase)
		c.Assert(result.Error, IsNil, Commentf("Couldn't parse '%s'", testCase))
		c.Assert(result.Rest, Equals, "", Commentf("Coulnd't parse '%s'", testCase))
	}
}

func (p *parserSuite) Test_ParseExpression_fail_table(c *C) {
	cases := []string{
		"",
		"-",
		"$test()",
		`$test.test(`,
		`$test.test(12`,
		`$test.test(12, `,
		`$test.test(12, "test"`,
		`$test.test[`,
		`$test.test[12`,
		`$test.test[12:`,
	}
	for _, testCase := range cases {
		result := parseExpression(testCase)
		c.Assert(result.Error != nil || result.Rest != "", Equals, true, Commentf("Shouldn't be able to parse '%s' (error: %s, rest: %s)", testCase, result.Error, result.Rest))
	}
}
