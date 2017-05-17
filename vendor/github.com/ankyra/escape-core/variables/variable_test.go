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
	. "gopkg.in/check.v1"
	"testing"
)

type variableSuite struct{}

var _ = Suite(&variableSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *variableSuite) Test_GetValue_String_Variable(c *C) {
	unit := NewVariableFromString("test", "string")
	variableCtx := map[string]interface{}{
		"test": "test value",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "test value")
}

func (s *variableSuite) Test_GetValue_Uses_Default(c *C) {
	unit := NewVariableFromString("test", "string")
	defaultVal := "test value"
	unit.SetDefault(&defaultVal)
	val, err := unit.GetValue(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "test value")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable(c *C) {
	unit := NewVariableFromString("test", "string")
	unit.SetOneOfItems([]interface{}{"valid", "also valid"})
	variableCtx := map[string]interface{}{
		"test": "valid",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "valid")
}

func (s *variableSuite) Test_GetValue_OneOf_Variable_Fails(c *C) {
	unit := NewVariableFromString("test", "string")
	unit.SetOneOfItems([]interface{}{"valid", "also valid"})
	variableCtx := map[string]interface{}{
		"test": "not valid",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), DeepEquals, "Expecting one of [\"valid\",\"also valid\"] for variable 'test'")
}

func (s *variableSuite) Test_String_Variable_Converts_To_String_Value(c *C) {
	unit := NewVariableFromString("test", "string")
	variableCtx := map[string]interface{}{
		"test": 12,
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, "12")
}

func (s *variableSuite) Test_GetValue_Integer_Variable(c *C) {
	unit := NewVariableFromString("test", "integer")
	variableCtx := map[string]interface{}{
		"test": 12,
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, 12)
}

func (s *variableSuite) Test_Integer_Variable_Expects_Integer_Value(c *C) {
	unit := NewVariableFromString("test", "integer")
	variableCtx := map[string]interface{}{
		"test": "test",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), Not(Equals), "")
}

func (s *variableSuite) Test_Integer_Variable_Expects_Integer_Value_Or_Convertable_String(c *C) {
	unit := NewVariableFromString("test", "integer")
	variableCtx := map[string]interface{}{
		"test": "12",
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, 12)
}

func (s *variableSuite) Test_GetValue_List_Variable(c *C) {
	unit := NewVariableFromString("test", "list")
	variableCtx := map[string]interface{}{
		"test": []interface{}{"test value", "test value 2"},
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, []interface{}{"test value", "test value 2"})
}

func (s *variableSuite) Test_GetValue_List_Variable_Checks_String_Values(c *C) {
	unit := NewVariableFromString("test", "list")
	variableCtx := map[string]interface{}{
		"test": []interface{}{12},
	}
	val, err := unit.GetValue(&variableCtx, nil)
	c.Assert(val, IsNil)
	c.Assert(err.Error(), Equals, "Unexpected 'integer' value in list, expecting 'string' for variable 'test'")
}
