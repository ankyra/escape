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

package escape_plan

import (
	. "gopkg.in/check.v1"
)

type docSuite struct{}

var _ = Suite(&docSuite{})

func (s *docSuite) Test_Doc(c *C) {
	doc := GetDoc("name")
	c.Assert(string(doc), Equals, nameDoc)
	doc = GetDoc("version")
	c.Assert(string(doc), Equals, versionDoc)
	doc = GetDoc("depends")
	c.Assert(string(doc), Equals, dependsDoc)
	doc = GetDoc("includes")
	c.Assert(string(doc), Equals, includesDoc)
	doc = GetDoc("consumes")
	c.Assert(string(doc), Equals, consumesDoc)
	doc = GetDoc("provides")
	c.Assert(string(doc), Equals, providesDoc)
	doc = GetDoc("metadata")
	c.Assert(string(doc), Equals, metadataDoc)
	doc = GetDoc("inputs")
	c.Assert(string(doc), Equals, inputsDoc)
	doc = GetDoc("outputs")
	c.Assert(string(doc), Equals, outputsDoc)
	doc = GetDoc("errands")
	c.Assert(string(doc), Equals, errandsDoc)
	doc = GetDoc("pre_build")
	c.Assert(string(doc), Equals, prebuildDoc)
	doc = GetDoc("post_build")
	c.Assert(string(doc), Equals, postbuildDoc)
}

func (s *docSuite) Test_Doc_doesnt_exist(c *C) {
	doc := GetDoc("")
	c.Assert(doc, DeepEquals, []byte{})
}
