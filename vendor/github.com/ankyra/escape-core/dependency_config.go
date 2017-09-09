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
	"fmt"
)

type DependencyConfig struct {
	ReleaseId     string                 `json:"release_id" yaml:"release_id"`
	BuildMapping  map[string]interface{} `json:"build_mapping" yaml:"build_mapping"`
	DeployMapping map[string]interface{} `json:"deploy_mapping" yaml:"deploy_mapping"`
	Consumes      map[string]string      `json:"consumes" yaml:"consumes"`
	Scopes        []string               `json:"scopes" yaml:"scopes"`
}

func NewDependencyConfig(releaseId string) *DependencyConfig {
	return &DependencyConfig{
		ReleaseId:     releaseId,
		BuildMapping:  map[string]interface{}{},
		DeployMapping: map[string]interface{}{},
		Scopes:        []string{"build", "deploy"},
		Consumes:      map[string]string{},
	}
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
