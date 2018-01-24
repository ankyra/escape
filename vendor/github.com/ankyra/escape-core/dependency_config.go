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
	"fmt"
)

/*

## Escape Plan

Dependencies are configured in the [`depends`](/docs/reference/escape-plan/#depends)
field of the Escape plan.

*/
type DependencyConfig struct {
	// The release id is required and is resolved at *build* time and then
	// persisted in the release metadata ensuring that deployments always use
	// the same versions.
	//
	// Examples:
	// - To always use the latest version: `my-organisation/my-dependency-latest`
	// - To always use version 0.1.1: `my-organisation/my-dependency-v0.1.1`
	// - To always use the latest version in the 0.1 series: `my-organisation/my-dependency-v0.1.@`
	// - To make it possible to reference a dependency using a different name: `my-organisation/my-dependency-latest as my-name`
	ReleaseId string `json:"release_id" yaml:"release_id"`

	// Define the values of dependency inputs using Escape Script.
	Mapping map[string]interface{} `json:"mapping" yaml:"mapping"`

	// Define the values of dependency inputs using Escape Script when running
	// stages in the build scope.
	BuildMapping map[string]interface{} `json:"build_mapping" yaml:"build_mapping"`

	// Define the values of dependency inputs using Escape Script when running
	// stages in the deploy scope.
	DeployMapping map[string]interface{} `json:"deploy_mapping" yaml:"deploy_mapping"`

	// Map providers from the parent to dependencies.
	//
	// Example:
	// ```
	// consumes:
	// - my-provider
	// depends:
	// - release_id: my-org/my-dep-latest
	//     consumes:
	//       provider: $my-provider.deployment
	// ```
	Consumes map[string]string `json:"consumes" yaml:"consumes"`

	// The name of the (sub)-deployment. This defaults to the versionless release id;
	// e.g. if the release_id is `my-org/my-dep-v1.0` then the DeploymentName will be
	// `my-org/my-dep` by default.
	DeploymentName string `json:"deployment_name" yaml:"deployment_name"`

	// The variable used to reference this dependency. By default the variable
	// name is the versionless release id of the dependency, but this can be
	// overruled by renaming the dependency (e.g. `my-org/my-release-latest as
	// my-variable`. This field will be set automatically at build time.
	// Overwriting this field in the Escape plan has no effect.
	VariableName string `json:"variable" yaml:"variable"`

	// A list of scopes (`build`, `deploy`) that defines during which stage(s)
	// this dependency should be fetched and deployed. *Currently not implemented!*
	Scopes []string `json:"scopes" yaml:"scopes"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"my-org"`.
	Project string `json:"-" yaml:"-"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"my-name"`.
	Name string `json:"-" yaml:"-"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"1.0"`.
	Version string `json:"-" yaml:"-"`
}

func NewDependencyConfig(releaseId string) *DependencyConfig {
	return &DependencyConfig{
		ReleaseId:     releaseId,
		Mapping:       map[string]interface{}{},
		BuildMapping:  map[string]interface{}{},
		DeployMapping: map[string]interface{}{},
		Scopes:        []string{"build", "deploy"},
		Consumes:      map[string]string{},
	}
}

func DependencyNeedsResolvingError(dependencyReleaseId string) error {
	return fmt.Errorf("The dependency '%s' needs its version resolved.", dependencyReleaseId)
}

func (d *DependencyConfig) Validate(m *ReleaseMetadata) error {
	if d.BuildMapping == nil {
		d.BuildMapping = map[string]interface{}{}
	}
	if d.DeployMapping == nil {
		d.DeployMapping = map[string]interface{}{}
	}
	if d.Scopes == nil || len(d.Scopes) == 0 {
		d.Scopes = []string{"build", "deploy"}
	}
	dep, err := NewDependencyFromString(d.ReleaseId)
	if err != nil {
		return err
	}
	if dep.NeedsResolving() {
		return DependencyNeedsResolvingError(d.ReleaseId)
	}
	return nil
}

func (d *DependencyConfig) GetMapping(scope string) map[string]interface{} {
	if scope == "build" {
		return d.BuildMapping
	}
	if scope == "deploy" {
		return d.DeployMapping
	}
	return nil
}

func (d *DependencyConfig) AddVariableMapping(scopes []string, id, key string) {
	for _, scope := range scopes {
		mapping := d.GetMapping(scope)
		if mapping != nil {
			_, found := mapping[id]
			if !found {
				mapping[id] = key
			}
		}
	}
}

func (d *DependencyConfig) InScope(scope string) bool {
	for _, s := range d.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func NewDependencyConfigFromMap(dep map[interface{}]interface{}) (*DependencyConfig, error) {
	var releaseId string
	buildMapping := map[string]interface{}{}
	deployMapping := map[string]interface{}{}
	consumes := map[string]string{}
	scopes := []string{}
	for key, val := range dep {
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string key in dependency")
		}
		if keyStr == "release_id" {
			valString, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("Expecting string for dependency 'release_id' got '%T'", val)
			}
			releaseId = valString
		} else if key == "mapping" { // backwards compatibility with release metadata <= 6
			valMap, ok := val.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting dict for dependency 'mapping' got '%T'", val)
			}
			for k, v := range valMap {
				kStr, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string key in dependency 'mapping'")
				}
				buildMapping[kStr] = v
				deployMapping[kStr] = v
			}
		} else if key == "build_mapping" {
			valMap, ok := val.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting dict for dependency 'build_mapping' got '%T'", val)
			}
			for k, v := range valMap {
				kStr, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string key in dependency 'build_mapping'")
				}
				buildMapping[kStr] = v
			}
		} else if key == "deploy_mapping" {
			valMap, ok := val.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting dict for dependency 'deploy_mapping' got '%T'", val)
			}
			for k, v := range valMap {
				kStr, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string key in dependency 'deploy_mapping'")
				}
				deployMapping[kStr] = v
			}
		} else if key == "consumes" {
			valMap, ok := val.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting dict for dependency 'consumes' got '%T'", val)
			}
			for k, v := range valMap {
				kStr, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string key in dependency consumer mapping")
				}
				vStr, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string value in dependency consumer mapping")
				}
				consumes[kStr] = vStr
			}
		} else if key == "scopes" {
			s, err := parseScopesFromInterface(val)
			if err != nil {
				return nil, err
			}
			scopes = s
		}
	}
	if releaseId == "" {
		return nil, fmt.Errorf("Missing 'release_id' in dependency")
	}
	cfg := NewDependencyConfig(releaseId)
	cfg.BuildMapping = buildMapping
	cfg.DeployMapping = deployMapping
	cfg.Scopes = scopes
	cfg.Consumes = consumes
	return cfg, nil
}
