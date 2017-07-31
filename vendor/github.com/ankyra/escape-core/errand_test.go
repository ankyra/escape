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

package core

import (
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

func (s *metadataSuite) Test_Diff_NewErrand(c *C) {
	var1, _ := variables.NewVariableFromString("var1", "string")
	testCases := [][]interface{}{
		[]interface{}{map[interface{}]interface{}{
			"description": "Description",
			"script":      "test.sh",
		}, Errand{Description: "Description", Script: "test.sh", Inputs: []*variables.Variable{}}},
		[]interface{}{map[interface{}]interface{}{
			"description": "Description",
			"script":      "test.sh",
			"inputs":      []interface{}{"var1"},
		}, Errand{Description: "Description", Script: "test.sh", Inputs: []*variables.Variable{var1}}},
		[]interface{}{map[interface{}]interface{}{
			"description": "Description",
			"script":      "test.sh",
			"inputs":      []interface{}{map[interface{}]interface{}{"id": "var1"}},
		}, Errand{Description: "Description", Script: "test.sh", Inputs: []*variables.Variable{var1}}},
	}
	for _, test := range testCases {
		errand, err := NewErrandFromDict("name", test[0])
		c.Assert(err, IsNil)
		expected := test[1].(Errand)
		c.Assert(errand.Name, Equals, "name")
		c.Assert(errand.Description, Equals, expected.Description)
		c.Assert(errand.Script, Equals, expected.Script)
		c.Assert(errand.Inputs, DeepEquals, expected.Inputs)
	}
}

func (s *metadataSuite) Test_Diff_NewErrand_invalid(c *C) {
	testCases := [][]interface{}{
		[]interface{}{"", map[interface{}]interface{}{}, "Missing name in errand"},
		[]interface{}{"name", true, "Expecting a dictionary for errand name"},
		[]interface{}{"name", map[interface{}]interface{}{}, "Missing script in errand 'name'"},
		[]interface{}{"name", map[interface{}]interface{}{
			true: true,
		}, "Expecting string key for errand name"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script": true,
		}, "Expecting string value for script field in errand name"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script":      "test.sh",
			"description": true,
		}, "Expecting string value for description field in errand name"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script": "test.sh",
			"inputs": []interface{}{"$invalid-var"},
		}, "Not a valid errand variable: $invalid-var"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script": "test.sh",
			"inputs": []interface{}{map[interface{}]interface{}{"id": "$invalid-var"}},
		}, "Invalid variable format '$invalid-var'"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script": "test.sh",
			"inputs": true,
		}, "Expecting list type for inputs key in errand name"},
		[]interface{}{"name", map[interface{}]interface{}{
			"script": "test.sh",
			"inputs": []interface{}{true},
		}, "Expecting dict or string type for input item in errand name"},
	}
	for _, test := range testCases {
		_, err := NewErrandFromDict(test[0].(string), test[1])
		c.Assert(err, Not(IsNil))
		c.Assert(err.Error(), Equals, test[2].(string))
	}
}

func (s *metadataSuite) Test_GetInputs(c *C) {
	var1, _ := variables.NewVariableFromString("var1", "string")
	errand := NewErrand("test", "script.sh", "description")
	errand.Inputs = append(errand.Inputs, var1)
	c.Assert(errand.GetInputs(), DeepEquals, []*variables.Variable{var1})
}

func (s *metadataSuite) Test_Validate_Inputs(c *C) {
	var1, _ := variables.NewVariableFromString("var1", "string")
	errand := NewErrand("test", "script.sh", "description")
	errand.Inputs = append(errand.Inputs, var1)
	var1.Id = "$test"
	c.Assert(errand.Validate().Error(), Equals, "Error in errand 'test' variable: Invalid variable format '$test'")
	errand.Inputs = nil
	c.Assert(errand.Validate(), IsNil)
}
