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
	. "github.com/ankyra/escape-client/model/interfaces"
	. "gopkg.in/check.v1"
	"testing"
)

type exprSuite struct{}

var _ = Suite(&exprSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *exprSuite) Test_Eval_String(c *C) {
	v := LiftString("test")
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test")
}

func (s *exprSuite) Test_Eval_Integer(c *C) {
	v := LiftInteger(12)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, 12)
}

func (s *exprSuite) Test_Eval_Function(c *C) {
	id := LiftFunction(builtinId)
	_, err := EvalToGoValue(id, nil)
	c.Assert(err, IsNil)
}

func (s *exprSuite) Test_Eval_Function_Apply(c *C) {
	id := LiftFunction(builtinId)
	v := LiftString("test")
	apply := NewApply(id, []Script{v})

	result, err := EvalToGoValue(apply, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test")
}

func (s *exprSuite) Test_Eval_Function_Apply_Recursive(c *C) {
	id := LiftFunction(builtinId)
	v := LiftString("test")
	apply2 := NewApply(id, []Script{v})
	apply1 := NewApply(id, []Script{apply2})

	result, err := EvalToGoValue(apply1, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test")
}

func (s *exprSuite) Test_Eval_Map_Apply(c *C) {
	v := LiftString("test")
	dict := LiftDict(map[string]Script{
		"test": LiftString("test value"),
	})
	apply := NewApply(dict, []Script{v})

	result, err := EvalToGoValue(apply, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test value")
}

func (s *exprSuite) Test_Eval_Env_Lookup_BuiltIn(c *C) {
	envLookup := LiftFunction(builtinEnvLookup)
	v := LiftString("test")
	apply := NewApply(envLookup, []Script{v})

	env := NewScriptEnvironment()
	(*env)["test"] = LiftString("test value")

	result, err := EvalToGoValue(apply, env)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test value")
}

// $test.version
//   =>
// apply(apply(apply(env_lookup, "$"), "test"), "version")

func (s *exprSuite) Test_Eval(c *C) {
	globals := LiftString("$")
	key1 := LiftString("test")
	key2 := LiftString("version")
	envLookup := LiftFunction(builtinEnvLookup)
	apply3 := NewApply(envLookup, []Script{globals})
	apply2 := NewApply(apply3, []Script{key1})
	apply1 := NewApply(apply2, []Script{key2})

	testDict := LiftDict(map[string]Script{
		"version": LiftString("1.0"),
	})
	globalsDict := LiftDict(map[string]Script{
		"test": testDict,
	})
	env := NewScriptEnvironment()
	(*env)["$"] = globalsDict
	result, err := EvalToGoValue(apply1, env)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "1.0")
}
