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

package release

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/parsers"
	"github.com/ankyra/escape-client/model/release/variable_types"
	"github.com/ankyra/escape-client/model/script"
	"gopkg.in/yaml.v2"
)

type variable struct {
	Id          string                 `json:"id"`
	Type        string                 `json:"type"`
	Default     interface{}            `json:"default,omitempty"`
	Description string                 `json:"description,omitempty"`
	Friendly    string                 `json:"friendly,omitempty"`
	Visible     bool                   `json:"visible"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Sensitive   bool                   `json:"sensitive,omitempty"`
	Items       []interface{}          `json:"items"` // Only set for one_of variables
}

func (v *variable) GetId() string {
	return v.Id
}

type UntypedVariable map[interface{}]interface{}

func NewVariable() Variable {
	return &variable{
		Visible: true,
	}
}

func NewVariableFromString(id, typ string) Variable {
	v := NewVariable().(*variable)
	v.Id = id
	v.Type = typ
	if v.Id == "version" || v.Id == "deployment" || v.Id == "client" || v.Id == "project" || v.Id == "environment" {
		v.Type = v.Id
	}
	return v
}

func NewVariableFromDict(input UntypedVariable) (Variable, error) {
	str, err := yaml.Marshal(input)
	if err != nil {
		return nil, errors.New("Invalid input variable format: " + err.Error())
	}
	result := NewVariable().(*variable)
	err = yaml.Unmarshal(str, result)
	if err != nil {
		return nil, errors.New("Invalid input variable format: " + err.Error())
	}
	if result.Id == "" {
		return nil, errors.New("Missing 'id' field in variable")
	}
	if err = result.parseType(); err != nil {
		return nil, err
	}
	return result, nil
}

func (v *variable) GetType() string {
	return v.Type
}

func (v *variable) SetDefault(def interface{}) {
	v.Default = def
}
func (v *variable) HasDefault() bool {
	return v.Default != nil
}
func (v *variable) SetSensitive(s bool) {
	v.Sensitive = s
}
func (v *variable) SetVisible(s bool) {
	v.Visible = s
}
func (v *variable) SetDescription(desc string) {
	v.Description = desc
}

func (v *variable) SetOneOfItems(items []interface{}) {
	v.Items = items
}

func (v *variable) AskUserInput() interface{} {
	if v.Default != nil {
		return nil
	}
	if v.Type == "version" {
		return nil
	}
	if v.Type == "one_of" { // backwards compatible
		v.Type = "string"
	}
	if v.Type == "string" {
		return ""
	}
	if v.Type == "integer" {
		return 0
	}
	if v.Type == "list" {
		return []interface{}{}
	}
	return nil
}

func (v *variable) GetValue(variableCtx *map[string]interface{}, env *ScriptEnvironment) (interface{}, error) {
	var vars map[string]interface{}
	if variableCtx == nil {
		vars = map[string]interface{}{}
	} else {
		vars = *variableCtx
	}
	if v.Type == "one_of" { // backwards compatible
		v.Type = "string"
	}
	if v.Type == "string" || v.Type == "integer" || v.Type == "list" {
		var val interface{}
		val, ok := vars[v.Id]
		if !ok {
			if v.Default != nil {
				var err error
				val, err = v.validateDefault(env)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("Missing value for variable '%s'", v.Id)
			}
		}
		typ, err := variable_types.GetVariableType(v.Type)
		if err != nil {
			return nil, err
		}
		val, err = typ.Validate(val, v.Options)
		if err != nil {
			return nil, errors.New(err.Error() + " for variable '" + v.Id + "'")
		}
		return v.validateOneOf(val)
	} else if v.Type == "version" {
		script, _ := script.ParseScript("$this.version")
		result, err := script.Eval(env)
		if err != nil {
			panic(err)
		}
		return result.Value()
	} else if v.Type == "client" { // backwards compatibility
		script, _ := script.ParseScript("$this.project")
		result, err := script.Eval(env)
		if err != nil {
			panic(err)
		}
		return result.Value()
	} else if v.Type == "project" {
		script, _ := script.ParseScript("$this.project")
		result, err := script.Eval(env)
		if err != nil {
			panic(err)
		}
		return result.Value()
	} else if v.Type == "deployment" {
		script, _ := script.ParseScript("$this.deployment")
		result, err := script.Eval(env)
		if err != nil {
			panic(err)
		}
		return result.Value()
	} else if v.Type == "environment" {
		script, _ := script.ParseScript("$this.environment")
		result, err := script.Eval(env)
		if err != nil {
			panic(err)
		}
		return result.Value()
	}
	return nil, errors.New("Variable type " + v.Type + " not implemented")
}

func (v *variable) validateDefault(env *ScriptEnvironment) (interface{}, error) {
	switch v.Default.(type) {
	case (*string):
		return v.parseEvalAndGetValue(*v.Default.(*string), env)
	case string:
		return v.parseEvalAndGetValue(v.Default.(string), env)
	case []interface{}:
		lst := []interface{}{}
		for _, k := range v.Default.([]interface{}) {
			switch k.(type) {
			case string:
				val, err := v.parseEvalAndGetValue(k.(string), env)
				if err != nil {
					return nil, err
				} else {
					lst = append(lst, val)
				}
			default:
				lst = append(lst, k)
			}
		}
		return lst, nil
	}
	return nil, fmt.Errorf("Unexpected type '%T' for default field of variable '%s'", v.Default, v.Id)
}

func (v *variable) parseEvalAndGetValue(str string, env *ScriptEnvironment) (interface{}, error) {
	script, err := script.ParseScript(str)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse default field of variable '%s': %s in '%s'", v.Id, err.Error(), str)
	}
	result, err := script.Eval(env)
	if err != nil {
		return nil, fmt.Errorf("Couldn't run expression in default field of variable '%s': %s in '%s'", v.Id, err.Error(), str)
	}
	value, err := result.Value()
	if err != nil {
		return nil, fmt.Errorf("Couldn't run expression in default field of variable '%s': %s in '%s'", v.Id, err.Error(), str)
	}
	return value, nil

}

func (v *variable) validateOneOf(item interface{}) (interface{}, error) {
	items := v.Items
	if items == nil {
		return item, nil
	}
	for _, i := range items {
		if i == item {
			return item, nil
		}
	}
	oneOfString, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Expecting one of %s for variable '%s'", oneOfString, v.Id)
}

func (v *variable) parseType() error {
	if v.Type == "" || v.Type == "string" {
		switch v.Id {
		case
			"version",
			"project",
			"environment",
			"deployment",
			"client":
			v.Type = v.Id
		default:
			v.Type = "string"
		}
	}
	parsed, err := parsers.ParseVariableType(v.Type)
	if err != nil {
		return err
	}
	v.Type = parsed.Type
	v.Options = parsed.Options
	return nil
}
