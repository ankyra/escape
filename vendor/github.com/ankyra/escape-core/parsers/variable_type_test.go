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
	"github.com/ankyra/escape-core/variables/variable_types"
	. "gopkg.in/check.v1"
)

type varTypSuite struct{}

var _ = Suite(&varTypSuite{})

func (s *varTypSuite) Test_Parse_VariableType1(c *C) {
	for _, typ := range variable_types.GetSupportedTypes() {
		t, err := ParseVariableType(typ)
		c.Assert(err, IsNil)
		c.Assert(t.Type, Equals, typ)
	}
}

func (s *varTypSuite) Test_Parse_VariableType2(c *C) {
	t, err := ParseVariableType("string[min=10,max=14]")
	c.Assert(err, IsNil)
	c.Assert(t.Type, Equals, "string")
	c.Assert(t.Options["min"], Equals, 10)
	c.Assert(t.Options["max"], Equals, 14)
}
