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

package parsers

import (
	. "gopkg.in/check.v1"
)

type optionsSuite struct{}

var _ = Suite(&optionsSuite{})

func (s *optionsSuite) Test_Parse_Options1(c *C) {
	opts, err := ParseOptions("key = 12")
	c.Assert(err, IsNil)
	c.Assert(opts["key"], Equals, 12)
}

func (s *optionsSuite) Test_Parse_Options_List(c *C) {
	opts, err := ParseOptions("key = 12, key2 = 640")
	c.Assert(err, IsNil)
	c.Assert(opts["key"], Equals, 12)
	c.Assert(opts["key2"], Equals, 640)
}

func (s *optionsSuite) Test_Parse_Options_Missing_Expression(c *C) {
	opts, err := ParseOptions("key = ")
	c.Assert(err.Error(), Equals, "Expecting expression key=value, got: 'key = '")
	c.Assert(opts, IsNil)
}

func (s *optionsSuite) Test_Parse_Options_Missing_Expression_In_Second_Item(c *C) {
	opts, err := ParseOptions("key = 12, key2 =")
	c.Assert(err.Error(), Equals, "Expecting expression key=value, got: 'key2 ='")
	c.Assert(opts, IsNil)
}

func (s *optionsSuite) Test_Parse_Operator(c *C) {
	result, rest := parseOperator("=")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "=")
}

func (s *optionsSuite) Test_Parse_Operator2(c *C) {
	result, rest := parseOperator("=12")
	c.Assert(rest, Equals, "12")
	c.Assert(result, Equals, "=")
}

func (s *optionsSuite) Test_Parse_Operator_Whitespace_Prefix(c *C) {
	result, rest := parseOperator("   =")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "=")
}

func (s *optionsSuite) Test_Parse_Operator_Whitespace_Suffix(c *C) {
	result, rest := parseOperator("= ")
	c.Assert(rest, Equals, " ")
	c.Assert(result, Equals, "=")
}

func (s *optionsSuite) Test_Parse_Boolean_Expression(c *C) {
	result, rest := parseExpression("test")
	c.Assert(rest, Equals, "")
	c.Assert(result["test"], Equals, true)
}

func (s *optionsSuite) Test_Parse_Expression_Missing_Value(c *C) {
	result, rest := parseExpression("test=")
	c.Assert(rest, Equals, "test=")
	c.Assert(result, IsNil)
}

func (s *optionsSuite) Test_Parse_Expression1(c *C) {
	result, rest := parseExpression("test=1")
	c.Assert(rest, Equals, "")
	c.Assert(result["test"], Equals, 1)
}
func (s *optionsSuite) Test_Parse_Expression2(c *C) {
	result, rest := parseExpression("test=1, test=12")
	c.Assert(rest, Equals, ", test=12")
	c.Assert(result["test"], Equals, 1)
}
func (s *optionsSuite) Test_Parse_Expression_Empty_String(c *C) {
	result, rest := parseExpression("")
	c.Assert(rest, Equals, "")
	c.Assert(result, IsNil)
}
