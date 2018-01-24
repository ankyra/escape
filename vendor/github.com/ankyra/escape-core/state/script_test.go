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

package state

import (
	"errors"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

type scriptSuite struct{}

var _ = Suite(&scriptSuite{})

func (s *scriptSuite) SetUpTest(c *C) {
	var err error
	prj, err := NewProjectStateFromFile("prj", "testdata/project_script.json", nil)
	c.Assert(err, IsNil)
	env := prj.GetEnvironmentStateOrMakeNew("dev")
	depl = env.GetOrCreateDeploymentState("archive-release")
	fullDepl = env.GetOrCreateDeploymentState("archive-full")
	dep := env.GetOrCreateDeploymentState("archive-release-with-deps")
	deplWithDeps = dep.GetDeploymentOrMakeNew("deploy", "archive-release")
}

func (s *scriptSuite) Test_ToScript(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Metadata["value"] = "yo"
	input, err := variables.NewVariableFromString("user_level", "string")
	c.Assert(err, IsNil)
	metadata.AddInputVariable(input)
	metadata.AddOutputVariable(input)
	unit := newStateCompiler(nil).compileState(depl, metadata, "deploy", true)
	dicts := map[string][]string{
		"inputs":   []string{"user_level"},
		"outputs":  []string{"user_level"},
		"metadata": []string{"value"},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *scriptSuite) Test_ToScript_doesnt_include_variable_that_are_not_defined_in_release_metadata(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	unit := newStateCompiler(nil).compileState(depl, metadata, "deploy", true)
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, unit, dicts, "archive-release")
}

func (s *scriptSuite) Test_ToScriptEnvironment_adds_dependencies(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"_/archive-dep-v1.0":  core.NewReleaseMetadata("test", "1.0"),
		"_/archive-dep2-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"archive-dep-v1.0 as archive-dep", "archive-dep2-v1.0"})

	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	test_helper_check_script_environment(c, dict["archive-dep"], dicts, "archive-full:archive-dep")
	test_helper_check_script_environment(c, dict["_/archive-dep2"], dicts, "archive-full:_/archive-dep2")
}

func (s *scriptSuite) Test_ToScriptEnvironment_doesnt_add_dependencies_that_are_not_in_metadata(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	env, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-full")
	c.Assert(dict["test-v1.0"], IsNil)
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_deployment_state_is_missing(c *C) {
	metadata := core.NewReleaseMetadata("test", "1.0")
	_, err := ToScriptEnvironment(nil, metadata, "build", nil)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_dependency_metadata_is_missing(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetDependencies([]string{"archive-dep-v1.0"})
	_, err := ToScriptEnvironment(fullDepl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_adds_consumers(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "archive-full")
	env, err := ToScriptEnvironment(depl, metadata, "build", resolver)
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

func (s *scriptSuite) Test_ToScriptEnvironment_adds_renamed_consumers(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	cfg, _ := core.NewConsumerConfigFromString("test as t")
	cfg2, _ := core.NewConsumerConfigFromString("test as t2")
	metadata.Consumes = []*core.ConsumerConfig{cfg, cfg2}
	depl.SetProvider("build", "t", "archive-full")
	depl.SetProvider("build", "t2", "archive-full")
	env, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, IsNil)
	c.Assert(script.IsDictAtom((*env)["$"]), Equals, true)
	dict := script.ExpectDictAtom((*env)["$"])
	dicts := map[string][]string{
		"inputs":   []string{},
		"outputs":  []string{},
		"metadata": []string{},
	}
	test_helper_check_script_environment(c, dict["this"], dicts, "archive-release")
	test_helper_check_script_environment(c, dict["t"], dicts, "archive-full")
	test_helper_check_script_environment(c, dict["t2"], dicts, "archive-full")
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_renamed_consumer_not_configured(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	cfg, _ := core.NewConsumerConfigFromString("test as t")
	metadata.Consumes = []*core.ConsumerConfig{cfg}
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, DeepEquals, errors.New(`Provider 't' of type 'test' has not been configured in the deployment state.`))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_metadata(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "archive-full")
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func (s *scriptSuite) Test_ToScriptEnvironment_fails_if_missing_provider_state_in_environment(c *C) {
	resolver := newResolverFromMap(map[string]*core.ReleaseMetadata{
		"archive-full-v1.0": core.NewReleaseMetadata("test", "1.0"),
	})
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.SetConsumes([]string{"test"})
	depl.SetProvider("build", "test", "this-doesnt-exist")
	_, err := ToScriptEnvironment(depl, metadata, "build", resolver)
	c.Assert(err, Not(IsNil))
}

func test_helper_check_script_environment(c *C, unit script.Script, dicts map[string][]string, name string) {
	c.Assert(script.IsDictAtom(unit), Equals, true)
	dict := script.ExpectDictAtom(unit)
	strings := map[string]string{
		"version":     "1.0",
		"description": "",
		"logo":        "",
		"release":     "test-v1.0",
		"id":          "_/test-v1.0",
		"name":        "test",
		"branch":      "",
		"revision":    "",
		"project":     "project_name",
		"environment": "dev",
		"deployment":  name,
	}
	for key, val := range strings {
		c.Assert(script.IsStringAtom(dict[key]), Equals, true, Commentf("Expecting %s to be of type string, but was %T", key, dict[key]))
		c.Assert(script.ExpectStringAtom(dict[key]), Equals, val, Commentf("Expecting '%s' to be '%s'. Got '%s'", key, val, script.ExpectStringAtom(dict[key])))
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
