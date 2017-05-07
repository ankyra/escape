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

package script

import (
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
)

func RunScriptForCompileStep(script string, variableCtx map[string]ReleaseMetadata) (string, error) {
	parsedScript, err := ParseScript(script)
	if err != nil {
		return "", err
	}
	env := NewScriptEnvironmentForCompileStep(variableCtx)
	val, err := parsedScript.Eval(env)
	if err != nil {
		return "", err
	}
	if val.Type().IsString() {
		v, err := val.Value()
		if err != nil {
			return "", err
		}
		return v.(string), nil
	}
	if val.Type().IsInteger() {
		v, err := val.Value()
		if err != nil {
			return "", err
		}
		return string(v.(int)), nil
	}
	return "", fmt.Errorf("Expression '%s' did not return a string value", script)
}
