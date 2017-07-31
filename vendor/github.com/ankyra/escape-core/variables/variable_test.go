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

package variables

import (
	"github.com/ankyra/escape-core/script"
	. "gopkg.in/check.v1"
	"testing"
)

type variableSuite struct{}

var _ = Suite(&variableSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *variableSuite) Test_GetValue_fails_with_invalid_id(c *C) {
	tbl := []string{"_test", "123123", "//", "", "#"}
	for _, testCase := range tbl {
		_, err := NewVariableFromString(testCase, "string")
		c.Assert(err, Not(IsNil))
	}
}

func (s *variableSuite) Test_GetValue_fails_with_id_starting_with_previous(c *C) {
	tbl := []string{"previous_test", "PREVIOUS_test", "preVIOUS_test"}
	for _, testCase := range tbl {
		_, err := NewVariableFromString(testCase, "string")
		c.Assert(err, Not(IsNil))
	}
}

func (s *variableSuite) Test_GetValue_String_Variable(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": "test value",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "test value")
}

func (s *variableSuite) Test_GetValue_Uses_Default(c *C) {
	stringValue := "test"
	testCases := [][]interface{}{
		[]interface{}{"test value", "test value"},
		[]interface{}{12, "12"},
		[]interface{}{12.0, "12"},
		[]interface{}{true, "1"},
		[]interface{}{&stringValue, "test"},
		[]interface{}{[]interface{}{"test"}, `["test"]`},
	}
	for _, test := range testCases {
		unit, err := NewVariableFromString("test", "string")
		c.Assert(err, IsNil)
		unit.Default = test[0]
		val, err := unit.GetValue(nil, nil)
		c.Assert(err, IsNil)
		c.Assert(val, DeepEquals, test[1])
	}
}

func (s *variableSuite) Test_GetValue_OneOf_Variable(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = []interface{}{"valid", "also valid"}
	variableCtx := map[string]interface{}{
		"test": "valid",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "valid")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_Script(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = `$__split("valid,also valid", ",")`
	variableCtx := map[string]interface{}{
		"test": "valid",
	}
	globalsDict := map[string]script.Script{}
	env := script.NewScriptEnvironmentWithGlobals(globalsDict)
	val, err := unit.GetValue(&variableCtx, env)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "valid")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_List_Script(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = []interface{}{`$__concat("val", "id")`, `$__concat("also ", "valid")`}
	variableCtx := map[string]interface{}{
		"test": "valid",
	}
	globalsDict := map[string]script.Script{}
	env := script.NewScriptEnvironmentWithGlobals(globalsDict)
	val, err := unit.GetValue(&variableCtx, env)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "valid")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_Script_fails(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = `$__split("valid,also valid", ",")`
	variableCtx := map[string]interface{}{
		"test": "not valid",
	}
	globalsDict := map[string]script.Script{}
	env := script.NewScriptEnvironmentWithGlobals(globalsDict)
	_, err = unit.GetValue(&variableCtx, env)
	c.Assert(err.Error(), DeepEquals, "Expecting one of [\"valid\",\"also valid\"] for variable 'test'")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_Fails(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = []interface{}{"valid", "also valid"}
	variableCtx := map[string]interface{}{
		"test": "not valid",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), DeepEquals, "Expecting one of [\"valid\",\"also valid\"] for variable 'test'")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_string(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	unit.Items = "valid"
	variableCtx := map[string]interface{}{
		"test": "not valid",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), DeepEquals, "Unexpected value 'not valid' for variable 'test', only 'valid' is allowed")
}

func (s *variableSuite) Test_String_Variable_Converts_To_String_Value(c *C) {
	unit, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": 12,
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "12")
}

func (s *variableSuite) Test_GetValue_Integer_Variable(c *C) {
	unit, err := NewVariableFromString("test", "integer")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": 12,
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, 12)
}

func (s *variableSuite) Test_Integer_Variable_Expects_Integer_Value(c *C) {
	unit, err := NewVariableFromString("test", "integer")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": "test",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), Not(Equals), "")
}

func (s *variableSuite) Test_Integer_Variable_Expects_Integer_Value_Or_Convertable_String(c *C) {
	unit, err := NewVariableFromString("test", "integer")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": "12",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, 12)
}

func (s *variableSuite) Test_GetValue_List_Variable(c *C) {
	unit, err := NewVariableFromString("test", "list")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": []interface{}{"test value", "test value 2"},
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, []interface{}{"test value", "test value 2"})
}

func (s *variableSuite) Test_GetValue_List_Variable_Checks_String_Values(c *C) {
	unit, err := NewVariableFromString("test", "list")
	c.Assert(err, IsNil)
	variableCtx := map[string]interface{}{
		"test": []interface{}{12},
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), Equals, "Unexpected 'integer' value in list, expecting 'string' for variable 'test'")
}

func (s *variableSuite) Test_NewVariable_Default_Visible_and_EvalBeforeDependencies(c *C) {
	v := NewVariable()
	c.Assert(v.Visible, Equals, true)
	c.Assert(v.EvalBeforeDependencies, Equals, true)
}

func (s *variableSuite) Test_NewVariableFromString_Default_Visible_and_EvalBeforeDependencies(c *C) {
	v, err := NewVariableFromString("test", "string")
	c.Assert(err, IsNil)
	c.Assert(v.Visible, Equals, true)
	c.Assert(v.EvalBeforeDependencies, Equals, true)
}

func (s *variableSuite) Test_NewVariableFromDict_Default_Visible_and_EvalBeforeDependencies(c *C) {
	dict := map[interface{}]interface{}{
		"id": "test",
	}
	v, err := NewVariableFromDict(dict)
	c.Assert(err, IsNil)
	c.Assert(v.Visible, Equals, true)
	c.Assert(v.EvalBeforeDependencies, Equals, true)
}

func (s *variableSuite) Test_NewVariableFromDict_fails_if_missing_id(c *C) {
	dict := map[interface{}]interface{}{}
	_, err := NewVariableFromDict(dict)
	c.Assert(err.Error(), Equals, "Missing 'id' field in variable")
}

func (s *variableSuite) Test_NewVariableFromDict_fails_if_type_invalid(c *C) {
	dict := map[interface{}]interface{}{
		"id":   "test",
		"type": "unknown",
	}
	_, err := NewVariableFromDict(dict)
	c.Assert(err.Error(), Equals, "Unknown variable type: unknown")
}
