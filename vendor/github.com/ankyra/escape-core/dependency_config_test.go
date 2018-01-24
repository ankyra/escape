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

func (s *metadataSuite) Test_NewDependencyConfig_happy_path(c *C) {
	metadata := NewReleaseMetadata("name", "1.0")
	dep := NewDependencyConfig("my-dependency-v1.1")
	dep.BuildMapping = nil
	dep.DeployMapping = nil
	c.Assert(dep.Validate(metadata), IsNil)
	c.Assert(dep.BuildMapping, Not(IsNil))
	c.Assert(dep.BuildMapping, HasLen, 0)
	c.Assert(dep.DeployMapping, Not(IsNil))
	c.Assert(dep.DeployMapping, HasLen, 0)
	c.Assert(dep.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(dep.Consumes, DeepEquals, map[string]string{})
}

func (s *metadataSuite) Test_NewDependencyConfig_fails_if_invalid_dependency_string(c *C) {
	cases := []string{
		"",
		"my",
		"my-dependency",
		"my-dependency-v£%%",
	}
	for _, test := range cases {
		metadata := NewReleaseMetadata("name", "1.0")
		dep := NewDependencyConfig(test)
		dep.BuildMapping = nil
		dep.DeployMapping = nil
		c.Assert(dep.Validate(metadata), NotNil)
	}
}

func (s *metadataSuite) Test_NewDependencyConfig_fails_if_version_needs_resolving(c *C) {
	cases := []string{
		"my-dependency-latest",
		"my-dependency-v1.0.@",
		"my-dependency-v0.@",
		"my-dependency-v@",
		"my-dependency-@",
	}
	for _, test := range cases {
		metadata := NewReleaseMetadata("name", "1.0")
		dep := NewDependencyConfig(test)
		dep.BuildMapping = nil
		dep.DeployMapping = nil
		c.Assert(dep.Validate(metadata), DeepEquals, DependencyNeedsResolvingError(test))
	}
}

func (s *metadataSuite) Test_NewDependencyConfigFromMap(c *C) {
	dep, err := NewDependencyConfigFromMap(map[interface{}]interface{}{
		"release_id": "test-latest",
		"build_mapping": map[interface{}]interface{}{
			"build": "building",
		},
		"deploy_mapping": map[interface{}]interface{}{
			"deploy": "deploying",
		},
		"mapping": map[interface{}]interface{}{
			"input_variable1": "test",
		},
		"scopes": []interface{}{"build"},
		"consumes": map[interface{}]interface{}{
			"test": "whatver",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(dep.ReleaseId, Equals, "test-latest")
	c.Assert(dep.BuildMapping, Not(IsNil))
	c.Assert(dep.BuildMapping, HasLen, 2)
	c.Assert(dep.BuildMapping["input_variable1"], Equals, "test")
	c.Assert(dep.BuildMapping["build"], Equals, "building")
	c.Assert(dep.DeployMapping, Not(IsNil))
	c.Assert(dep.DeployMapping, HasLen, 2)
	c.Assert(dep.DeployMapping["input_variable1"], Equals, "test")
	c.Assert(dep.DeployMapping["deploy"], Equals, "deploying")
	c.Assert(dep.Scopes, DeepEquals, []string{"build"})
	c.Assert(dep.Consumes, DeepEquals, map[string]string{"test": "whatver"})
}
