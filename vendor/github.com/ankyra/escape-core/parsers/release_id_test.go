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
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type releaseIdSuite struct{}

var _ = Suite(&releaseIdSuite{})

func (s *releaseIdSuite) Test_ReleaseId_Happy_Path(c *C) {
	id, err := ParseReleaseId("name-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Name, Equals, "name")
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Can_Have_Dashes(c *C) {
	id, err := ParseReleaseId("name-with-dashes-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Name, Equals, "name-with-dashes")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest1(c *C) {
	id, err := ParseReleaseId("type-name-latest")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest2(c *C) {
	id, err := ParseReleaseId("name-@")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest3(c *C) {
	id, err := ParseReleaseId("name-v@")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Version(c *C) {
	id, err := ParseReleaseId("name-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Tag(c *C) {
	id, err := ParseReleaseId("name:tag")
	c.Assert(err, IsNil)
	c.Assert(id.Tag, Equals, "tag")
}

func (s *releaseIdSuite) Test_NeedsResolving_true(c *C) {
	cases := []string{
		"name-latest",
		"name-v@",
		"name-v0.@",
		"name-v0.0.0.@",
		"name:tag",
	}
	for _, test := range cases {
		id, err := ParseReleaseId(test)
		c.Assert(err, IsNil)
		c.Assert(id.NeedsResolving(), Equals, true)
	}
}

func (s *releaseIdSuite) Test_NeedsResolving_false(c *C) {
	cases := []string{
		"name-v1",
		"name-v1.0",
		"name-v0.0.0",
	}
	for _, test := range cases {
		id, err := ParseReleaseId(test)
		c.Assert(err, IsNil)
		c.Assert(id.NeedsResolving(), Equals, false)
	}
}

func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format1(c *C) {
	_, err := ParseReleaseId("type")
	c.Assert(err, DeepEquals, InvalidReleaseFormatError("type"))
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format2(c *C) {
	_, err := ParseReleaseId("type-name")
	c.Assert(err, DeepEquals, InvalidVersionStringInReleaseIdError("type-name", "name"))
}
func (s *releaseIdSuite) Test_ReleaseId_Missing_Version(c *C) {
	_, err := ParseReleaseId("type-name-nope")
	c.Assert(err, DeepEquals, InvalidVersionStringInReleaseIdError("type-name-nope", "nope"))
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Version(c *C) {
	_, err := ParseReleaseId("type-name-vnope")
	c.Assert(err, DeepEquals, InvalidVersionStringInReleaseIdError("type-name-vnope", "vnope"))
}

func (s *releaseIdSuite) Test_QualifiedReleaseID(c *C) {
	q, err := ParseQualifiedReleaseId("project/type-name-v1")
	c.Assert(err, IsNil)
	c.Assert(q.Project, Equals, "project")
	c.Assert(q.ToString(), Equals, "project/type-name-v1")
}

func (s *releaseIdSuite) Test_QualifiedReleaseID_with_tag(c *C) {
	q, err := ParseQualifiedReleaseId("project/type-name:tag")
	c.Assert(err, IsNil)
	c.Assert(q.Project, Equals, "project")
	c.Assert(q.ToString(), Equals, "project/type-name:tag")
}

func (s *releaseIdSuite) Test_QualifiedReleaseID_default_project(c *C) {
	q, err := ParseQualifiedReleaseId("type-name-v1")
	c.Assert(err, IsNil)
	c.Assert(q.Project, Equals, "_")
	c.Assert(q.ToString(), Equals, "_/type-name-v1")
}

func (s *releaseIdSuite) Test_QualifiedReleaseID_fails_on_invalid_input(c *C) {
	cases := []string{
		"",
		"project",
		"project/type-name-vnope",
	}
	for _, test := range cases {
		_, err := ParseQualifiedReleaseId(test)
		c.Assert(err, Not(IsNil))
	}
}

func (s *releaseIdSuite) Test_IsValidVersion(c *C) {
	cases := []string{
		"latest",
		"0",
		"10",
		"0.0",
		"0.10",
		"0.0.0",
		"0.0.10",
		"0.@",
		"0.0.@",
	}
	for _, test := range cases {
		c.Assert(isValidVersion(test), Equals, true)
	}
}

func (s *releaseIdSuite) Test_IsValidVersion_false(c *C) {
	cases := []string{
		"",
		"whatsthisnow",
		"nope",
		"0.test",
		"0.0.test",
		"0.0.latest",
		"0-0",
		"0_0",
		"0@",
		"0.0@",
	}
	for _, test := range cases {
		c.Assert(isValidVersion(test), Equals, false)
	}
}
