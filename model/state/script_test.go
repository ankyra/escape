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

package state

import (
	. "github.com/ankyra/escape-client/model/state/types"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
	"testing"
)

type deplSuite struct{}

var _ = Suite(&deplSuite{})

func Test(t *testing.T) { TestingT(t) }

var depl *DeploymentState
var deplWithDeps *DeploymentState
var fullDepl *DeploymentState

func (s *deplSuite) SetUpTest(c *C) {
	var err error
	env, err := NewLocalStateProvider("testdata/project.json").Load("prj", "dev")
	c.Assert(err, IsNil)
	depl, err = env.GetDeploymentState([]string{"archive-release"})
	c.Assert(err, IsNil)

	deplWithDeps, err = env.GetDeploymentState([]string{"archive-release-with-deps", "archive-release"})
	c.Assert(err, IsNil)

	fullDepl, err = env.GetDeploymentState([]string{"archive-full"})
	c.Assert(err, IsNil)
}

func (s *deplSuite) Test_ToScript(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Metadata["value"] = "yo"
	input, err := variables.NewVariableFromString("user_level", "string")
	c.Assert(err, IsNil)
	metadata.AddInputVariable(input)
	metadata.AddOutputVariable(input)
	unit := toScript(depl, metadata, "deploy")
	dicts := map[string][]string{
		"inputs":   []string{"user_level"},
		"outputs":  []string{"user_level"},
		"metadata": []string{"value"},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *deplSuite) Test_ToScript_doesnt_include_variable_that_are_not_defined_in_release_metadata(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	unit := toScript(depl, metadata, "deploy")
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *deplSuite) Test_ToScriptEnvironment(c *C) {
	metadataMap := map[string]*core.ReleaseMetadata{
		"this": core.NewReleaseMetadata("test", "1.0"),
	}
	env, err := ToScriptEnvironment(depl, metadataMap, "build")
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-release")
}

func (s *deplSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	metadataMap := map[string]*core.ReleaseMetadata{
		"this": metadata,
	}
	_, err := ToScriptEnvironment(depl, metadataMap, "build")
	c.Assert(err, Not(IsNil))
}

func (s *deplSuite) Test_ToScriptEnvironment_fails_if_missing_provider_metadata(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	metadataMap := map[string]*core.ReleaseMetadata{
		"this": metadata,
	}
	depl.Providers["test"] = "archive-full"
	_, err := ToScriptEnvironment(depl, metadataMap, "build")
	c.Assert(err, Not(IsNil))
}

func (s *deplSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state_in_environment(c *C) {
	archiveMetadata := core.NewReleaseMetadata("test", "1.0")
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	metadataMap := map[string]*core.ReleaseMetadata{
		"this":            metadata,
		"archive-full-v1": archiveMetadata,
	}
	depl.Providers["test"] = "this-doesnt-exist"
	_, err := ToScriptEnvironment(depl, metadataMap, "build")
	c.Assert(err, Not(IsNil))
}

func (s *deplSuite) Test_ToScriptEnvironment_adds_consumers(c *C) {
	archiveMetadata := core.NewReleaseMetadata("test", "1.0")
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	metadataMap := map[string]*core.ReleaseMetadata{
		"this":         metadata,
		"archive-full": archiveMetadata,
	}
	depl.Providers["test"] = "archive-full"
	env, err := ToScriptEnvironment(depl, metadataMap, "build")
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-release")
	test_helper_check_script_environment(c, dict["test"], dicts, "archive-full")
}

func (s *deplSuite) Test_ToScriptEnvironment_fails_if_dependency_metadata_is_missing(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Depends = []string{"archive-dep-v1.0"}
	metadataMap := map[string]*core.ReleaseMetadata{
		"this": metadata,
	}
	_, err := ToScriptEnvironment(fullDepl, metadataMap, "build")
	c.Assert(err, Not(IsNil))
}

func (s *deplSuite) Test_ToScriptEnvironment_doesnt_add_dependencies_that_are_not_in_metadata(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadataMap := map[string]*core.ReleaseMetadata{
		"this": metadata,
	}
	env, err := ToScriptEnvironment(fullDepl, metadataMap, "build")
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	c.Assert(dict["archive-dep-v1.0"], IsNil)
}

func (s *deplSuite) Test_ToScriptEnvironment_adds_dependencies(c *C) {
	depMetadata := core.NewReleaseMetadata("test", "1.0")
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Depends = []string{"archive-dep-v1.0"}
	metadataMap := map[string]*core.ReleaseMetadata{
		"this":        metadata,
		"archive-dep": depMetadata,
	}
	env, err := ToScriptEnvironment(fullDepl, metadataMap, "build")
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	test_helper_check_script_environment(c, dict["archive-dep-v1.0"], dicts, "archive-dep")
}

func test_helper_check_script_environment(c *C, unit script.Script, dicts map[string][]string, name string) {
	c.Assert(script.IsDictAtom(unit), Equals, true)
	dict := script.ExpectDictAtom(unit)
	strings := map[string]string{
		"version":     "1.0",
		"description": "",
		"logo":        "",
		"id":          "test-v1.0",
		"name":        "test",
		"branch":      "",
		"revision":    "",
		"project":     "project_name",
		"environment": "dev",
		"deployment":  name,
	}
	for key, val := range strings {
		c.Assert(script.IsStringAtom(dict[key]), Equals, true, Commentf("Expecting %s to be of type string, but was %T", key, dict[key]))
		c.Assert(script.ExpectStringAtom(dict[key]), Equals, val)
	}
	for key, keys := range dicts {
		c.Assert(script.IsDictAtom(dict[key]), Equals, true, Commentf("Expecting %s to be of type dict, but was %T", key, dict[key]))
		d := script.ExpectDictAtom(dict[key])
		c.Assert(d, HasLen, len(keys), Commentf("Expecting %d values in %s dict.", len(keys), key))
		for _, k := range keys {
			c.Assert(script.IsStringAtom(d[k]), Equals, true, Commentf("Expecting %s to be of type string, but was %T", k, d[k]))
		}
	}
}
