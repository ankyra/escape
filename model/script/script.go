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
)

type Script interface {
	Eval(*ScriptEnvironment) (Script, error)
	Value() (interface{}, error)
	Type() ValueType
}

func EvalToGoValue(script Script, env *ScriptEnvironment) (interface{}, error) {
	evaled, err := script.Eval(env)
	if err != nil {
		return nil, err
	}
	return evaled.Value()
}

func ParseAndEvalToGoValue(scriptStr string, env *ScriptEnvironment) (interface{}, error) {
	parsed, err := ParseScript(scriptStr)
	if err != nil {
		return "", err
	}
	evaled, err := EvalToGoValue(parsed, env)
	if err != nil {
		return "", fmt.Errorf("Failed to evaluate '%s': %s", scriptStr, err.Error())
	}
	return evaled, nil
}
