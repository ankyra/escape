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

type primSuite struct{}

var _ = Suite(&primSuite{})

func (s *primSuite) Test_Parse_Ident1(c *C) {
	result, rest := ParseIdent("test")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "test")
}

func (s *primSuite) Test_Parse_Ident2(c *C) {
	result, rest := ParseIdent("test123")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "test123")
}

func (s *primSuite) Test_Parse_Ident3(c *C) {
	result, rest := ParseIdent("test_123")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "test_123")
}

func (s *primSuite) Test_Parse_Ident4(c *C) {
	result, rest := ParseIdent("test-123")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "test-123")
}

func (s *primSuite) Test_Parse_Ident5(c *C) {
	result, rest := ParseIdent("test=12")
	c.Assert(rest, Equals, "=12")
	c.Assert(result, Equals, "test")
}

func (s *primSuite) Test_Parse_Ident_Empty_String(c *C) {
	result, rest := ParseIdent("")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "")
}

func (s *primSuite) Test_Parse_Ident_Cant_Start_With_Integer(c *C) {
	result, rest := ParseIdent("1test")
	c.Assert(rest, Equals, "1test")
	c.Assert(result, Equals, "")
}

func (s *primSuite) Test_Parse_Ident_Cant_Start_With_Dash(c *C) {
	result, rest := ParseIdent("-test")
	c.Assert(rest, Equals, "-test")
	c.Assert(result, Equals, "")
}

func (s *primSuite) Test_Parse_Ident_Cant_Start_With_Underscore(c *C) {
	result, rest := ParseIdent("_test")
	c.Assert(rest, Equals, "_test")
	c.Assert(result, Equals, "")
}

func (s *primSuite) Test_Parse_Ident_Whitespace_Prefix(c *C) {
	result, rest := ParseIdent("  test")
	c.Assert(rest, Equals, "")
	c.Assert(result, Equals, "test")
}

func (s *primSuite) Test_Parse_Ident_Whitespace_Suffix(c *C) {
	result, rest := ParseIdent("test ")
	c.Assert(rest, Equals, " ")
	c.Assert(result, Equals, "test")
}

func (s *primSuite) Test_Parse_Integer(c *C) {
	result, rest := ParseInteger("12")
	c.Assert(rest, Equals, "")
	c.Assert(*result, Equals, 12)
}

func (s *primSuite) Test_Parse_Integer_Empty_String(c *C) {
	result, rest := ParseInteger("")
	c.Assert(rest, Equals, "")
	c.Assert(result, IsNil)
}

func (s *primSuite) Test_Parse_Integer_Whitespace_Prefix(c *C) {
	result, rest := ParseInteger("    12")
	c.Assert(rest, Equals, "")
	c.Assert(*result, Equals, 12)
}

func (s *primSuite) Test_Parse_Integer_Whitespace_Suffix(c *C) {
	result, rest := ParseInteger("12 ")
	c.Assert(rest, Equals, " ")
	c.Assert(*result, Equals, 12)
}

func (s *primSuite) Test_Parse_Integer_Big_Int(c *C) {
	result, rest := ParseInteger("123456789012345678901234567890")
	c.Assert(rest, Equals, "123456789012345678901234567890")
	c.Assert(result, IsNil) // TODO: should result in error
}
