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

package variable_types

import (
	. "gopkg.in/check.v1"
)

func (s *variableSuite) Test_ValidateInt(c *C) {
	testCases := map[interface{}]int{
		0:       0,
		1:       1,
		1000:    1000,
		-1000:   -1000,
		0.0:     0,
		1.0:     1,
		1000.0:  1000,
		-1000.0: -1000,
		"0":     0,
		"1":     1,
		"1000":  1000,
		"-1000": -1000,
	}
	for testCase, expected := range testCases {
		result, err := validateInt(testCase, nil)
		c.Assert(err, IsNil)
		c.Assert(result, Equals, expected, Commentf("'%v' should be '%v'", testCase, expected))
	}
}
