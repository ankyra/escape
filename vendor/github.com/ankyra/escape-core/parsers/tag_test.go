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

package parsers

import (
	. "gopkg.in/check.v1"
)

type tagSuite struct{}

var _ = Suite(&tagSuite{})

func (s *tagSuite) Test_IsValidTag_invalid_cases(c *C) {
	testCases := []string{"", "latest", "v@", "@", "v1", "v11", "v0.1", "v10.11", "v0.1.@", "v1.@", "0", "0.1", "0.0.1", "0.@"}
	for _, testCase := range testCases {
		c.Assert(IsValidTag(testCase), Equals, false, Commentf("The tag '%s' should be invalid", testCase))
	}
}

func (s *tagSuite) Test_IsValidTag_valid_cases(c *C) {
	testCases := []string{"ci", "production", "very", "v10.alpha"}
	for _, testCase := range testCases {
		c.Assert(IsValidTag(testCase), Equals, true, Commentf("The tag '%s' should be valid", testCase))
	}
}
