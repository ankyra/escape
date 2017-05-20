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

package core

import (
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
	"testing"
)

type metadataSuite struct{}

var _ = Suite(&metadataSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *metadataSuite) Test_GetReleaseId(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	releaseId := m.GetReleaseId()
	c.Assert(releaseId, Equals, "test-release-v0.1")
}

func (s *metadataSuite) Test_GetVersionlessReleaseId(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	releaseId := m.GetVersionlessReleaseId()
	c.Assert(releaseId, Equals, "test-release")
}

func (s *metadataSuite) Test_VariableContext(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	m.SetVariableInContext("test_key1", "test_value1")
	m.SetVariableInContext("test_key2", "test_value2")
	ctx := m.GetVariableContext()
	c.Assert(ctx, HasLen, 2)
	c.Assert(ctx["test_key1"], Equals, "test_value1")
	c.Assert(ctx["test_key2"], Equals, "test_value2")
}

func (s *metadataSuite) Test_InputVariables(c *C) {
	v1, _ := variables.NewVariableFromString("input_variable1", "string")
	v2, _ := variables.NewVariableFromString("input_variable2", "string")
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddInputVariable(v1)
	m.AddInputVariable(v2)
	vars := m.GetInputs()
	c.Assert(vars, HasLen, 2)
	c.Assert(vars[0], Equals, v1)
	c.Assert(vars[1], Equals, v2)
}

func (s *metadataSuite) Test_OutputVariables(c *C) {
	v1, _ := variables.NewVariableFromString("output_variable1", "string")
	v2, _ := variables.NewVariableFromString("output_variable2", "string")
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddOutputVariable(v1)
	m.AddOutputVariable(v2)
	vars := m.GetOutputs()
	c.Assert(len(vars), Equals, 2)
	c.Assert(vars[0], Equals, v1)
	c.Assert(vars[1], Equals, v2)
}

func (s *metadataSuite) Test_GetDirectores(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddFileWithDigest("test/file1.txt", "abcdef")
	m.AddFileWithDigest("test/file2.txt", "abcdef")
	m.AddFileWithDigest("test2/file3.txt", "abcdef")
	m.AddFileWithDigest("test2/test3/file4.txt", "abcdef")
	dirs := m.GetDirectories()
	c.Assert(len(dirs), Equals, 3)
	expectedDirs := map[string]bool{
		"test/":        false,
		"test2/":       false,
		"test2/test3/": false,
	}
	for _, dir := range dirs {
		alreadySeen, found := expectedDirs[dir]
		c.Assert(found, Equals, true)
		c.Assert(alreadySeen, Equals, false)
		expectedDirs[dir] = true
	}
}

func (s *metadataSuite) Test_FromJson(c *C) {
	json := `{
        "api_version": "1",
        "consumes": ["provider1", "provider2"],
        "name": "test-release",
        "description": "Test release",
        "version": "0.1",
        "variable_context": {
            "base": "test-depends-v1",
            "test-depends": "test-depends-v1"
        }
    }`
	m, err := NewReleaseMetadataFromJsonString(json)
	c.Assert(err, IsNil)
	c.Assert(m.GetApiVersion(), Equals, "1")
	c.Assert(m.GetName(), Equals, "test-release")
	c.Assert(m.GetDescription(), Equals, "Test release")
	c.Assert(m.GetVersion(), Equals, "0.1")
	c.Assert(m.GetConsumes()[0], Equals, "provider1")
	c.Assert(m.GetConsumes()[1], Equals, "provider2")
	c.Assert(m.GetVariableContext()["base"], Equals, "test-depends-v1")
	c.Assert(m.GetVariableContext()["test-depends"], Equals, "test-depends-v1")
}
