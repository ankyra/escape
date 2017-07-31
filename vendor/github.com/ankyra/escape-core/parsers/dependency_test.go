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
)

type dependencySuite struct{}

var _ = Suite(&dependencySuite{})

func (s *dependencySuite) Test_Dependency_Happy_Path1(c *C) {
	dep, err := ParseDependency("name-v1.0")
	c.Assert(err, IsNil)
	c.Assert(dep.Name, Equals, "name")
	c.Assert(dep.Version, Equals, "1.0")
	c.Assert(dep.VariableName, Equals, "")
	c.Assert(dep.Project, Equals, "_")
}

func (s *dependencySuite) Test_Dependency_Happy_Path2(c *C) {
	dep, err := ParseDependency("project/type-name-v1.0 as dep")
	c.Assert(err, IsNil)
	c.Assert(dep.Name, Equals, "type-name")
	c.Assert(dep.Version, Equals, "1.0")
	c.Assert(dep.VariableName, Equals, "dep")
	c.Assert(dep.Project, Equals, "project")
}

func (s *dependencySuite) Test_Dependency_WhiteSpace(c *C) {
	dep, err := ParseDependency("   name-v1.0    as   dep  ")
	c.Assert(err, IsNil)
	c.Assert(dep.Name, Equals, "name")
	c.Assert(dep.Version, Equals, "1.0")
	c.Assert(dep.VariableName, Equals, "dep")
}

func (s *dependencySuite) Test_Dependency_Missing_Id(c *C) {
	dep, err := ParseDependency("name-v1.0 as")
	c.Assert(dep, IsNil)
	c.Assert(err.Error(), Equals, "Malformed dependency string 'name-v1.0 as'")
}

func (s *dependencySuite) Test_Dependency_Second_Word_Not_As(c *C) {
	dep, err := ParseDependency("type-name-v1.0 oh identifier")
	c.Assert(dep, IsNil)
	c.Assert(err.Error(), Equals, "Unexpected 'oh' expecting 'as' in 'type-name-v1.0 oh identifier'")
}

func (s *dependencySuite) Test_Dependency_Malformed_Release_Id(c *C) {
	dep, err := ParseDependency("type-name-whatever")
	c.Assert(dep, IsNil)
	c.Assert(err.Error(), Equals, "Invalid version string in release ID 'type-name-whatever': whatever")
}
