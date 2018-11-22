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

type versionSuite struct{}

var _ = Suite(&versionSuite{})

func (s *versionSuite) Test_ParseVersionQuery_latest(c *C) {
	testCases := []string{"latest", "v@", "@"}
	for _, testCase := range testCases {
		vq, err := ParseVersionQuery(testCase)
		c.Assert(err, IsNil)
		c.Assert(vq.LatestVersion, Equals, true)
		c.Assert(vq.VersionPrefix, Equals, "")
		c.Assert(vq.SpecificVersion, Equals, "")
		c.Assert(vq.SpecificTag, Equals, "")
	}
}

func (s *versionSuite) Test_ParseVersionQuery_without_prefix(c *C) {
	testCases := map[string]string{
		"v0.1.0":         "0.1.0",
		"0.1.0":          "0.1.0",
		"0.1.10.100.1.8": "0.1.10.100.1.8",
	}
	for testCase, expected := range testCases {
		vq, err := ParseVersionQuery(testCase)
		c.Assert(err, IsNil)
		c.Assert(vq.LatestVersion, Equals, false)
		c.Assert(vq.VersionPrefix, Equals, "")
		c.Assert(vq.SpecificVersion, Equals, expected, Commentf("Expecting '%s' for '%s'", expected, testCase))
		c.Assert(vq.SpecificTag, Equals, "")
	}
}

func (s *versionSuite) Test_ParseVersionQuery_with_prefix(c *C) {
	testCases := map[string]string{
		"v0.1.@":         "0.1.",
		"0.1.@":          "0.1.",
		"0.1.10.100.1.@": "0.1.10.100.1.",
	}
	for testCase, expected := range testCases {
		vq, err := ParseVersionQuery(testCase)
		c.Assert(err, IsNil)
		c.Assert(vq.LatestVersion, Equals, false)
		c.Assert(vq.VersionPrefix, Equals, expected)
		c.Assert(vq.SpecificVersion, Equals, "")
		c.Assert(vq.SpecificTag, Equals, "")
	}
}

func (s *versionSuite) Test_ParseVersionQuery_tag(c *C) {
	testCases := []string{"production", "ci"}
	for _, testCase := range testCases {
		vq, err := ParseVersionQuery(testCase)
		c.Assert(err, IsNil)
		c.Assert(vq.LatestVersion, Equals, false)
		c.Assert(vq.VersionPrefix, Equals, "")
		c.Assert(vq.SpecificVersion, Equals, "")
		c.Assert(vq.SpecificTag, Equals, testCase)
	}
}
