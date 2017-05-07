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
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type releaseIdSuite struct{}

var _ = Suite(&releaseIdSuite{})

func (s *releaseIdSuite) Test_ReleaseId_Happy_Path(c *C) {
	id, err := ParseReleaseId("type-name-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Type, Equals, "type")
	c.Assert(id.Build, Equals, "name")
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Can_Have_Dashes(c *C) {
	id, err := ParseReleaseId("type-name-with-dashes-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Build, Equals, "name-with-dashes")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest1(c *C) {
	id, err := ParseReleaseId("type-name-latest")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest2(c *C) {
	id, err := ParseReleaseId("type-name-@")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Latest3(c *C) {
	id, err := ParseReleaseId("type-name-v@")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "latest")
}

func (s *releaseIdSuite) Test_ReleaseId_Parse_Version(c *C) {
	id, err := ParseReleaseId("type-name-v1.0")
	c.Assert(err, IsNil)
	c.Assert(id.Version, Equals, "1.0")
}

func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format1(c *C) {
	_, err := ParseReleaseId("type")
	c.Assert(err.Error(), Equals, "Invalid release format: type")
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format2(c *C) {
	_, err := ParseReleaseId("type-name")
	c.Assert(err.Error(), Equals, "Invalid release format: type-name")
}
func (s *releaseIdSuite) Test_ReleaseId_Missing_Version(c *C) {
	_, err := ParseReleaseId("type-name-nope")
	c.Assert(err.Error(), Equals, "Invalid version string in release ID 'type-name-nope': nope")
}
