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
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	func_builtinId        = "__id"
	func_builtinEnvLookup = "__envLookup"
	func_builtinConcat    = "__concat"
	func_builtinToLower   = "__lower"
	func_builtinToUpper   = "__upper"
	func_builtinTitle     = "__title"
)

var builtinToLower = ShouldLift(strings.ToLower)
var builtinToUpper = ShouldLift(strings.ToUpper)
var builtinTitle = ShouldLift(strings.ToTitle)

func LiftFunc_string_to_string(f func(string) string) Script {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return LiftFunction(ToScriptFuncType_one_arg("builtin function "+name, wrap_one_arg_string_to_string(f)))
}

func wrap_one_arg_string_to_string(f func(string) string) func(interface{}) (interface{}, error) {
	return func(arg_ interface{}) (interface{}, error) {
		arg, ok := arg_.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string argument")
		}
		return f(arg), nil
	}
}

func ToScriptFuncType_one_arg(name string, f func(interface{}) (interface{}, error)) scriptFuncType {
	return func(env *ScriptEnvironment, args []Script) (Script, error) {
		if err := builtinArgCheck(1, name, args); err != nil {
			return nil, err
		}
		arg := args[0]
		if !IsStringAtom(arg) {
			return nil, fmt.Errorf("Expecting string argument in call to %s, but got %s", name, arg.Type().Name())
		}
		argument := ExpectStringAtom(arg)
		result, err := f(argument)
		if err != nil {
			return nil, err
		}
		return Lift(result)
	}
}

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
	if err := builtinArgCheck(1, func_builtinId, inputValues); err != nil {
		return nil, err
	}
	return inputValues[0], nil
}

func builtinEnvLookup(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if err := builtinArgCheck(1, func_builtinEnvLookup, inputValues); err != nil {
		return nil, err
	}
	arg := inputValues[0]
	if !IsStringAtom(arg) {
		return nil, fmt.Errorf("Expecting string argument in environment lookup call, but got '%s'", arg.Type().Name())
	}
	key := ExpectStringAtom(arg)
	val, ok := (*env)[key]
	if !ok {
		return nil, fmt.Errorf("Field '%s' was not found in environment.", key)
	}
	return val, nil
}

func builtinConcat(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	result := ""
	for _, val := range inputValues {
		if IsStringAtom(val) {
			result += ExpectStringAtom(val)
		} else if IsIntegerAtom(val) {
			result += strconv.Itoa(ExpectIntegerAtom(val))
		} else {
			return nil, fmt.Errorf("Can't concatenate value of type %s", val.Type().Name())
		}
	}
	return LiftString(result), nil
}
