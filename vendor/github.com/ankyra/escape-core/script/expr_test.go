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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

type exprSuite struct{}

var _ = Suite(&exprSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *exprSuite) Test_Lift_nil_returns_empty_string(c *C) {
	v, err := Lift(nil)
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(v), Equals, true)
	c.Assert(ExpectStringAtom(v), Equals, "")
}

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

func (s *exprSuite) Test_Lift_Func_string_to_string(c *C) {
	v, err := Lift(strings.ToLower)
	c.Assert(err, IsNil)
	c.Assert(IsFunctionAtom(v), Equals, true)
	apply := NewApply(v, []Script{LiftString("TEST")})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "test")
}

func (s *exprSuite) Test_Lift_Func_string_string_to_string_slice(c *C) {
	v, err := Lift(strings.Split)
	c.Assert(err, IsNil)
	c.Assert(IsFunctionAtom(v), Equals, true)
	apply := NewApply(v, []Script{LiftString("test1 test2"), LiftString(" ")})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsListAtom(result), Equals, true)
	c.Assert(ExpectListAtom(result)[0], DeepEquals, LiftString("test1"))
	c.Assert(ExpectListAtom(result)[1], DeepEquals, LiftString("test2"))
}

func (s *exprSuite) Test_Lift_Func_string_slice_string_to_string(c *C) {
	v, err := Lift(strings.Join)
	c.Assert(err, IsNil)
	c.Assert(IsFunctionAtom(v), Equals, true)
	lst := []Script{
		LiftString("test1"),
		LiftString("test2"),
	}
	apply := NewApply(v, []Script{LiftList(lst), LiftString(" ")})
	result, err := apply.Eval(nil)
	c.Assert(err, IsNil)
	c.Assert(IsStringAtom(result), Equals, true)
	c.Assert(ExpectStringAtom(result), Equals, "test1 test2")
}

func (s *exprSuite) Test_Eval_String(c *C) {
	v := LiftString("test")
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "test")
}

func (s *exprSuite) Test_Equals_String(c *C) {
	v := LiftString("test")
	v2 := LiftString("test2")
	v3 := LiftString("test")
	c.Assert(v.Equals(v2), Equals, false)
	c.Assert(v3.Equals(v2), Equals, false)
	c.Assert(v.Equals(v3), Equals, true)
}
func (s *exprSuite) Test_Equals_Bool(c *C) {
	v := LiftBool(true)
	v2 := LiftBool(false)
	v3 := LiftBool(true)
	c.Assert(v.Equals(v2), Equals, false)
	c.Assert(v3.Equals(v2), Equals, false)
	c.Assert(v.Equals(v3), Equals, true)
}

func (s *exprSuite) Test_Equals_Integer(c *C) {
	v := LiftInteger(1)
	v2 := LiftInteger(2)
	v3 := LiftInteger(1)
	c.Assert(v.Equals(v2), Equals, false)
	c.Assert(v3.Equals(v2), Equals, false)
	c.Assert(v.Equals(v3), Equals, true)
}

func (s *exprSuite) Test_Equals_List(c *C) {
	v := LiftList([]Script{})
	v2 := LiftList([]Script{LiftString("test"), LiftString("test2")})
	v3 := LiftList([]Script{})
	v4 := LiftList([]Script{LiftString("test"), LiftString("test2")})
	v5 := LiftList([]Script{LiftString("test")})
	v6 := LiftList([]Script{LiftString("test"), LiftString("test5")})
	c.Assert(v.Equals(v2), Equals, false)
	c.Assert(v3.Equals(v2), Equals, false)
	c.Assert(v.Equals(v3), Equals, true)
	c.Assert(v2.Equals(v4), Equals, true)
	c.Assert(v2.Equals(v5), Equals, false)
	c.Assert(v2.Equals(v6), Equals, false)
	c.Assert(v5.Equals(v2), Equals, false)
}

func (s *exprSuite) Test_Equals_Dict(c *C) {
	v := LiftDict(map[string]Script{})
	v2 := LiftDict(map[string]Script{
		"test": LiftString("yo"),
	})
	v3 := LiftDict(map[string]Script{})
	v4 := LiftDict(map[string]Script{
		"test": LiftInteger(1),
	})
	v5 := LiftDict(map[string]Script{
		"test": LiftString("yo"),
	})
	v6 := LiftDict(map[string]Script{
		"test":  LiftString("yo"),
		"test2": LiftString("yop"),
	})
	c.Assert(v.Equals(v2), Equals, false)
	c.Assert(v3.Equals(v2), Equals, false)
	c.Assert(v.Equals(v3), Equals, true)
	c.Assert(v2.Equals(v4), Equals, false)
	c.Assert(v2.Equals(v5), Equals, true)
	c.Assert(v2.Equals(v6), Equals, false)
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

func (s *exprSuite) Test_Eval_Lambda(c *C) {
	body := NewApply(ShouldLift(builtinAdd), []Script{ShouldParse("$var1"), ShouldParse("$var2")})
	lambda := NewLambda([]string{"var1", "var2"}, body)
	v := NewApply(lambda, []Script{LiftInteger(1), LiftInteger(3)})
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, 4)
}

func (s *exprSuite) Test_Eval_List(c *C) {
	list := []Script{LiftString("test"), LiftInteger(12)}
	v := LiftList(list)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, []interface{}{"test", 12})
}

func (s *exprSuite) Test_Eval_List_evals_recursive(c *C) {
	f := LiftFunction(builtinId)
	t := LiftString("test")
	a := NewApply(f, []Script{t})
	list := []Script{a}
	v := LiftList(list)
	result, err := EvalToGoValue(v, nil)
	c.Assert(err, IsNil)
	c.Assert(result, DeepEquals, []interface{}{"test"})
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
