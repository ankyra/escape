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
	"io/ioutil"
	"path/filepath"

	"github.com/ankyra/escape-core/script"
	"github.com/cbroglie/mustache"
)

/*
Escape provides the Mustache templating language and integrates it with the
package's [Variables](/docs/reference/input-and-output-variables/), making for a quick
and easy way to render files at either build or deploy time.

## Escape Plan

Templates are configured in the Escape Plan under the
[`templates`](/docs/reference/escape-plan/#templates) field.

*/
type Template struct {
	// The file containing the template. This field is required.
	File string `json:"file"`

	// The target location for the rendered template. If the source location
	// specified in `file` has the `.tpl` extension this `target` will default
	// to source location minus that extension.
	//
	// For example: if `file` is `"hello.txt.tpl"` then the default value for
	// target will be `"hello.txt"`
	Target string `json:"target"`

	// A list of scopes (`build`, `deploy`) that defines during which stage(s)
	// the template should be rendered.
	Scopes []string `json:"scopes"`

	// This mapping can be used to relate template variables to Escape variables.
	Mapping map[string]interface{} `json:"mapping"`
}

func NewTemplate() *Template {
	return &Template{
		Mapping: map[string]interface{}{},
		Scopes:  []string{"build", "deploy"},
	}
}

func NewTemplateFromString(file string) *Template {
	template := NewTemplate()
	return template.SetFile(file).SetTarget(fileWithoutExtension(file))
}

func NewTemplateWithMapping(file string, mapping map[string]interface{}) *Template {
	return NewTemplateFromString(file).SetMapping(mapping)
}

func NewTemplateFromInterface(obj interface{}) (*Template, error) {
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

func (t *Template) SetFileFromInterface(obj interface{}) error {
	file, ok := obj.(string)
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	t.SetFile(file)
	return nil
}

func (t *Template) SetTargetFromInterface(obj interface{}) error {
	file, ok := obj.(string)
	if !ok {
		return fmt.Errorf("Unexpected type '%T'", obj)
	}
	t.SetTarget(file)
	return nil
}

func (t *Template) SetScopesFromInterface(obj interface{}) error {
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

func (t *Template) SetMappingFromInterface(obj interface{}) error {
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

func NewTemplateFromInterfaceMap(obj map[string]interface{}) (*Template, error) {
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
		template.SetTarget(fileWithoutExtension(template.File))
	}
	return template, nil
}

func (t *Template) SetMapping(mapping map[string]interface{}) *Template {
	t.Mapping = mapping
	return t
}
func (t *Template) SetFile(file string) *Template {
	t.File = file
	return t
}
func (t *Template) SetTarget(file string) *Template {
	t.Target = file
	return t
}
func (t *Template) SetScopes(scopes []string) *Template {
	t.Scopes = scopes
	return t
}

func (t *Template) Render(stage string, env *script.ScriptEnvironment) error {
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

func (t *Template) renderToString(env *script.ScriptEnvironment) (string, error) {
	mapping := map[string]interface{}{}
	for key, mappingValue := range t.Mapping {
		switch mappingValue.(type) {
		case string:
			scriptStr := mappingValue.(string)
			evaled, err := script.ParseAndEvalToGoValue(scriptStr, env)
			if err != nil {
				return "", fmt.Errorf("Error in template '%s' mapping key '%s': %s", t.File, key, err.Error())
			}
			mapping[key] = evaled
		default:
			mapping[key] = mappingValue

		}
	}
	mustache.AllowMissingVariables = false
	return mustache.RenderFile(t.File, mapping)
}

func (t *Template) InScope(scope string) bool {
	for _, s := range t.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func fileWithoutExtension(path string) (root string) {
	ext := filepath.Ext(path)
	root = path[:len(path)-len(ext)]
	return
}
