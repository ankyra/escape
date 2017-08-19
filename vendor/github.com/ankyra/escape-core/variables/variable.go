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

package variables

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/variables/variable_types"
	"gopkg.in/yaml.v2"
)

type Variable struct {
	Id                     string                 `json:"id"`
	Type                   string                 `json:"type"`
	Default                interface{}            `json:"default,omitempty"`
	Description            string                 `json:"description,omitempty"`
	Friendly               string                 `json:"friendly,omitempty"`
	Visible                bool                   `json:"visible"`
	Options                map[string]interface{} `json:"options,omitempty"`
	Sensitive              bool                   `json:"sensitive,omitempty"`
	Items                  interface{}            `json:"items"`
	EvalBeforeDependencies bool                   `json:"eval_before_dependencies,omitempty"`
	Scopes                 []string               `json:"scopes"`
}

type UntypedVariable map[interface{}]interface{}

func NewVariable() *Variable {
	return &Variable{
		Visible:                true,
		EvalBeforeDependencies: true,
		Scopes:                 []string{"build", "deploy"},
	}
}

func NewVariableFromInterface(v interface{}) (*Variable, error) {
	switch v.(type) {
	case string:
		return NewVariableFromString(v.(string), "string")
	case map[interface{}]interface{}:
		return NewVariableFromDict(v.(map[interface{}]interface{}))
	}
	return nil, fmt.Errorf("Expecting dict or string type")
}

func NewVariableFromString(id, typ string) (*Variable, error) {
	v := NewVariable()
	v.Id = id
	v.Type = typ
	return v, v.Validate()
}

func NewVariableFromDict(input UntypedVariable) (*Variable, error) {
	str, err := yaml.Marshal(input)
	if err != nil {
		return nil, errors.New("Invalid input variable format: " + err.Error())
	}
	result := NewVariable()
	if err = yaml.Unmarshal(str, result); err != nil {
		return nil, errors.New("Invalid input variable format: " + err.Error())
	}
	if result.Id == "" {
		return nil, errors.New("Missing 'id' field in variable")
	}
	if err = result.parseType(); err != nil {
		return nil, err
	}
	return result, result.Validate()
}

func (v *Variable) Validate() error {
	if v.Id == "" {
		return fmt.Errorf("Variable object is missing an 'id'")
	}
	_, rest := parsers.ParseIdent(v.Id)
	if strings.TrimSpace(rest) != "" {
		return fmt.Errorf("Invalid variable id format '%s'", v.Id)
	}
	v.Id = strings.TrimSpace(v.Id)
	if strings.HasPrefix(strings.ToUpper(v.Id), "PREVIOUS_") {
		return fmt.Errorf("Invalid variable format '%s'. Variable is not allowed to start with '%s'",
			v.Id, v.Id[:len("PREVIOUS_")])
	}
	if v.Scopes == nil || len(v.Scopes) == 0 {
		v.Scopes = []string{"build", "deploy"}
	}
	if variable_types.VariableIdIsReservedType(v.Id) {
		return fmt.Errorf("The variable name '%s' is reserved", v.Id)
	}
	return nil
}

func (v *Variable) InScope(scope string) bool {
	for _, s := range v.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (v *Variable) HasDefault() bool {
	return v.Default != nil
}

func (v *Variable) AskUserInput() interface{} {
	if v.HasDefault() {
		return nil
	}
	if v.Type == "version" {
		return nil
	}
	if v.Type == "string" {
		return ""
	}
	if v.Type == "integer" {
		return 0
	}
	if v.Type == "bool" {
		return false
	}
	if v.Type == "list" {
		return []interface{}{}
	}
	return nil
}

func (v *Variable) GetValue(variableCtx *map[string]interface{}, env *script.ScriptEnvironment) (interface{}, error) {
	typ, err := variable_types.GetVariableType(v.Type)
	if err != nil {
		return nil, err
	}
	if typ.UserCanOverride {
		return v.getValueForUserManagedVariable(variableCtx, env)
	}
	return script.ParseAndEvalToGoValue(typ.Script, env)
}

func (v *Variable) getValueForUserManagedVariable(variableCtx *map[string]interface{}, env *script.ScriptEnvironment) (interface{}, error) {
	typ, err := variable_types.GetVariableType(v.Type)
	if err != nil {
		return nil, err
	}
	val, err := v.getValue(variableCtx, env)
	if err != nil {
		return nil, err
	}
	val, err = typ.Validate(val, v.Options)
	if err != nil {
		return nil, errors.New(err.Error() + " for variable '" + v.Id + "'")
	}
	return v.validateOneOf(env, val)
}

func (v *Variable) getValue(variableCtx *map[string]interface{}, env *script.ScriptEnvironment) (interface{}, error) {
	if variableCtx == nil {
		variableCtx = &map[string]interface{}{}
	}
	val, ok := (*variableCtx)[v.Id]
	if ok {
		return val, nil
	}
	return v.getDefaultValue(env)
}

func (v *Variable) getDefaultValue(env *script.ScriptEnvironment) (interface{}, error) {
	if v.Default == nil {
		return nil, fmt.Errorf("Missing value for variable '%s'", v.Id)
	}
	switch v.Default.(type) {
	case int:
		return v.Default.(int), nil
	case float64:
		return v.Default, nil
	case bool:
		return v.Default, nil
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

func (v *Variable) parseEvalAndGetValue(str string, env *script.ScriptEnvironment) (interface{}, error) {
	result, err := script.ParseAndEvalToGoValue(str, env)
	if err != nil {
		return nil, fmt.Errorf("Couldn't run expression in default field of variable '%s': %s in '%s'", v.Id, err.Error(), str)
	}
	return result, nil
}

func (v *Variable) validateOneOf(env *script.ScriptEnvironment, item interface{}) (interface{}, error) {
	items := v.Items
	return v.validateOneOfInterface(env, item, items)
}

func (v *Variable) validateOneOfInterface(env *script.ScriptEnvironment, item interface{}, items interface{}) (interface{}, error) {
	if items == nil {
		return item, nil
	}
	switch items.(type) {
	case string:
		pv, err := v.parseEvalAndGetValue(items.(string), env)
		if err != nil {
			return nil, fmt.Errorf("In items field of variable '%s': %s", v.Id, err.Error())
		}
		_, isString := pv.(string)
		if isString {
			if pv == item {
				return item, nil
			}
			return nil, fmt.Errorf("Unexpected value '%s' for variable '%s', only '%s' is allowed", item, v.Id, pv)
		}
		return v.validateOneOfInterface(env, item, pv)
	case []interface{}:
		return v.validateOneOfList(env, item, items.([]interface{}))
	}
	return nil, fmt.Errorf("Unexpected type '%T' for 'items' field of variable '%s'", items, v.Id)
}

func (v *Variable) validateOneOfList(env *script.ScriptEnvironment, item interface{}, items []interface{}) (interface{}, error) {
	for _, i := range items {
		_, isString := i.(string)
		if isString {
			pv, err := v.parseEvalAndGetValue(i.(string), env)
			if err != nil {
				return nil, fmt.Errorf("In items field of variable '%s': %s", v.Id, err.Error())
			}
			i = pv
		}
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

func (v *Variable) parseType() error {
	if v.Type == "" {
		v.Type = "string"
		if variable_types.VariableIdIsReservedType(v.Id) {
			v.Type = v.Id
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
