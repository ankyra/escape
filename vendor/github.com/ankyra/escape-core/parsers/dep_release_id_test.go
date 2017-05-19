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

func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format1(c *C) {
	_, err := ParseReleaseId("type")
	c.Assert(err.Error(), Equals, "Invalid release format: type")
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Format2(c *C) {
	_, err := ParseReleaseId("type-name")
	c.Assert(err.Error(), Equals, "Invalid version string in release ID 'type-name': name")
}
func (s *releaseIdSuite) Test_ReleaseId_Missing_Version(c *C) {
	_, err := ParseReleaseId("type-name-nope")
	c.Assert(err.Error(), Equals, "Invalid version string in release ID 'type-name-nope': nope")
}
func (s *releaseIdSuite) Test_ReleaseId_Invalid_Version(c *C) {
	_, err := ParseReleaseId("type-name-vnope")
	c.Assert(err.Error(), Equals, "Invalid release ID 'type-name-vnope': Invalid version format: nope")
}

func (s *releaseIdSuite) Test_ValidateVersion(c *C) {
	c.Assert(ValidateVersion("latest"), IsNil)
	c.Assert(ValidateVersion("0"), IsNil)
	c.Assert(ValidateVersion("10"), IsNil)
	c.Assert(ValidateVersion("0.0"), IsNil)
	c.Assert(ValidateVersion("0.10"), IsNil)
	c.Assert(ValidateVersion("0.0.0"), IsNil)
	c.Assert(ValidateVersion("0.0.10"), IsNil)
	c.Assert(ValidateVersion("0.@"), IsNil)
	c.Assert(ValidateVersion("0.0.@"), IsNil)
}

func (s *releaseIdSuite) Test_ValidateVersion_Error(c *C) {
	c.Assert(ValidateVersion("whatsthisnow"), Not(IsNil))
	c.Assert(ValidateVersion("nope"), Not(IsNil))
	c.Assert(ValidateVersion("0.test"), Not(IsNil))
	c.Assert(ValidateVersion("0.0.test"), Not(IsNil))
	c.Assert(ValidateVersion("0.0.latest"), Not(IsNil))
	c.Assert(ValidateVersion("0-0"), Not(IsNil))
	c.Assert(ValidateVersion("0_0"), Not(IsNil))
	c.Assert(ValidateVersion("0@"), Not(IsNil))
	c.Assert(ValidateVersion("0.0@"), Not(IsNil))
}
