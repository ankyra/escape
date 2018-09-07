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
	. "gopkg.in/check.v1"
)

type semverSuite struct{}

var _ = Suite(&semverSuite{})

func (s *semverSuite) Test_LessOrEqual(c *C) {
	unit := NewSemanticVersion("0.0.3")
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.4.0")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3.1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3.0")), Equals, true)
	c.Assert(NewSemanticVersion("0.0.3.0").LessOrEqual(unit), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.1")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.2")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.3")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.4")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0.5")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.0")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0.1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("0")), Equals, false)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("1")), Equals, true)
	c.Assert(unit.LessOrEqual(NewSemanticVersion("2")), Equals, true)
}

func (s *semverSuite) Test_IncrementSmallest(c *C) {
	unit := NewSemanticVersion("0.0.3")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "0.0.4")
	unit = NewSemanticVersion("0.1")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "0.2")
	unit = NewSemanticVersion("0")
	unit.IncrementSmallest()
	c.Assert(unit.ToString(), Equals, "1")
}

func (s *semverSuite) Test_OnlyKeepLeadingVersionPart(c *C) {
	unit := NewSemanticVersion("0.0.3")
	unit.OnlyKeepLeadingVersionPart()
	c.Assert(unit.ToString(), Equals, "0")
	unit = NewSemanticVersion("10.23.33")
	unit.OnlyKeepLeadingVersionPart()
	c.Assert(unit.ToString(), Equals, "10")
}
