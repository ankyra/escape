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
	"encoding/base64"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	func_builtinId           = "__id"
	func_builtinEnvLookup    = "__envLookup"
	func_builtinConcat       = "__concat"
	func_builtinToLower      = "__lower"
	func_builtinToUpper      = "__upper"
	func_builtinTitle        = "__title"
	func_builtinSplit        = "__split"
	func_builtinJoin         = "__join"
	func_builtinReplace      = "__replace"
	func_builtinBase64Encode = "__base64_encode"
	func_builtinBase64Decode = "__base64_decode"
	func_builtinTrim         = "__trim"
	func_builtinListIndex    = "__list_index"
	func_builtinListSlice    = "__list_slice"
)

var builtinToLower = ShouldLift(strings.ToLower)
var builtinToUpper = ShouldLift(strings.ToUpper)
var builtinTitle = ShouldLift(strings.ToTitle)
var builtinSplit = ShouldLift(strings.Split)
var builtinJoin = ShouldLift(strings.Join)
var builtinReplace = ShouldLift(strings.Replace)
var builtinTrim = ShouldLift(strings.TrimSpace)
var builtinBase64Encode = ShouldLift(base64.StdEncoding.EncodeToString)
var builtinBase64Decode = ShouldLift(base64.StdEncoding.DecodeString)

func LiftGoFunc(f interface{}) Script {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	typ := reflect.TypeOf(f)
	nInputs := typ.NumIn()
	nOutputs := typ.NumOut()
	scriptFunc := func(env *ScriptEnvironment, args []Script) (Script, error) {
		if err := builtinArgCheck(nInputs, name, args); err != nil {
			return nil, err
		}

		goArgs := []reflect.Value{}
		for i := 0; i < nInputs; i++ {
			argType := typ.In(i)
			arg := args[i]

			if argType.Kind() == reflect.String {
				if !IsStringAtom(arg) {
					return nil, fmt.Errorf("Expecting string argument in call to %s, but got %s", name, arg.Type().Name())
				} else {
					goArgs = append(goArgs, reflect.ValueOf(ExpectStringAtom(arg)))
				}
			} else if argType.Kind() == reflect.Int {
				if !IsIntegerAtom(arg) {
					return nil, fmt.Errorf("Expecting integer argument in call to %s, but got %s", name, arg.Type().Name())
				} else {
					goArgs = append(goArgs, reflect.ValueOf(ExpectIntegerAtom(arg)))
				}
			} else if argType.Kind() == reflect.Slice {
				if !IsListAtom(arg) {
					if argType.Elem().Kind() == reflect.Uint8 && IsStringAtom(arg) {
						goArgs = append(goArgs, reflect.ValueOf([]byte(ExpectStringAtom(arg))))
					} else {
						return nil, fmt.Errorf("Expecting list argument in call to %s, but got %s", name, arg.Type().Name())
					}
				} else {
					lst := ExpectListAtom(arg) // []Script
					if argType.Elem().Kind() == reflect.String {
						strArg := []string{}
						for k := 0; k < len(lst); k++ {
							if !IsStringAtom(lst[k]) {
								return nil, fmt.Errorf("Expecting string value in list in call to %s, but got %s", name, arg.Type().Name())
							} else {
								strArg = append(strArg, ExpectStringAtom(lst[k]))
							}
						}
						goArgs = append(goArgs, reflect.ValueOf(strArg))
					} else {
						return nil, fmt.Errorf("Unsupported slice type in function %s", name)
					}
				}
			} else {
				return nil, fmt.Errorf("Unsupported argument type '%s' in function %s", argType.Kind(), name)
			}
		}

		outputs := reflect.ValueOf(f).Call(goArgs)
		if nOutputs == 1 {
			return Lift(outputs[0].Interface())
		}
		if nOutputs == 2 {
			_, isError := outputs[1].Interface().(error)
			if isError && outputs[1].Interface() != nil {
				return nil, fmt.Errorf("Error in call to %s: %s", name, outputs[1].Interface().(error))
			}
			return Lift(outputs[0].Interface())
		}
		if nOutputs != 1 {
			return nil, fmt.Errorf("Go functions with multiple outputs are not supported at this time")
		}
		return Lift(outputs[0].Interface())
	}
	return LiftFunction(scriptFunc)
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

func builtinListIndex(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if err := builtinArgCheck(2, func_builtinListIndex, inputValues); err != nil {
		return nil, err
	}
	lstArg := inputValues[0]
	if !IsListAtom(lstArg) {
		return nil, fmt.Errorf("Expecting list argument in list index call, but got '%s'", lstArg.Type().Name())
	}
	indexArg := inputValues[1]
	if !IsIntegerAtom(indexArg) {
		return nil, fmt.Errorf("Expecting integer argument in list index call, but got '%s'", indexArg.Type().Name())
	}
	lst := ExpectListAtom(inputValues[0])
	index := ExpectIntegerAtom(inputValues[1])
	if index < 0 || index >= len(lst) {
		return nil, fmt.Errorf("Index '%d' out of range (len: %d)", index, len(lst))
	}
	return Lift(lst[index])
}

func builtinListSlice(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if len(inputValues) < 2 || len(inputValues) > 3 {
		return nil, fmt.Errorf("Expecting at least %d argument(s) (but not more than 3) in call to '%s', got %d",
			2, "list slice", len(inputValues))
	}
	lstArg := inputValues[0]
	if !IsListAtom(lstArg) {
		return nil, fmt.Errorf("Expecting list argument in list slice call, but got '%s'", lstArg.Type().Name())
	}
	indexArg := inputValues[1]
	if !IsIntegerAtom(indexArg) {
		return nil, fmt.Errorf("Expecting integer argument in list slice call, but got '%s'", indexArg.Type().Name())
	}
	lst := ExpectListAtom(inputValues[0])
	index := ExpectIntegerAtom(inputValues[1])

	if len(inputValues) == 3 {
		endSliceArg := inputValues[2]
		if !IsIntegerAtom(endSliceArg) {
			return nil, fmt.Errorf("Expecting integer argument in list slice call, but got '%s'", endSliceArg.Type().Name())
		}
		endIndex := ExpectIntegerAtom(inputValues[2])
		if endIndex < 0 {
			endIndex = len(lst) + endIndex
		}
		return Lift(lst[index:endIndex])
	}
	return Lift(lst[index:])
}
