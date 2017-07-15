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
	ReleaseId string                 `json:"release_id" yaml:"release_id"`
	Mapping   map[string]interface{} `json:"mapping" yaml:"mapping"`
}

func (d *DependencyConfig) Validate(m *ReleaseMetadata) error {
	if d.Mapping == nil {
		d.Mapping = map[string]interface{}{}
	}
	for _, input := range m.Inputs {
		_, alreadySet := d.Mapping[input.Id]
		if alreadySet {
			continue
		}
		if input.EvalBeforeDependencies {
			d.Mapping[input.Id] = "$this.inputs." + input.Id
		}
	}
	return nil
}

func NewDependencyConfig(releaseId string) *DependencyConfig {
	return &DependencyConfig{
		ReleaseId: releaseId,
		Mapping:   map[string]interface{}{},
	}
}

func NewDependencyConfigFromMap(dep map[interface{}]interface{}) (*DependencyConfig, error) {
	var releaseId string
	mapping := map[string]interface{}{}
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
		} else if key == "mapping" {
			valMap, ok := val.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting dict for dependency 'mapping' got '%T'", val)
			}
			for k, v := range valMap {
				kStr, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string key in dependency mapping")
				}
				mapping[kStr] = v
			}
		}
	}
	if releaseId == "" {
		return nil, fmt.Errorf("Missing 'release_id' in dependency")
	}
	cfg := NewDependencyConfig(releaseId)
	cfg.Mapping = mapping
	return cfg, nil
}
