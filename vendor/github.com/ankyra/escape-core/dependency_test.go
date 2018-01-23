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

package core

import (
	"github.com/ankyra/escape-core/parsers"
	. "gopkg.in/check.v1"
)

func (s *metadataSuite) Test_NewDependencyFromMetadata(c *C) {
	metadata := NewReleaseMetadata("test", "1.0")
	dep := NewDependencyFromMetadata(metadata)
	c.Assert(dep.Name, Equals, "test")
	c.Assert(dep.Version, Equals, "1.0")
	c.Assert(dep.Project, Equals, "_")
}

func (s *metadataSuite) Test_NewDependencyFromMetadata_adds_project(c *C) {
	metadata := NewReleaseMetadata("test", "1.0")
	metadata.Project = "prj"
	dep := NewDependencyFromMetadata(metadata)
	c.Assert(dep.Name, Equals, "test")
	c.Assert(dep.Version, Equals, "1.0")
	c.Assert(dep.Project, Equals, "prj")
}

func (s *metadataSuite) Test_NewDependencyFromString(c *C) {
	testCases := map[string][]string{
		"prj/test-v1.0":                  []string{"prj", "test", "1.0", ""},
		"test-v1.0":                      []string{"_", "test", "1.0", ""},
		"test-latest":                    []string{"_", "test", "latest", ""},
		"test-v@":                        []string{"_", "test", "latest", ""},
		"test-v0.@":                      []string{"_", "test", "0.@", ""},
		"prj/test-v0.@ as var":           []string{"prj", "test", "0.@", "var"},
		"   prj/test-v0.@   as    var  ": []string{"prj", "test", "0.@", "var"},
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.Project, Equals, expected[0])
		c.Assert(dep.Name, Equals, expected[1])
		c.Assert(dep.Version, Equals, expected[2])
		c.Assert(dep.VariableName, Equals, expected[3])
	}
}

func (s *metadataSuite) Test_NewDependencyFromString_invalid(c *C) {
	testCases := []string{
		"",
		"   ",
		"wfiwpef piwfje pwaeifo kwae",
		"test-v1.0 as",
		"tes v.0",
	}
	for _, testCase := range testCases {
		_, err := NewDependencyFromString(testCase)
		c.Assert(err, Not(IsNil))
	}
}

func (s *metadataSuite) Test_NewDependencyFromQualifiedReleaseId(c *C) {
	testCases := map[string][]string{
		"prj/test-v1.0": []string{"prj", "test", "1.0", ""},
		"test-v1.0":     []string{"_", "test", "1.0", ""},
		"test-latest":   []string{"_", "test", "latest", ""},
		"test-v@":       []string{"_", "test", "latest", ""},
		"test-v0.@":     []string{"_", "test", "0.@", ""},
	}
	for testCase, expected := range testCases {
		release, err := parsers.ParseQualifiedReleaseId(testCase)
		c.Assert(err, IsNil)
		dep := NewDependencyFromQualifiedReleaseId(release)
		c.Assert(dep.Project, Equals, expected[0])
		c.Assert(dep.Name, Equals, expected[1])
		c.Assert(dep.Version, Equals, expected[2])
		c.Assert(dep.VariableName, Equals, expected[3])
	}
}

func (s *metadataSuite) Test_GetVersionAsString(c *C) {
	testCases := map[string]string{
		"prj/test-v1.0":                  "v1.0",
		"test-v1.0":                      "v1.0",
		"test-latest":                    "latest",
		"test-v@":                        "latest",
		"test-v0.@":                      "v0.@",
		"prj/test-v0.@ as var":           "v0.@",
		"   prj/test-v0.@   as    var  ": "v0.@",
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.GetVersionAsString(), Equals, expected)
	}
}

func (s *metadataSuite) Test_GetDependencyReleaseId(c *C) {
	testCases := map[string]string{
		"prj/test-v1.0":                  "test-v1.0",
		"test-v1.0":                      "test-v1.0",
		"test-latest":                    "test-latest",
		"test-v@":                        "test-latest",
		"test-v0.@":                      "test-v0.@",
		"prj/test-v0.@ as var":           "test-v0.@",
		"   prj/test-v0.@   as    var  ": "test-v0.@",
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.GetReleaseId(), Equals, expected)
	}
}

func (s *metadataSuite) Test_GetDependencyQualifiedReleaseId(c *C) {
	testCases := map[string]string{
		"prj/test-v1.0":                  "prj/test-v1.0",
		"test-v1.0":                      "_/test-v1.0",
		"test-latest":                    "_/test-latest",
		"test-v@":                        "_/test-latest",
		"test-v0.@":                      "_/test-v0.@",
		"prj/test-v0.@ as var":           "prj/test-v0.@",
		"   prj/test-v0.@   as    var  ": "prj/test-v0.@",
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.GetQualifiedReleaseId(), Equals, expected)
	}
}

func (s *metadataSuite) Test_GetDependencyVersionlessReleaseId(c *C) {
	testCases := map[string]string{
		"prj/test-v1.0":                  "prj/test",
		"test-v1.0":                      "_/test",
		"test-latest":                    "_/test",
		"test-v@":                        "_/test",
		"test-v0.@":                      "_/test",
		"prj/test-v0.@ as var":           "prj/test",
		"   prj/test-v0.@   as    var  ": "prj/test",
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.GetVersionlessReleaseId(), Equals, expected)
	}
}

func (s *metadataSuite) Test_NeedsResolving(c *C) {
	testCases := map[string]bool{
		"prj/test-v1.0":                  false,
		"test-v1.0":                      false,
		"test-latest":                    true,
		"test-v@":                        true,
		"test-v0.@":                      true,
		"prj/test-v0.@ as var":           true,
		"   prj/test-v0.@   as    var  ": true,
	}
	for testCase, expected := range testCases {
		dep, err := NewDependencyFromString(testCase)
		c.Assert(err, IsNil)
		c.Assert(dep.NeedsResolving(), Equals, expected)
	}
}
