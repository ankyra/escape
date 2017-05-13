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

func (s *exprSuite) Test_Builtin_Id(c *C) {
	result, err := builtinId(nil, []Script{LiftString("test")})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "test")
}
func (s *exprSuite) Test_Builtin_Id_fails_without_args(c *C) {
	_, err := builtinId(nil, []Script{})
	c.Assert(err, Not(IsNil))
}

func (s *exprSuite) Test_Builtin_Id_fails_too_many_args(c *C) {
	_, err := builtinId(nil, []Script{LiftString("t"), LiftString("t2")})
	c.Assert(err, Not(IsNil))
}

func (s *exprSuite) Test_Builtin_env_lookup(c *C) {
	e := map[string]Script{
		"key": LiftString("value"),
	}
	env := NewScriptEnvironmentFromMap(e)
	result, err := builtinEnvLookup(env, []Script{LiftString("key")})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "value")
}

func (s *exprSuite) Test_Builtin_env_lookup_fails_without_args(c *C) {
	e := map[string]Script{
		"key": LiftString("value"),
	}
	env := NewScriptEnvironmentFromMap(e)
	_, err := builtinEnvLookup(env, []Script{})
	c.Assert(err, Not(IsNil))
}

func (s *exprSuite) Test_Builtin_env_lookup_fails_if_key_is_not_string(c *C) {
	e := map[string]Script{
		"key": LiftString("value"),
	}
	env := NewScriptEnvironmentFromMap(e)
	_, err := builtinEnvLookup(env, []Script{LiftInteger(12)})
	c.Assert(err, Not(IsNil))
}

func (s *exprSuite) Test_Builtin_env_lookup_fails_if_key_is_not_found(c *C) {
	e := map[string]Script{
		"key": LiftString("value"),
	}
	env := NewScriptEnvironmentFromMap(e)
	_, err := builtinEnvLookup(env, []Script{LiftString("not found")})
	c.Assert(err, Not(IsNil))
}
