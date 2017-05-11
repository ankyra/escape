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
	"io/ioutil"
	"strings"
)

type scriptFuncType func(*ScriptEnvironment, []Script) (Script, error)

/*
   Lift
*/

func Lift(val interface{}) (Script, error) {
	if val == nil {
		return LiftString(""), nil
	}
	switch val.(type) {
	case string:
		return LiftString(val.(string)), nil
	case float64:
		return LiftInteger(int(val.(float64))), nil
	case int:
		return LiftInteger(val.(int)), nil
	case []interface{}:
		vals := []Script{}
		for _, k := range val.([]interface{}) {
			v, err := Lift(k)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}
		return LiftList(vals), nil
	}
	return nil, fmt.Errorf("Couldn't lift value of type '%T': %v", val, val)
}

/*
   Strings
*/
type stringAtom struct {
	String string
}

func LiftString(s string) Script {
	return &stringAtom{String: s}
}
func (s *stringAtom) Eval(env *ScriptEnvironment) (Script, error) {
	return s, nil
}
func (s *stringAtom) Value() (interface{}, error) {
	return s.String, nil
}
func (s *stringAtom) Type() ValueType {
	return NewType("string")
}

/*
   Integers
*/
type integerAtom struct {
	Integer int
}

func LiftInteger(i int) Script {
	return &integerAtom{Integer: i}
}
func (i *integerAtom) Eval(env *ScriptEnvironment) (Script, error) {
	return i, nil
}
func (i *integerAtom) Value() (interface{}, error) {
	return i.Integer, nil
}
func (i *integerAtom) Type() ValueType {
	return NewType("integer")
}

/*
   Lists
*/
type list struct {
	List []Script
}

func LiftList(l []Script) Script {
	return &list{List: l}
}
func (l *list) Eval(env *ScriptEnvironment) (Script, error) {
	return l, nil
}
func (l *list) Value() (interface{}, error) {
	return l.List, nil
}
func (l *list) Type() ValueType {
	return NewType("list")
}

/*
   Dicts
*/
type dict struct {
	Dict map[string]Script
}

func LiftDict(d map[string]Script) Script {
	return &dict{Dict: d}
}
func (d *dict) Eval(env *ScriptEnvironment) (Script, error) {
	return d, nil
}
func (d *dict) Value() (interface{}, error) {
	return d.Dict, nil
}
func (d *dict) Type() ValueType {
	return NewType("map")
}

/*
   Functions
*/
type function struct {
	Func scriptFuncType
}

func LiftFunction(f scriptFuncType) Script {
	return &function{
		Func: f,
	}
}

func (f *function) Eval(env *ScriptEnvironment) (Script, error) {
	return f, nil
}
func (f *function) Value() (interface{}, error) {
	return f.Func, nil
}
func (f *function) Type() ValueType {
	return NewType("func")
}

/*
   Apply
*/
type apply struct {
	To        Script
	Arguments []Script
}

func NewApply(to Script, args []Script) Script {
	return &apply{
		To:        to,
		Arguments: args,
	}
}
func (f *apply) Eval(env *ScriptEnvironment) (Script, error) {
	evaledTo, err := f.To.Eval(env)
	if err != nil {
		return nil, err
	}
	typ := evaledTo.Type()
	args := []Script{}
	for _, arg := range f.Arguments {
		evaledArg, err := arg.Eval(env)
		if err != nil {
			return nil, err
		}
		args = append(args, evaledArg)
	}
	if typ.IsFunc() {
		return f.evalFuncApply(evaledTo, args, env)
	} else if typ.IsMap() {
		return f.evalDictApply(evaledTo, args)
	} else if typ.IsString() {
		return f.evalStringApply(evaledTo, args)
	}
	return nil, fmt.Errorf("Expecting function, map or string for apply, but got '%s'", typ.Name())
}

func (f *apply) evalFuncApply(script Script, args []Script, env *ScriptEnvironment) (Script, error) {
	fun, err := script.Value()
	if err != nil {
		return nil, err
	}
	return fun.(scriptFuncType)(env, args)
}

func (f *apply) evalDictApply(dict Script, args []Script) (Script, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Expecting one argument in dict lookup call, but got '%d'", len(args))
	}
	arg := args[0]
	typ := arg.Type()
	if !typ.IsString() {
		return nil, fmt.Errorf("Expecting string argument in dict lookup call, but got '%s'", typ.Name())
	}
	key, _ := arg.Value()
	d, _ := dict.Value()
	result, ok := d.(map[string]Script)[key.(string)]
	if !ok {
		keys := []string{}
		for k, _ := range d.(map[string]Script) {
			keys = append(keys, k)
		}
		expects := strings.Join(keys, ", ")
		if len(keys) == 0 {
			expects = "target collection was empty"
		}
		return nil, fmt.Errorf("Field '%s' was not found (%s)", key, expects)
	}
	return result, nil
}

func (f *apply) evalStringApply(str Script, args []Script) (Script, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Expecting one argument in string call, but got '%d'", len(args))
	}
	arg := args[0]
	typ := arg.Type()
	if !typ.IsString() {
		return nil, fmt.Errorf("Expecting string argument in string call, but got '%s'", typ.Name())
	}
	s, _ := str.Value()
	fun, _ := arg.Value()
	if fun.(string) == "file" {
		result, err := builtinFileStringFunc(s.(string))
		if err != nil {
			return nil, err
		}
		return LiftString(result), nil
	}
	return nil, fmt.Errorf("Calling unknown function '%s' on string", fun.(string))
}

func (f *apply) Value() (interface{}, error) {
	return nil, fmt.Errorf("Function application can not be converted to Go value (forgot to eval?)")
}

func (f *apply) Type() ValueType {
	return f.To.Type() // TODO
}

func builtinFileStringFunc(str string) (string, error) {
	tmp, err := ioutil.TempFile("", "escape_input_")
	if err != nil {
		return "", fmt.Errorf("Could not create temporary file: %s", err.Error())
	}
	if _, err := tmp.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("Could not write to temporary file: %s", err.Error())
	}
	return tmp.Name(), nil
}
