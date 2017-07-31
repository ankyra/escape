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

package variable_types

import (
	. "gopkg.in/check.v1"
	"testing"
)

type variableSuite struct{}

var _ = Suite(&variableSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *variableSuite) Test_ValidateList_empty_list(c *C) {
	lst, err := validateList([]interface{}{}, map[string]interface{}{})
	c.Assert(err, IsNil)
	c.Assert(lst, HasLen, 0)
}

func (s *variableSuite) Test_ValidateList_list(c *C) {
	lst, err := validateList([]interface{}{"test", "test2"}, map[string]interface{}{})
	c.Assert(err, IsNil)
	c.Assert(lst, HasLen, 2)
	c.Assert(lst, DeepEquals, []interface{}{"test", "test2"})
}

func (s *variableSuite) Test_ValidateList_empty_string(c *C) {
	lst, err := validateList("", map[string]interface{}{})
	c.Assert(err, IsNil)
	c.Assert(lst, HasLen, 0)
}

func (s *variableSuite) Test_ValidateList_json_string(c *C) {
	lst, err := validateList("[\"test\", \"test2\"]", map[string]interface{}{})
	c.Assert(err, IsNil)
	c.Assert(lst, HasLen, 2)
	c.Assert(lst, DeepEquals, []interface{}{"test", "test2"})
}
