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
	"errors"
	"fmt"
	"github.com/ankyra/escape-core/variables"
)

type Errand struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Script      string                `json:"script"`
	Inputs      []*variables.Variable `json:"inputs"`
}

func (e *Errand) GetInputs() []*variables.Variable {
	result := []*variables.Variable{}
	for _, i := range e.Inputs {
		result = append(result, i)
	}
	return result
}

func (e *Errand) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("Missing name in errand")
	} else if e.Script == "" {
		return fmt.Errorf("Missing script in errand '%s'", e.Name)
	}
	if e.Inputs == nil {
		return nil
	}
	for _, v := range e.Inputs {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("Error in errand '%s' variable: %s", e.Name, err.Error())
		}
	}
	return nil
}

func NewErrand(name, script, description string) *Errand {
	result := &Errand{
		Name:        name,
		Script:      script,
		Description: description,
		Inputs:      []*variables.Variable{},
	}
	return result
}

func NewErrandFromDict(name string, dict interface{}) (*Errand, error) {
	switch dict.(type) {
	case map[interface{}]interface{}:
		errandMap := dict.(map[interface{}]interface{})
		description := ""
		script := ""
		inputs := []*variables.Variable{}
		for key, val := range errandMap {
			switch key.(type) {
			case string:
				if key == "description" {
					str, err := getString(val)
					if err != nil {
						return nil, errors.New("Expecting string value for description field in errand " + name)
					}
					description = str

				} else if key == "script" {
					str, err := getString(val)
					if err != nil {
						return nil, errors.New("Expecting string value for script field in errand " + name)
					}
					script = str
				} else if key == "inputs" {
					switch val.(type) {
					case []interface{}:
						inputDicts := val.([]interface{})
						for _, inputDict := range inputDicts {
							switch inputDict.(type) {
							case map[interface{}]interface{}:
								dict := inputDict.(map[interface{}]interface{})
								v, err := variables.NewVariableFromDict(dict)
								if err != nil {
									return nil, err
								}
								inputs = append(inputs, v)
							case string:
								stringVar := inputDict.(string)
								v, err := variables.NewVariableFromString(stringVar, "string")
								if err != nil {
									return nil, errors.New("Not a valid errand variable: " + stringVar)
								}
								inputs = append(inputs, v)
							default:
								return nil, errors.New("Expecting dict or string type for input item in errand " + name)
							}
						}
					default:
						return nil, errors.New("Expecting list type for inputs key in errand " + name)

					}
				}
			default:
				return nil, errors.New("Expecting string key for errand " + name)
			}

		}
		result := &Errand{
			Name:        name,
			Description: description,
			Script:      script,
			Inputs:      inputs,
		}
		return result, result.Validate()
	}
	return nil, errors.New("Expecting a dictionary for errand " + name)
}

func getString(val interface{}) (string, error) {
	switch val.(type) {
	case string:
		return val.(string), nil
	}
	return "", errors.New("Expecting string")
}
