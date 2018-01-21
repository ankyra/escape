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

package controllers

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *suite) Test_NewControllerResult(c *C) {
	unit := NewControllerResult()
	c.Assert(unit.HumanOutput, NotNil)
	c.Assert(unit.MarshalableOutput, IsNil)
	c.Assert(unit.Error, IsNil)
}

func (s *suite) Test_NewControllerResult_Print_ReturnsErrorIfSet(c *C) {
	unit := &ControllerResult{
		Error: fmt.Errorf("test error"),
	}

	c.Assert(unit.Print(false), NotNil)
	c.Assert(unit.Print(false), ErrorMatches, "test error")
}

func (s *suite) Test_NewHumanOutput(c *C) {
	unit := NewHumanOutput("%s", "test string")

	c.Assert(unit.value, Equals, "test string")
}

func (s *suite) Test_NewHumanOutput_AcceptsSingleString(c *C) {
	unit := NewHumanOutput("test string")

	c.Assert(unit.value, Equals, "test string")
}

func (s *suite) Test_NewHumanOutput_AddLine(c *C) {
	unit := NewHumanOutput("")
	unit.AddLine("%s", "test string")

	c.Assert(unit.value, Equals, "test string")
}

func (s *suite) Test_NewHumanOutput_AddLine_AcceptsSingleString(c *C) {
	unit := NewHumanOutput("")
	unit.AddLine("test string")

	c.Assert(unit.value, Equals, "test string")
}

func (s *suite) Test_NewHumanOutput_AddLine_AddsNewLineOnMultipleInvocations(c *C) {
	unit := NewHumanOutput("")
	unit.AddLine("%s", "test")
	unit.AddLine("%s", "string")

	c.Assert(unit.value, Equals, "test\nstring")
}

func (s *suite) Test_NewHumanOutput_AddLine_AcceptsSingleString_And_AddsNewLineOnMultipleInvocations(c *C) {
	unit := NewHumanOutput("")
	unit.AddLine("test")
	unit.AddLine("string")

	c.Assert(unit.value, Equals, "test\nstring")
}

func (s *suite) Test_NewHumanOutput_AddMap(c *C) {
	unit := NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test": "value",
	})

	c.Assert(unit.value, Matches, ".*test: value.*")

	unit = NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test": 12,
	})

	c.Assert(unit.value, Matches, ".*test: 12.*")

	unit = NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test": 12.12,
	})

	c.Assert(unit.value, Matches, ".*test: 12.12.*")

	unit = NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test": true,
	})

	c.Assert(unit.value, Matches, ".*test: true.*")

	unit = NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test": []string{"hello", "world"},
	})

	c.Assert(unit.value, Equals, "test: [hello world]")
}

func (s *suite) Test_NewHumanOutput_AddMap_AddsNewLineOnMultipleInvocations(c *C) {
	unit := NewHumanOutput("")
	unit.AddMap(map[string]interface{}{
		"test1": "value",
	})

	unit.AddMap(map[string]interface{}{
		"test2": 15,
	})

	c.Assert(unit.value, Matches, "test1: value\n\ntest2: 15")
}

func (s *suite) Test_NewHumanOutput_AddList(c *C) {
	unit := NewHumanOutput("")
	unit.AddList([]interface{}{"one", "two"})

	c.Assert(unit.value, Equals, "one\ntwo")
}

func (s *suite) Test_NewHumanOutput_AddList_AddsNewLineOnMultipleInvocations(c *C) {
	unit := NewHumanOutput("")
	unit.AddList([]interface{}{"one", "two"})
	unit.AddList([]interface{}{"three", "four"})

	c.Assert(unit.value, Equals, "one\ntwo\n\nthree\nfour")
}
