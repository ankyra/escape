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

/*
   Builtins
*/
func builtinArgCheck(expected int, funcName string, inputValues []Script) error {
	if len(inputValues) != expected {
		return fmt.Errorf("Expecting %d argument(s) in call to '%s', got %d",
			expected, funcName, len(inputValues))
	}
	return nil
}

func builtinId(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if err := builtinArgCheck(1, "id", inputValues); err != nil {
		return nil, err
	}
	return inputValues[0], nil
}

func builtinEnvLookup(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if err := builtinArgCheck(1, "__env_lookup__", inputValues); err != nil {
		return nil, err
	}
	arg := inputValues[0]
	if !arg.Type().IsString() {
		return nil, fmt.Errorf("Expecting string argument in environment lookup call, but got '%s'", arg.Type().Name())
	}
	key, _ := arg.Value()
	val, ok := (*env)[key.(string)]
	if !ok {
		return nil, fmt.Errorf("Field '%s' was not found in environment.", key.(string))
	}
	return val, nil
}
