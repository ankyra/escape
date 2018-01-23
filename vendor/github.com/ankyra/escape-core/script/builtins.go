/*
Copyright 2017, 2018 Ankyra

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
	"io/ioutil"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type StdlibFunc struct {
	Id     string
	Func   Script
	Doc    string
	ActsOn string
	Args   string
}

var trackMajorVersion = ShouldParse(`$func(v) { $v.split(".")[:1].join(".").concat(".@") }`)
var trackMinorVersion = ShouldParse(`$func(v) { $v.split(".")[:2].join(".").concat(".@") }`)
var trackPatchVersion = ShouldParse(`$func(v) { $v.split(".")[:3].join(".").concat(".@") }`)
var trackVersion = ShouldParse(`$func(v) { $v.concat(".@") }`)

var Stdlib = []StdlibFunc{
	StdlibFunc{"id", LiftFunction(builtinId), "Returns its argument", "everything", "parameter :: *"},
	StdlibFunc{"env_lookup", LiftFunction(builtinEnvLookup), "Lookup key in environment. Usually called implicitly when using '$'", "lists", "key :: string"},
	StdlibFunc{"concat", LiftFunction(builtinConcat), "Concatate stringable arguments", "strings", "v1 :: string, v2 :: string, ..."},
	StdlibFunc{"lower", ShouldLift(strings.ToLower), "Returns a copy of the string v with all Unicode characters mapped to their lower case", "strings", "v :: string"},
	StdlibFunc{"upper", ShouldLift(strings.ToUpper), "Returns a copy of the string v with all Unicode characters mapped to their upper case", "strings", "v :: string"},
	StdlibFunc{"title", ShouldLift(strings.ToTitle), "Returns a copy of the string v with all Unicode characters mapped to their title case", "strings", "v :: string"},
	StdlibFunc{"split", ShouldLift(strings.Split), "Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators. If sep is empty, Split splits after each UTF-8 sequence.", "strings", "sep :: string"},
	StdlibFunc{"join", ShouldLift(strings.Join), "Join concatenates the elements of a to create a single string. The separator string sep is placed between elements in the resulting string. ", "lists", "sep :: string"},
	StdlibFunc{"replace", ShouldLift(strings.Replace), "Replace returns a copy of the string s with the first n non-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to k+1 replacements for a k-rune string. If n < 0, there is no limit on the number of replacements.", "strings", "old :: string, new :: string, n :: integer"},
	StdlibFunc{"base64_encode", ShouldLift(base64.StdEncoding.EncodeToString), "Encode string to base64", "strings", ""},
	StdlibFunc{"base64_decode", ShouldLift(base64.StdEncoding.DecodeString), "Decode string from base64", "strings", ""},
	StdlibFunc{"trim", ShouldLift(strings.TrimSpace), "Returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode. ", "strings", ""},
	StdlibFunc{"list_index", LiftFunction(builtinListIndex), "Index a list at position `n`. Usually accessed implicitly using indexing syntax (eg. `list[0]`)", "lists", "n :: integer"},
	StdlibFunc{"list_slice", LiftFunction(builtinListSlice), "Slice a list. Usually accessed implicitly using slice syntax (eg. `list[0:5]`)", "lists", "i :: integer, j :: integer"},
	StdlibFunc{"add", ShouldLift(builtinAdd), "Add two integers", "integers", "y :: integer"},
	StdlibFunc{"timestamp", ShouldLift(builtinTimestamp), "Returns a UNIX timestamp", "", ""},
	StdlibFunc{"read_file", ShouldLift(builtinReadfile), "Read the contents of a file", "strings", ""},
	StdlibFunc{"track_major_version", trackMajorVersion, "Track major version", "strings", ""},
	StdlibFunc{"track_minor_version", trackMinorVersion, "Track minor version", "strings", ""},
	StdlibFunc{"track_patch_version", trackPatchVersion, "Track patch version", "strings", ""},
	StdlibFunc{"track_version", trackVersion, "Track version", "strings", ""},
}

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
	if err := builtinArgCheck(1, "id", inputValues); err != nil {
		return nil, err
	}
	return inputValues[0], nil
}

func builtinEnvLookup(env *ScriptEnvironment, inputValues []Script) (Script, error) {
	if err := builtinArgCheck(1, "env_lookup", inputValues); err != nil {
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
	if err := builtinArgCheck(2, "list_index", inputValues); err != nil {
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

func builtinReadfile(arg string) (string, error) {
	bytes, err := ioutil.ReadFile(arg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func builtinTimestamp() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func builtinAdd(x, y int) int {
	return x + y
}
