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
	"encoding/json"
	"strconv"
	"testing"

	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

type metadataSuite struct{}

var _ = Suite(&metadataSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *metadataSuite) Test_NewReleaseMetadata(c *C) {
	m := NewReleaseMetadata("test", "1.0")
	c.Assert(m.Name, Equals, "test")
	c.Assert(m.Version, Equals, "1.0")
	c.Assert(m.Project, Equals, "_")
	c.Assert(m.BuiltWithCoreVersion, Equals, CoreVersion)
}

func (s *metadataSuite) Test_NewReleaseMetadata_name_check(c *C) {
	testCases := map[string]bool{
		"valid":                                   true,
		"valid-name":                              true,
		"valid_name":                              true,
		"valid1":                                  true,
		"valid-name-2":                            true,
		"  invalid":                               false,
		"_invalid":                                false,
		"invalid ":                                false,
		"":                                        false,
		"$":                                       false,
		"a$":                                      false,
		"1invalid-start":                          false,
		"project/project-should-have-been-parsed": false,

		// protected
		"this":    false,
		"string":  false,
		"integer": false,
		"list":    false,
		"dict":    false,
		"func":    false,
	}
	for testCase, expected := range testCases {
		obj := map[string]string{
			"name":    testCase,
			"version": "1.0",
		}
		payload, err := json.Marshal(obj)
		c.Assert(err, IsNil)
		m, err := NewReleaseMetadataFromJsonString(string(payload))
		if expected {
			c.Assert(err, IsNil)
			c.Assert(m.Name, Equals, testCase)
		} else {
			c.Assert(err, Not(IsNil), Commentf("'%s' is not a valid name", testCase))
		}
	}
}

func (s *metadataSuite) Test_NewReleaseMetadata_project_check(c *C) {
	testCases := map[string]bool{
		"valid":                                   true,
		"valid-name":                              true,
		"valid_name":                              true,
		"valid1":                                  true,
		"valid-name-2":                            true,
		"  invalid":                               false,
		"_invalid":                                false,
		"invalid ":                                false,
		"$":                                       false,
		"a$":                                      false,
		"1invalid-start":                          false,
		"project/project-should-have-been-parsed": false,

		// protected
		"this":    false,
		"string":  false,
		"integer": false,
		"list":    false,
		"dict":    false,
	}
	for testCase, expected := range testCases {
		obj := map[string]string{
			"name":    "name",
			"version": "1.0",
			"project": testCase,
		}
		payload, err := json.Marshal(obj)
		c.Assert(err, IsNil)
		m, err := NewReleaseMetadataFromJsonString(string(payload))
		if expected {
			c.Assert(err, IsNil)
			c.Assert(m.Project, Equals, testCase)
		} else {
			c.Assert(err, Not(IsNil), Commentf("'%s' is not a valid project name", testCase))
		}
	}
}

func (s *metadataSuite) Test_validate(c *C) {
	testCases := map[string]string{
		`null`:                                                      "Missing release metadata",
		`{}`:                                                        "Missing name field in release metadata",
		`{"name": "1"}`:                                             "Invalid name '1'",
		`{"name": "test"}`:                                          "Missing version field in release metadata",
		`{"name": "test", "version": "1", "api_version": 1000}`:     "The release metadata is compiled with a version of Escape targetting API version v1000, but this build supports up to v" + strconv.Itoa(CurrentApiVersion),
		`{"name": "name", "version": "@ASD"}`:                       "Invalid version string '@ASD'.",
		`{"name": "name", "version": "1", "inputs": [{"id": ""}]}`:  "Variable object is missing an 'id'",
		`{"name": "name", "version": "1", "outputs": [{"id": ""}]}`: "Variable object is missing an 'id'",
	}
	for testCase, expected := range testCases {
		_, err := NewReleaseMetadataFromJsonString(testCase)
		c.Assert(err.Error(), Equals, expected)
	}
}

func (s *metadataSuite) Test_GetReleaseId(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	releaseId := m.GetReleaseId()
	c.Assert(releaseId, Equals, "test-release-v0.1")
}

func (s *metadataSuite) Test_GetVersionlessReleaseId(c *C) {
	m := NewReleaseMetadata("test-release", "0.1")
	releaseId := m.GetVersionlessReleaseId()
	c.Assert(releaseId, Equals, "_/test-release")
}

func (s *metadataSuite) Test_InputVariables(c *C) {
	v1, _ := variables.NewVariableFromString("input_variable1", "string")
	v2, _ := variables.NewVariableFromString("input_variable2", "string")
	v2.Scopes = []string{"deploy"}
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddInputVariable(v1)
	m.AddInputVariable(v2)
	vars := m.GetInputs("deploy")
	c.Assert(vars, HasLen, 2)
	c.Assert(vars[0], Equals, v1)
	c.Assert(vars[1], Equals, v2)
	vars = m.GetInputs("build")
	c.Assert(vars, HasLen, 1)
	c.Assert(vars[0], Equals, v1)
}

func (s *metadataSuite) Test_GetInputsInScopes(c *C) {
	v1, _ := variables.NewVariableFromString("input_variable1", "string")
	v2, _ := variables.NewVariableFromString("input_variable2", "string")
	v2.Scopes = []string{"deploy"}
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddInputVariable(v1)
	m.AddInputVariable(v2)
	vars := m.GetInputsInScopes([]string{"deploy"})
	c.Assert(vars, HasLen, 2)
	c.Assert(vars[0], Equals, v1)
	c.Assert(vars[1], Equals, v2)
	vars = m.GetInputsInScopes([]string{"build"})
	c.Assert(vars, HasLen, 1)
	c.Assert(vars[0], Equals, v1)
	vars = m.GetInputsInScopes([]string{"build", "deploy"})
	c.Assert(vars, HasLen, 1)
	c.Assert(vars[0], Equals, v1)
}

func (s *metadataSuite) Test_OutputVariables(c *C) {
	v1, _ := variables.NewVariableFromString("output_variable1", "string")
	v2, _ := variables.NewVariableFromString("output_variable2", "string")
	v2.Scopes = []string{"deploy"}
	m := NewReleaseMetadata("test-release", "0.1")
	m.AddOutputVariable(v1)
	m.AddOutputVariable(v2)
	vars := m.GetOutputs("deploy")
	c.Assert(len(vars), Equals, 2)
	c.Assert(vars[0], Equals, v1)
	c.Assert(vars[1], Equals, v2)
	vars = m.GetOutputs("build")
	c.Assert(len(vars), Equals, 1)
	c.Assert(vars[0], Equals, v1)
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
        "api_version": 1,
        "project": "my-project",
        "consumes": [{ "Name": "provider1" }, 
					 { "name" : "provider2", "scopes": ["deploy"] }],
        "name": "test-release",
        "description": "Test release",
        "version": "0.1",
        "variable_context": {
            "base": "test-depends-v1",
            "test-depends": "test-depends-v1"
        },
		"stages": {
			"deploy": {
				"script": "deploy.sh"
			}
		}
    }`
	m, err := NewReleaseMetadataFromJsonString(json)
	c.Assert(err, IsNil)
	c.Assert(m.ApiVersion, Equals, 1)
	c.Assert(m.GetProject(), Equals, "my-project")
	c.Assert(m.Name, Equals, "test-release")
	c.Assert(m.Description, Equals, "Test release")
	c.Assert(m.Version, Equals, "0.1")
	c.Assert(m.GetConsumes("deploy"), HasLen, 2)
	c.Assert(m.GetConsumes("deploy")[0], Equals, "provider1")
	c.Assert(m.GetConsumes("deploy")[1], Equals, "provider2")
	c.Assert(m.GetConsumes("build"), HasLen, 1)
	c.Assert(m.GetConsumes("build")[0], Equals, "provider1")
	c.Assert(m.Stages["deploy"].RelativeScript, Equals, "deploy.sh")
}

func (s *metadataSuite) Test_AddInputVariable(c *C) {
	m := NewReleaseMetadata("test", "1.0")
	c.Assert(m.GetInputs("build"), HasLen, 0)
	variable, _ := variables.NewVariableFromString("testing", "string")
	m.AddInputVariable(variable)
	c.Assert(m.GetInputs("build"), HasLen, 1)
	c.Assert(m.GetInputs("deploy"), HasLen, 1)
	variable2, _ := variables.NewVariableFromString("test", "string")
	variable2.Scopes = []string{"deploy"}
	m.AddInputVariable(variable2)
	c.Assert(m.GetInputs("build"), HasLen, 1)
	c.Assert(m.GetInputs("deploy"), HasLen, 2)
	dep, _ := variables.NewVariableFromString("test", "string")
	dep.Scopes = []string{"build", "deploy"}
	m.AddInputVariable(dep)
	c.Assert(m.GetInputs("build"), HasLen, 2)
	c.Assert(m.GetInputs("deploy"), HasLen, 2)
	dep2, _ := variables.NewVariableFromString("test", "string")
	dep2.Scopes = []string{"deploy"}
	m.AddInputVariable(dep2)
	c.Assert(m.GetInputs("build"), HasLen, 2, Commentf("Most scopes win"))
	c.Assert(m.GetInputs("deploy"), HasLen, 2)
}

func (s *metadataSuite) Test_AddOutputVariable(c *C) {
	m := NewReleaseMetadata("test", "1.0")
	c.Assert(m.GetOutputs("build"), HasLen, 0)
	variable, _ := variables.NewVariableFromString("testing", "string")
	m.AddOutputVariable(variable)
	c.Assert(m.GetOutputs("build"), HasLen, 1)
	c.Assert(m.GetOutputs("deploy"), HasLen, 1)
	variable2, _ := variables.NewVariableFromString("test", "string")
	variable2.Scopes = []string{"deploy"}
	m.AddOutputVariable(variable2)
	c.Assert(m.GetOutputs("build"), HasLen, 1)
	c.Assert(m.GetOutputs("deploy"), HasLen, 2)
	dep, _ := variables.NewVariableFromString("test", "string")
	dep.Scopes = []string{"build", "deploy"}
	m.AddOutputVariable(dep)
	c.Assert(m.GetOutputs("build"), HasLen, 2)
	c.Assert(m.GetOutputs("deploy"), HasLen, 2)
	dep2, _ := variables.NewVariableFromString("test", "string")
	dep2.Scopes = []string{"deploy"}
	m.AddOutputVariable(dep2)
	c.Assert(m.GetOutputs("build"), HasLen, 2, Commentf("Most scopes win"))
	c.Assert(m.GetOutputs("deploy"), HasLen, 2)
}

func (s *metadataSuite) Test_AddConsumes(c *C) {
	m := NewReleaseMetadata("test", "1.0")
	c.Assert(m.GetConsumes("build"), HasLen, 0)
	consumer := NewConsumerConfig("all-scopes")
	m.AddConsumes(consumer)
	c.Assert(m.GetConsumes("build"), HasLen, 1)
	c.Assert(m.GetConsumes("deploy"), HasLen, 1)
	consumer2 := NewConsumerConfig("deploy-scope")
	consumer2.Scopes = []string{"deploy"}
	m.AddConsumes(consumer2)
	c.Assert(m.GetConsumes("build"), HasLen, 1)
	c.Assert(m.GetConsumes("deploy"), HasLen, 2)
	dep := NewConsumerConfig("deploy-scope")
	dep.Scopes = []string{"build", "deploy"}
	m.AddConsumes(dep)
	c.Assert(m.GetConsumes("build"), HasLen, 2)
	c.Assert(m.GetConsumes("deploy"), HasLen, 2)
	dep2 := NewConsumerConfig("deploy-scope")
	dep2.Scopes = []string{"deploy"}
	m.AddConsumes(dep2)
	c.Assert(m.GetConsumes("build"), HasLen, 2, Commentf("Most scopes win"))
	c.Assert(m.GetConsumes("deploy"), HasLen, 2)
	dep3 := NewConsumerConfig("deploy-scope")
	dep3.VariableName = "t"
	m.AddConsumes(dep3)
	c.Assert(m.GetConsumes("build"), HasLen, 3, Commentf("Variable name is part of key"))
	c.Assert(m.GetConsumes("deploy"), HasLen, 3)
	m.AddConsumes(dep3)
	c.Assert(m.GetConsumes("build"), HasLen, 3)
	c.Assert(m.GetConsumes("deploy"), HasLen, 3)
}
