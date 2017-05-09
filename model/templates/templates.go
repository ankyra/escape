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

package templates

import (
	"fmt"
	"github.com/ankyra/escape-client/model/script"
	"github.com/ankyra/escape-client/util"
	"github.com/cbroglie/mustache"
	"io/ioutil"
)

type template struct {
	File    string
	Target  string
	Scopes  []string
	Mapping map[string]interface{}
}

func NewTemplate() *template {
	return &template{
		Mapping: map[string]interface{}{},
		Scopes:  []string{"build", "deploy"},
	}
}

func NewTemplateFromString(file string) *template {
	template := NewTemplate()
	target, _ := util.SplitExt(file)
	return template.SetFile(file).SetTarget(target)
}

func NewTemplateWithMapping(file string, mapping map[string]interface{}) *template {
	return NewTemplateFromString(file).SetMapping(mapping)
}

func NewTemplateFromInterface(obj interface{}) (*template, error) {
	switch obj.(type) {
	case string:
		return NewTemplateFromString(obj.(string)), nil
	case map[interface{}]interface{}:
		mapObj := obj.(map[interface{}]interface{})
		resultMap := map[string]interface{}{}
		for key, val := range mapObj {
			switch key.(type) {
			case string:
				keyStr := key.(string)
				resultMap[keyStr] = val
			default:
				return nil, fmt.Errorf("Unexpected type '%T' for key in template dict.", key)
			}
		}
		return NewTemplateFromInterfaceMap(resultMap)
	}
	return nil, fmt.Errorf("Unexpected type '%T' for template", obj)
}

func (t *template) SetFileFromInterface(obj interface{}) error {
	file, ok := obj.(string)
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	t.SetFile(file)
	return nil
}

func (t *template) SetTargetFromInterface(obj interface{}) error {
	file, ok := obj.(string)
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	t.SetTarget(file)
	return nil
}

func (t *template) SetScopesFromInterface(obj interface{}) error {
	strScope, ok := obj.(string)
	if ok {
		t.SetScopes([]string{strScope})
		return nil
	}
	listScope, ok := obj.([]interface{})
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	scopes := []string{}
	for _, scopeObj := range listScope {
		scopeStr, ok := scopeObj.(string)
		if !ok {
			return fmt.Errorf("Unexpected type '%T'", obj)
		}
		scopes = append(scopes, scopeStr)
	}
	t.SetScopes(scopes)
	return nil
}

func (t *template) SetMappingFromInterface(obj interface{}) error {
	mapping, ok := obj.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	resultMap := map[string]interface{}{}
	for key, val := range mapping {
		keyStr, ok := key.(string)
		if !ok {
			return fmt.Errorf("Unexpected type '%T' for key", key)
		}
		resultMap[keyStr] = val
	}
	t.SetMapping(resultMap)
	return nil
}

func NewTemplateFromInterfaceMap(obj map[string]interface{}) (*template, error) {
	template := NewTemplate()
	for key, obj := range obj {
		switch key {
		case "file":
			if err := template.SetFileFromInterface(obj); err != nil {
				return nil, fmt.Errorf("%s in field '%s' of template dict.", err.Error(), key)
			}
		case "target":
			if err := template.SetTargetFromInterface(obj); err != nil {
				return nil, fmt.Errorf("%s in field '%s' of template dict.", err.Error(), key)
			}
		case "scopes":
			if err := template.SetScopesFromInterface(obj); err != nil {
				return nil, fmt.Errorf("%s in field '%s' of template dict.", err.Error(), key)
			}
		case "mapping":
			if err := template.SetMappingFromInterface(obj); err != nil {
				return nil, fmt.Errorf("%s in field '%s' of template dict.", err.Error(), key)
			}
		}
	}
	if template.Target == "" && template.File != "" {
		target, _ := util.SplitExt(template.File)
		template.SetTarget(target)
	}
	return template, nil
}

func (t *template) SetMapping(mapping map[string]interface{}) *template {
	t.Mapping = mapping
	return t
}
func (t *template) SetFile(file string) *template {
	t.File = file
	return t
}
func (t *template) SetTarget(file string) *template {
	t.Target = file
	return t
}
func (t *template) SetScopes(scopes []string) *template {
	t.Scopes = scopes
	return t
}

func (t *template) Render(stage string, env *script.ScriptEnvironment) error {
	if t.File == "" {
		return fmt.Errorf("Can't run template. Template file has not been defined (missing 'file' key in Escape plan?)")
	}
	if t.Target == "" {
		return fmt.Errorf("Can't run template. Template target has not been defined (empty 'target' key in Escape plan?)")
	}
	var inScope bool
	for _, scope := range t.Scopes {
		if scope == stage {
			inScope = true
			break
		}
	}
	if !inScope {
		return nil
	}
	result, err := t.renderToString(env)
	if err != nil {
		return fmt.Errorf("Failed to compile template %s: %s", t.File, err.Error())
	}
	if err := ioutil.WriteFile(t.Target, []byte(result), 0644); err != nil {
		return fmt.Errorf("Couldn't write output of template %s to %s: %s", t.File, t.Target, err.Error())
	}
	return nil
}

func (t *template) renderToString(env *script.ScriptEnvironment) (string, error) {
	mapping := map[string]interface{}{}
	for key, mappingValue := range t.Mapping {
		switch mappingValue.(type) {
		case string:
			scriptStr := mappingValue.(string)
			parsed, err := script.ParseScript(scriptStr)
			if err != nil {
				return "", fmt.Errorf("Couldn't parse script in template '%s' mapping key '%s': %s", t.File, key, err.Error())
			}
			evaled, err := parsed.Eval(env)
			if err != nil {
				return "", fmt.Errorf("Couldn't evaluate script in template '%s' mapping key '%s': %s", t.File, key, err.Error())
			}
			value, err := evaled.Value()
			if err != nil {
				return "", fmt.Errorf("Couldn't evaluate script in template '%s' mapping key '%s': %s", t.File, key, err.Error())
			}
			mapping[key] = value
		default:
			mapping[key] = mappingValue

		}
	}
	return mustache.RenderFile(t.File, mapping)
}
