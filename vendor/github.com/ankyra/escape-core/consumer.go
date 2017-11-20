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

import "fmt"

/*

Unlike Dependencies, which are resolved at build time and provide tight
coupling, we can use Consumers and Providers to resolve and loosely couple
dependencies at deployment time.

## Escape Plan

Consumers are configured in the [`consumes`](/docs/escape-plan/#consumes)
field of the Escape Plan.

Providers are configured in the [`provides`](/docs/escape-plan/#provides)
field of the Escape Plan.

*/
type ConsumerConfig struct {
	Name   string   `json:"name" yaml:"name"`
	Scopes []string `json:"scopes" yaml:"scopes"`
}

func NewConsumerConfig(name string) *ConsumerConfig {
	return &ConsumerConfig{
		Name:   name,
		Scopes: []string{"build", "deploy"},
	}
}

func NewConsumerConfigFromInterface(v interface{}) (*ConsumerConfig, error) {
	switch v.(type) {
	case string:
		return NewConsumerConfig(v.(string)), nil
	case map[interface{}]interface{}:
		return NewConsumerConfigFromMap(v.(map[interface{}]interface{}))
	}
	return nil, fmt.Errorf("Expecting dict or string type")
}

func NewConsumerConfigFromMap(dep map[interface{}]interface{}) (*ConsumerConfig, error) {
	var name string
	scopes := []string{}
	for key, val := range dep {
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string key in consumer")
		}
		if keyStr == "name" {
			valString, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("Expecting string for consumer 'name' got '%T'", val)
			}
			name = valString
		} else if key == "scopes" {
			s, err := parseScopesFromInterface(val)
			if err != nil {
				return nil, err
			}
			scopes = s
		}
	}
	if name == "" {
		return nil, fmt.Errorf("Missing 'name' in consumer")
	}
	cfg := NewConsumerConfig(name)
	cfg.Scopes = scopes
	return cfg, cfg.ValidateAndFix()
}

func parseScopesFromInterface(val interface{}) ([]string, error) {
	valList, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Expecting string in scopes, got '%v' (%T)", val, val)
	}
	scopes := []string{}
	for _, val := range valList {
		kStr, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string in scopes, got '%v' (%T)", val, val)
		}
		scopes = append(scopes, kStr)
	}
	return scopes, nil
}

func (c *ConsumerConfig) ValidateAndFix() error {
	if c.Scopes == nil || len(c.Scopes) == 0 {
		c.Scopes = []string{"build", "deploy"}
	}
	if c.Name == "" {
		return fmt.Errorf("Missing name for Consumer")
	}
	return nil
}

func (c *ConsumerConfig) InScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}
