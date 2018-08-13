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

func (s *exprSuite) Test_Builtin_Concat_empty(c *C) {
	result, err := builtinConcat(nil, []Script{})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "")
}
func (s *exprSuite) Test_Builtin_Concat_1(c *C) {
	result, err := builtinConcat(nil, []Script{LiftString("test")})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "test")
}
func (s *exprSuite) Test_Builtin_Concat_2(c *C) {
	result, err := builtinConcat(nil, []Script{LiftString("test"), LiftString(" testing"), LiftString(" testing")})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "test testing testing")
}
func (s *exprSuite) Test_Builtin_Concat_with_integer(c *C) {
	result, err := builtinConcat(nil, []Script{LiftInteger(12), LiftInteger(100), LiftString("test")})
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "12100test")
}

func (s *exprSuite) Test_Builtin_Concat_fails_with_wrong_type(c *C) {
	_, err := builtinConcat(nil, []Script{LiftList([]Script{})})
	c.Assert(err, Not(IsNil))
}

func (s *exprSuite) Test_Builtin_base64_encode(c *C) {
	for _, f := range Stdlib {
		if f.Id == "base64_encode" {
			apply := NewApply(f.Func, []Script{LiftString("TEST")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsStringAtom(result), Equals, true)
			c.Assert(ExpectStringAtom(result), Equals, "VEVTVA==")
		}
	}
}

func (s *exprSuite) Test_Builtin_base64_decode(c *C) {
	for _, f := range Stdlib {
		if f.Id == "base64_decode" {
			apply := NewApply(f.Func, []Script{LiftString("VEVTVA==")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsStringAtom(result), Equals, true)
			c.Assert(ExpectStringAtom(result), Equals, "TEST")
		}
	}
}

func (s *exprSuite) Test_Builtin_base64_decode_fails_if_invalid(c *C) {
	for _, f := range Stdlib {
		if f.Id == "base64_decode" {
			apply := NewApply(f.Func, []Script{LiftString("1")})
			_, err := apply.Eval(nil)
			c.Assert(err, Not(IsNil))
		}
	}
}

func (s *exprSuite) Test_Builtin_replace(c *C) {
	for _, f := range Stdlib {
		if f.Id == "replace" {
			apply := NewApply(f.Func, []Script{
				LiftString("TEST"), LiftString("T"), LiftString("B"), LiftInteger(1)})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsStringAtom(result), Equals, true)
			c.Assert(ExpectStringAtom(result), Equals, "BEST")
		}
	}
}

func (s *exprSuite) Test_Builtin_slice_no_end(c *C) {
	lst := LiftList([]Script{LiftString("first"), LiftString("second")})
	apply := NewApply(LiftFunction(builtinListSlice), []Script{lst, LiftInteger(0)})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(result), Equals, true)
	lstResult := ExpectListAtom(result)
	c.Assert(lstResult, HasLen, 2)
	c.Assert(ExpectStringAtom(lstResult[0]), Equals, "first")
	c.Assert(ExpectStringAtom(lstResult[1]), Equals, "second")

	apply = NewApply(LiftFunction(builtinListSlice), []Script{lst, LiftInteger(1)})
	result, err = apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(result), Equals, true)
	lstResult = ExpectListAtom(result)
	c.Assert(lstResult, HasLen, 1)
	c.Assert(ExpectStringAtom(lstResult[0]), Equals, "second")
}

func (s *exprSuite) Test_Builtin_slice(c *C) {
	lst := LiftList([]Script{LiftString("first"), LiftString("second")})
	apply := NewApply(LiftFunction(builtinListSlice), []Script{lst, LiftInteger(0), LiftInteger(2)})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(result), Equals, true)
	lstResult := ExpectListAtom(result)
	c.Assert(lstResult, HasLen, 2)
	c.Assert(ExpectStringAtom(lstResult[0]), Equals, "first")
	c.Assert(ExpectStringAtom(lstResult[1]), Equals, "second")

	apply = NewApply(LiftFunction(builtinListSlice), []Script{lst, LiftInteger(1), LiftInteger(2)})
	result, err = apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(result), Equals, true)
	lstResult = ExpectListAtom(result)
	c.Assert(lstResult, HasLen, 1)
	c.Assert(ExpectStringAtom(lstResult[0]), Equals, "second")
}

func (s *exprSuite) Test_Builtin_length(c *C) {
	lst := LiftList([]Script{LiftString("first"), LiftString("second")})
	apply := NewApply(LiftFunction(builtinListLength), []Script{lst})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsIntegerAtom(result), Equals, true)
	intResult := ExpectIntegerAtom(result)
	c.Assert(intResult, Equals, 2)
}

func (s *exprSuite) Test_Builtin_length_empty_list(c *C) {
	lst := LiftList(nil)
	apply := NewApply(LiftFunction(builtinListLength), []Script{lst})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsIntegerAtom(result), Equals, true)
	intResult := ExpectIntegerAtom(result)
	c.Assert(intResult, Equals, 0)
}

func (s *exprSuite) Test_Builtin_length_on_strings(c *C) {
	lst := LiftString("hello")
	apply := NewApply(LiftFunction(builtinListLength), []Script{lst})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsIntegerAtom(result), Equals, true)
	intResult := ExpectIntegerAtom(result)
	c.Assert(intResult, Equals, 5)
}

func (s *exprSuite) Test_Builtin_not(c *C) {
	lst := LiftBool(true)
	apply := NewApply(LiftFunction(builtinNOT), []Script{lst})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsBoolAtom(result), Equals, true)
	c.Assert(ExpectBoolAtom(result), Equals, false)
}

func (s *exprSuite) Test_Builtin_equals(c *C) {
	v1 := LiftBool(true)
	v2 := LiftBool(false)
	apply := NewApply(LiftFunction(builtinEquals), []Script{v1, v2})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsBoolAtom(result), Equals, true)
	c.Assert(ExpectBoolAtom(result), Equals, false)
}

func (s *exprSuite) Test_Builtin_equals_true(c *C) {
	v1 := LiftString("test")
	v2 := LiftString("test")
	apply := NewApply(LiftFunction(builtinEquals), []Script{v1, v2})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsBoolAtom(result), Equals, true)
	c.Assert(ExpectBoolAtom(result), Equals, true)
}

func (s *exprSuite) Test_Builtin_path_exists_true(c *C) {
	for _, f := range Stdlib {
		if f.Id == "path_exists" {
			apply := NewApply(f.Func, []Script{LiftString("builtins_test.go")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, true)
		}
	}
}

func (s *exprSuite) Test_Builtin_path_exists_false(c *C) {
	for _, f := range Stdlib {
		if f.Id == "path_exists" {
			apply := NewApply(f.Func, []Script{LiftString("non_existent_path")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, false)
		}
	}
}

func (s *exprSuite) Test_Builtin_file_exists_true(c *C) {
	for _, f := range Stdlib {
		if f.Id == "file_exists" {
			apply := NewApply(f.Func, []Script{LiftString("builtins_test.go")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, true)
		}
	}
}
func (s *exprSuite) Test_Builtin_file_exists_false(c *C) {
	for _, f := range Stdlib {
		if f.Id == "file_exists" {
			apply := NewApply(f.Func, []Script{LiftString("non_existent_path")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, false)
		}
	}
}

func (s *exprSuite) Test_Builtin_dir_exists_true(c *C) {
	for _, f := range Stdlib {
		if f.Id == "dir_exists" {
			apply := NewApply(f.Func, []Script{LiftString("builtins_test.go")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, false)
		}
	}
}
func (s *exprSuite) Test_Builtin_dir_exists_false(c *C) {
	for _, f := range Stdlib {
		if f.Id == "dir_exists" {
			apply := NewApply(f.Func, []Script{LiftString("non_existent_path")})
			result, err := apply.Eval(nil)
			c.Assert(err, IsNil)
			c.Assert(IsBoolAtom(result), Equals, true)
			c.Assert(ExpectBoolAtom(result), Equals, false)
		}
	}
}
