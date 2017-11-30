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

/*

Variables can be used to defined inputs and outputs for the build and
deployment stages. They can also be used to make [Errands](/docs/reference/errands/)
configurable.

Variables are strongly typed, which is checked at both build and deploy
time.  A task can't succeed if the required variables have not been
configured correctly.

## Escape Plan

Variables can be configured in the Escape Plan under the
[`inputs`](/docs/reference/escape-plan/#inputs),
[`build_inputs`](/docs/reference/escape-plan/#build_inputs),
[`deploy_inputs`](/docs/reference/escape-plan/#deploy_inputs) and
[`outputs`](/docs/reference/escape-plan/#outputs) fields.

*/
type Variable struct {
	// A unique name for this variable. Required field.
	Id string `json:"id"`

	// The variable type. Before executing any steps Escape will make sure that
	// all the values match the types that are set on the variables.
	//
	// One of: `string`, `list`, `integer`, `bool`.
	//
	// Default: `string`
	Type string `json:"type"`

	// A default value for this variable. This value will be used if no value
	// has been specified by the user.
	Default interface{} `json:"default,omitempty"`

	// A description of the variable.
	Description string `json:"description,omitempty"`

	// A friendly name for this variable for presentational purposes only.
	Friendly string `json:"friendly,omitempty"`

	// Control whether or not this variable should be visible when deploying
	// interactively. In other words: should the user be asked to input this
	// value?  It only really makes sense to set this to `true` if there a
	// `default` is set.
	Visible bool `json:"visible"`

	// Options that put more constraints on the type.
	Options map[string]interface{} `json:"options,omitempty"`

	// Is this sensitive data?
	Sensitive bool `json:"sensitive,omitempty"`

	// If set, this should contain all the valid values for this variable.
	Items interface{} `json:"items"`

	// Should the variables be evaluated before the dependencies are deployed?
	EvalBeforeDependencies bool `json:"eval_before_dependencies,omitempty"`

	// A list of scopes (`build`, `deploy`) that defines during which stage(s)
	// this variable should be active. You wouldn't usually use this field
	// directly, but use something like
	// [`build_inputs`](/docs/escape-plan/#build_inputs) or
	// [`deploy_inputs`](/docs/escape-plan/#deploy_inputs), which usually
	// express intent better.
	Scopes []string `json:"scopes"`
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
		//fmt.Errorf("The variable name '%s' is reserved", v.Id)
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
		switch i.(type) {
		case string:
			pv, err := v.parseEvalAndGetValue(i.(string), env)
			if err != nil {
				return nil, fmt.Errorf("In items field of variable '%s': %s", v.Id, err.Error())
			}
			if pv == item {
				return item, nil
			}
		case float64:
			_, itemInt := item.(int)
			if itemInt && int(i.(float64)) == item {
				return item, nil
			} else if i == item {
				return item, nil
			}
		default:
			if i == item {
				return item, nil
			}
		}
	}
	oneOfString, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Expecting one of %s for variable '%s', got: %v (%T)", oneOfString, v.Id, item, item)
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
