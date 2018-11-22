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
	"strings"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/scopes"
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
	Scopes scopes.Scopes `json:"scopes" yaml:"scopes"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"my-org"`.
	Project string `json:"-" yaml:"-"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"my-name"`.
	Name string `json:"-" yaml:"-"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name-v1.0"` this value is `"1.0"`.
	Version string `json:"-" yaml:"-"`

	// Parsed out of the release ID. For example: when release id is
	// `"my-org/my-name:tag"` this value is `"tag"`.
	Tag string `json:"-" yaml:"-"`
}

type ResolvedDependencyConfig struct {
	*DependencyConfig
	ReleaseMetadata *ReleaseMetadata
}

func NewDependencyConfig(releaseId string) *DependencyConfig {
	return &DependencyConfig{
		ReleaseId:     releaseId,
		Mapping:       map[string]interface{}{},
		BuildMapping:  map[string]interface{}{},
		DeployMapping: map[string]interface{}{},
		Scopes:        scopes.AllScopes,
		Consumes:      map[string]string{},
	}
}

func (d *DependencyConfig) Resolve(m *ReleaseMetadata) *ResolvedDependencyConfig {
	return &ResolvedDependencyConfig{
		DependencyConfig: d,
		ReleaseMetadata:  m,
	}
}

func DependencyNeedsResolvingError(dependencyReleaseId string) error {
	return fmt.Errorf("The dependency '%s' needs its version resolved.", dependencyReleaseId)
}

func (d *DependencyConfig) Copy() *DependencyConfig {
	result := NewDependencyConfig(d.ReleaseId)
	for k, v := range d.Mapping {
		result.Mapping[k] = v
	}
	for k, v := range d.BuildMapping {
		result.BuildMapping[k] = v
	}
	for k, v := range d.DeployMapping {
		result.DeployMapping[k] = v
	}
	for k, v := range d.Consumes {
		result.Consumes[k] = v
	}
	result.DeploymentName = d.DeploymentName
	result.VariableName = d.VariableName
	result.Scopes = d.Scopes.Copy()
	result.Project = d.Project
	result.Name = d.Name
	result.Version = d.Version
	result.Tag = d.Tag
	return result
}

func (d *DependencyConfig) EnsureConfigIsParsed() error {
	parsed, err := parsers.ParseDependency(d.ReleaseId)
	if err != nil {
		return err
	}
	d.ReleaseId = parsed.QualifiedReleaseId.ToString()
	d.Project = parsed.Project
	d.Name = parsed.Name
	d.Version = parsed.Version
	d.Tag = parsed.Tag
	if d.VariableName == "" {
		d.VariableName = parsed.VariableName
	}
	return nil
}

func (d *DependencyConfig) NeedsResolving() bool {
	return d.Tag != "" || d.Version == "latest" || strings.HasSuffix(d.Version, ".@")
}

func (d *DependencyConfig) GetVersionAsString() (version string) {
	if d.Tag != "" {
		return d.Tag
	}
	version = "v" + d.Version
	if d.Version == "latest" {
		version = d.Version
	}
	return version
}

func (d *DependencyConfig) Validate(m *ReleaseMetadata) error {
	if d.BuildMapping == nil {
		d.BuildMapping = map[string]interface{}{}
	}
	if d.DeployMapping == nil {
		d.DeployMapping = map[string]interface{}{}
	}
	if d.Scopes == nil || len(d.Scopes) == 0 {
		d.Scopes = scopes.AllScopes
	}
	if err := d.EnsureConfigIsParsed(); err != nil {
		return err
	}
	if d.NeedsResolving() {
		return DependencyNeedsResolvingError(d.ReleaseId)
	}
	d.DeploymentName = d.VariableName
	if d.DeploymentName == "" {
		d.DeploymentName = d.Project + "/" + d.Name
	}
	if d.VariableName == "" {
		d.VariableName = d.Project + "/" + d.Name
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
	return d.Scopes.InScope(scope)
}

func ExpectingTypeForDependencyFieldError(typ, field string, val interface{}) error {
	return fmt.Errorf("Expecting %s for dependency '%s'; got '%T'", typ, field, val)
}

func ExpectingStringKeyInMapError(field string, val interface{}) error {
	return fmt.Errorf("Expecting string key in dependency '%s'; got '%T'", field, val)
}

func stringFromInterface(field string, val interface{}) (string, error) {
	valString, ok := val.(string)
	if !ok {
		return "", ExpectingTypeForDependencyFieldError("string", field, val)
	}
	return valString, nil
}

func mapFromInterface(field string, val interface{}) (map[string]interface{}, error) {
	valMap, ok := val.(map[interface{}]interface{})
	if !ok {
		return nil, ExpectingTypeForDependencyFieldError("dict", field, val)
	}
	result := map[string]interface{}{}
	for k, v := range valMap {
		kStr, ok := k.(string)
		if !ok {
			return nil, ExpectingStringKeyInMapError(field, k)
		}
		result[kStr] = v
	}
	return result, nil
}

func NewDependencyConfigFromMap(dep map[interface{}]interface{}) (*DependencyConfig, error) {
	var releaseId, deploymentName, variable string
	buildMapping := map[string]interface{}{}
	deployMapping := map[string]interface{}{}
	consumes := map[string]string{}
	depScopes := []string{}
	for key, val := range dep {
		var err error
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string key in dependency")
		}
		if keyStr == "release_id" {
			releaseId, err = stringFromInterface("release_id", val)
			if err != nil {
				return nil, err
			}
		} else if keyStr == "deployment_name" {
			deploymentName, err = stringFromInterface("deployment_name", val)
			if err != nil {
				return nil, err
			}
		} else if keyStr == "variable" {
			variable, err = stringFromInterface("variable", val)
			if err != nil {
				return nil, err
			}
		} else if key == "mapping" { // backwards compatibility with release metadata <= 6
			valMap, err := mapFromInterface("mapping", val)
			if err != nil {
				return nil, err
			}
			for k, v := range valMap {
				buildMapping[k] = v
				deployMapping[k] = v
			}
		} else if key == "build_mapping" {
			valMap, err := mapFromInterface("build_mapping", val)
			if err != nil {
				return nil, err
			}
			for k, v := range valMap {
				buildMapping[k] = v
			}
		} else if key == "deploy_mapping" {
			valMap, err := mapFromInterface("deploy_mapping", val)
			if err != nil {
				return nil, err
			}
			for k, v := range valMap {
				deployMapping[k] = v
			}
		} else if key == "consumes" {
			valMap, err := mapFromInterface("consumes", val)
			if err != nil {
				return nil, err
			}
			for k, v := range valMap {
				vStr, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string value in dependency consumer mapping")
				}
				consumes[k] = vStr
			}
		} else if key == "scopes" {
			s, err := scopes.NewScopesFromInterface(val)
			if err != nil {
				return nil, err
			}
			depScopes = s
		}
	}
	if releaseId == "" {
		return nil, fmt.Errorf("Missing 'release_id' in dependency")
	}
	cfg := NewDependencyConfig(releaseId)
	cfg.DeploymentName = deploymentName
	cfg.VariableName = variable
	cfg.BuildMapping = buildMapping
	cfg.DeployMapping = deployMapping
	cfg.Scopes = depScopes
	cfg.Consumes = consumes
	return cfg, nil
}
