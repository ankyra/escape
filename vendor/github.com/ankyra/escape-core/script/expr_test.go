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
	"io/ioutil"
	"os"
	"testing"
)

type exprSuite struct{}

var _ = Suite(&exprSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *exprSuite) Test_Lift_ScriptString(c *C) {
	v, err := Lift(LiftString("string"))
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(v), Equals, true)
	c.Assert(ExpectStringAtom(v), Equals, "string")
}
func (s *exprSuite) Test_Lift_String(c *C) {
	v, err := Lift("string")
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(v), Equals, true)
	c.Assert(ExpectStringAtom(v), Equals, "string")
}
func (s *exprSuite) Test_Lift_Bool(c *C) {
	v, err := Lift(true)
	c.Assert(err, IsNil)
	c.Assert(IsBoolAtom(v), Equals, true)
	c.Assert(ExpectBoolAtom(v), Equals, true)
}
func (s *exprSuite) Test_Lift_Integer(c *C) {
	v, err := Lift(12)
	c.Assert(err, IsNil)
	c.Assert(IsIntegerAtom(v), Equals, true)
	c.Assert(ExpectIntegerAtom(v), Equals, 12)
}
func (s *exprSuite) Test_Lift_Float(c *C) {
	v, err := Lift(12.6)
	c.Assert(err, IsNil)
	c.Assert(IsIntegerAtom(v), Equals, true)
	c.Assert(ExpectIntegerAtom(v), Equals, 12)
}
func (s *exprSuite) Test_Lift_List(c *C) {
	list := []interface{}{"test", 12}
	v, err := Lift(list)
	expected := []Script{LiftString("test"), LiftInteger(12)}
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(v), Equals, true)
	c.Assert(ExpectListAtom(v), DeepEquals, expected)
}
func (s *exprSuite) Test_Lift_List_fails_if_item_type_not_supported(c *C) {
	list := []interface{}{"test", struct{}{}}
	_, err := Lift(list)
	c.Assert(err, Not(IsNil))
}
func (s *exprSuite) Test_Lift_Script_List(c *C) {
	list := []Script{LiftString("test"), LiftInteger(12)}
	v, err := Lift(list)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(v), Equals, true)
	c.Assert(ExpectListAtom(v), DeepEquals, list)
}
func (s *exprSuite) Test_Lift_ScriptDict(c *C) {
	dict := map[string]Script{
		"test": LiftString("value"),
	}
	v, err := Lift(dict)
	c.Assert(err, IsNil)
	c.Assert(IsDictAtom(v), Equals, true)
	c.Assert(ExpectDictAtom(v), DeepEquals, dict)
}
func (s *exprSuite) Test_Lift_Dict(c *C) {
	dict := map[string]interface{}{
		"test": "value",
		"recurse": map[string]interface{}{
			"test": "value2",
		},
	}
	v, err := Lift(dict)
	expected := map[string]Script{
		"test": LiftString("value"),
		"recurse": LiftDict(map[string]Script{
			"test": LiftString("value2"),
		}),
	}
	c.Assert(err, IsNil)
	c.Assert(IsDictAtom(v), Equals, true)
	c.Assert(ExpectDictAtom(v), DeepEquals, expected)
}
func (s *exprSuite) Test_Lift_Dict_interface(c *C) {
	dict := map[interface{}]interface{}{
		"test": "value",
		"recurse": map[string]interface{}{
			"test": "value2",
		},
	}
	v, err := Lift(dict)
	expected := map[string]Script{
		"test": LiftString("value"),
		"recurse": LiftDict(map[string]Script{
			"test": LiftString("value2"),
		}),
	}
	c.Assert(err, IsNil)
	c.Assert(IsDictAtom(v), Equals, true)
	c.Assert(ExpectDictAtom(v), DeepEquals, expected)
}

func (s *exprSuite) Test_Eval_String(c *C) {
	v := LiftString("test")
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test")
}

func (s *exprSuite) Test_IsStringAtom(c *C) {
	c.Assert(IsStringAtom(LiftString("test")), Equals, true)
	c.Assert(IsStringAtom(LiftInteger(12)), Equals, false)
	c.Assert(IsStringAtom(LiftFunction(builtinId)), Equals, false)
	c.Assert(IsStringAtom(NewApply(LiftFunction(builtinId), nil)), Equals, false)
}

func (s *exprSuite) Test_ExpectStringAtom(c *C) {
	c.Assert(ExpectStringAtom(LiftString("test")), Equals, "test")
	c.Assert(func() { ExpectStringAtom(LiftInteger(12)) }, Panics, "Expecting string type, got integer")
}

func (s *exprSuite) Test_Eval_Bool(c *C) {
	v := LiftBool(false)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, false)
}

func (s *exprSuite) Test_IsBoolAtom(c *C) {
	c.Assert(IsBoolAtom(LiftBool(true)), Equals, true)
	c.Assert(IsBoolAtom(LiftInteger(12)), Equals, false)
	c.Assert(IsBoolAtom(LiftFunction(builtinId)), Equals, false)
	c.Assert(IsBoolAtom(NewApply(LiftFunction(builtinId), nil)), Equals, false)
}

func (s *exprSuite) Test_ExpectBoolAtom(c *C) {
	c.Assert(ExpectBoolAtom(LiftBool(false)), Equals, false)
	c.Assert(func() { ExpectBoolAtom(LiftInteger(12)) }, Panics, "Expecting bool type, got integer")
}

func (s *exprSuite) Test_Eval_Integer(c *C) {
	v := LiftInteger(12)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, 12)
}

func (s *exprSuite) Test_IsIntegerAtom(c *C) {
	c.Assert(IsIntegerAtom(LiftInteger(12)), Equals, true)
	c.Assert(IsIntegerAtom(LiftString("test")), Equals, false)
	c.Assert(IsIntegerAtom(LiftFunction(builtinId)), Equals, false)
	c.Assert(IsIntegerAtom(NewApply(LiftFunction(builtinId), nil)), Equals, false)
}

func (s *exprSuite) Test_ExpectIntegerAtom(c *C) {
	c.Assert(ExpectIntegerAtom(LiftInteger(12)), Equals, 12)
	c.Assert(func() { ExpectIntegerAtom(LiftString("test")) }, Panics, "Expecting integer type, got string")
}

func (s *exprSuite) Test_Eval_List(c *C) {
	list := []Script{LiftString("test"), LiftInteger(12)}
	v := LiftList(list)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, list)
}

func (s *exprSuite) Test_Eval_Dict(c *C) {
	dict := map[string]Script{
		"test": LiftString("value"),
	}
	v := LiftDict(dict)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, dict)
}

func (s *exprSuite) Test_Eval_IsDictAtom(c *C) {
	dict := map[string]Script{
		"test": LiftString("value"),
	}
	v := LiftDict(dict)
	c.Assert(IsDictAtom(v), Equals, true)
	c.Assert(IsDictAtom(LiftInteger(12)), Equals, false)
	c.Assert(IsDictAtom(LiftString("test")), Equals, false)
	c.Assert(IsDictAtom(LiftFunction(builtinId)), Equals, false)
	c.Assert(IsDictAtom(NewApply(LiftFunction(builtinId), nil)), Equals, false)
}

func (s *exprSuite) Test_ExpectDictAtom(c *C) {
	dict := map[string]Script{
		"test": LiftString("value"),
	}
	v := LiftDict(dict)
	c.Assert(ExpectDictAtom(v), DeepEquals, dict)
	c.Assert(func() { ExpectDictAtom(LiftString("test")) }, Panics, "Expecting dict type, got string")
}

func (s *exprSuite) Test_ExpectDict(c *C) {
	dict := map[string]Script{
		"test": LiftString("value"),
	}
	v := LiftDict(dict)
	expect := map[string]interface{}{
		"test": "value",
	}
	c.Assert(ExpectDict(v), DeepEquals, expect)
}
func (s *exprSuite) Test_ExpectDict_fails_with_wrong_type(c *C) {
	c.Assert(func() { ExpectDict(LiftString("test")) }, Panics, "Expecting dict type, got string")
}
func (s *exprSuite) Test_ExpectDict_fails_with_unapplied_func(c *C) {
	dict := map[string]Script{
		"test": NewApply(LiftFunction(builtinId), []Script{}),
	}
	v := LiftDict(dict)
	c.Assert(func() { ExpectDict(v) }, Panics, "Function application can not be converted to Go value (forgot to eval?)")
}

func (s *exprSuite) Test_Eval_Function(c *C) {
	id := LiftFunction(builtinId)
	_, err := EvalToGoValue(id, nil)
	c.Assert(err, IsNil)
}
func (s *exprSuite) Test_Eval_IsFunction(c *C) {
	c.Assert(IsFunctionAtom(LiftFunction(builtinId)), Equals, true)
	c.Assert(IsFunctionAtom(LiftInteger(12)), Equals, false)
	c.Assert(IsFunctionAtom(LiftString("test")), Equals, false)
	c.Assert(IsFunctionAtom(NewApply(LiftFunction(builtinId), nil)), Equals, false)
}
func (s *exprSuite) Test_Eval_ExpectFunctionAtom(c *C) {
	id := LiftFunction(builtinId)
	ExpectFunctionAtom(id)
}
func (s *exprSuite) Test_ExpectFunctionAtom_fails_with_wrong_type(c *C) {
	c.Assert(func() { ExpectFunctionAtom(LiftString("test")) }, Panics, "Expecting function type, got string")
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

func (s *exprSuite) Test_Eval_String_Apply_file(c *C) {
	apply := NewApply(LiftString("test content"), []Script{LiftString("file")})
	filePath, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(filePath), Equals, true)
	result, err := ioutil.ReadFile(ExpectStringAtom(filePath))
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, []byte("test content"))
	os.RemoveAll(ExpectStringAtom(filePath))
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
