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
	"strings"

	. "gopkg.in/check.v1"
)

type variableSuite struct{}

var _ = Suite(&variableSuite{})

func (s *variableSuite) Test_ParseVariableIdent_happy_path(c *C) {
	cases := []string{
		"id",
		"test",
		"whatever",
		" v1",
		"  v100  ",
		"abc  ",
		"test_test",
		"test-test",
		" my-variable ",
	}
	for _, test := range cases {
		v, err := ParseVariableIdent(test)
		c.Assert(err, IsNil)
		c.Assert(v, Equals, strings.TrimSpace(test))
	}
}

func (s *variableSuite) Test_ParseVariableIdent_fails_on_empty_string(c *C) {
	cases := []string{
		"",
	}
	for _, test := range cases {
		_, err := ParseVariableIdent(test)
		c.Assert(err, DeepEquals, InvalidVariableIdEmptyError)
	}
}

func (s *variableSuite) Test_ParseVariableIdent_fails_on_trailing_characters(c *C) {
	cases := []string{
		"id  1",
		"test   t",
		"whatever  $test",
		" v1  v2 ",
		"  v100  as 3",
		"abc  sd ",
		"test_test test-test -",
		"test-test - 3",
		" my-variable  + 34",
	}
	for _, test := range cases {
		_, err := ParseVariableIdent(test)
		c.Assert(err, DeepEquals, InvalidVariableIdFormatError(test))
	}
}

func (s *variableSuite) Test_ParseVariableIdent_fails_if_starting_with_PREVIOUS_(c *C) {
	cases := []string{
		"PREVIOUS_id",
		"PREVIOUS_test",
		"PREVIOUS_whatever",
		" PREVIOUS_v1",
		"  PREVIOUS_v100  ",
		"PREVIOUS_abc  ",
		"PREVIOUS_test_test",
		"PREVIOUS_test-test",
		" PREVIOUS_my-variable ",
	}
	for _, test := range cases {
		_, err := ParseVariableIdent(test)
		c.Assert(err, DeepEquals, InvalidVariableIdPreviousError(test))
	}
}
