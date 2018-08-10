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
	case []byte:
		return LiftString(string(val.([]byte))), nil
	case bool:
		return LiftBool(val.(bool)), nil
	case float64:
		return LiftInteger(int(val.(float64))), nil
	case int:
		return LiftInteger(val.(int)), nil
	case Script:
		return val.(Script), nil
	case map[string]Script:
		return LiftDict(val.(map[string]Script)), nil
	case []Script:
		return LiftList(val.([]Script)), nil
	case scriptFuncType:
		return LiftFunction(val.(scriptFuncType)), nil
	case func() string:
		return LiftGoFunc(val), nil
	case func([]byte) string:
		return LiftGoFunc(val), nil
	case func(string) ([]byte, error):
		return LiftGoFunc(val), nil
	case func(string) (string, error):
		return LiftGoFunc(val), nil
	case func(string) bool:
		return LiftGoFunc(val), nil
	case func(string) string:
		return LiftGoFunc(val), nil
	case func(string, string) []string:
		return LiftGoFunc(val), nil
	case func(string, string) string:
		return LiftGoFunc(val), nil
	case func(string, string, string, int) string:
		return LiftGoFunc(val), nil
	case func([]string, string) string:
		return LiftGoFunc(val), nil
	case func(int) int:
		return LiftGoFunc(val), nil
	case func(int, int) int:
		return LiftGoFunc(val), nil
	case []string:
		vals := []Script{}
		for _, k := range val.([]string) {
			vals = append(vals, LiftString(k))
		}
		return LiftList(vals), nil
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
	case map[string]interface{}:
		resultMap := map[string]Script{}
		for key, val := range val.(map[string]interface{}) {
			v, err := Lift(val)
			if err != nil {
				return nil, err
			}
			resultMap[key] = v
		}
		return LiftDict(resultMap), nil
	case map[interface{}]interface{}:
		resultMap := map[string]Script{}
		for k, val := range val.(map[interface{}]interface{}) {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Expecting string key for dictionary type, but got %T", k)
			}
			v, err := Lift(val)
			if err != nil {
				return nil, err
			}
			resultMap[key] = v
		}
		return LiftDict(resultMap), nil
	}
	return nil, fmt.Errorf("Couldn't lift value of type '%T': %v", val, val)
}

func ShouldLift(v interface{}) Script {
	result, err := Lift(v)
	if err != nil {
		panic(err)
	}
	return result
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
func (s *stringAtom) Equals(s2 Script) bool {
	if !s2.Type().IsString() {
		return false
	}
	s2Str := ExpectStringAtom(s2)
	return s.String == s2Str
}
func IsStringAtom(s Script) (ok bool) {
	_, ok = s.(*stringAtom)
	return ok
}
func ExpectStringAtom(s Script) string {
	if IsStringAtom(s) {
		return s.(*stringAtom).String
	}
	panic("Expecting string type, got " + s.Type().Name())
}

/*
   Booleans
*/
type boolAtom struct {
	Bool bool
}

func LiftBool(b bool) Script {
	return &boolAtom{Bool: b}
}

func (b *boolAtom) Eval(env *ScriptEnvironment) (Script, error) {
	return b, nil
}
func (i *boolAtom) Value() (interface{}, error) {
	return i.Bool, nil
}
func (i *boolAtom) Type() ValueType {
	return NewType("bool")
}
func (s *boolAtom) Equals(s2 Script) bool {
	if !s2.Type().IsBool() {
		return false
	}
	s2Val := ExpectBoolAtom(s2)
	return s.Bool == s2Val
}

func IsBoolAtom(s Script) (ok bool) {
	_, ok = s.(*boolAtom)
	return ok
}
func ExpectBoolAtom(s Script) bool {
	if IsBoolAtom(s) {
		return s.(*boolAtom).Bool
	}
	panic("Expecting bool type, got " + s.Type().Name())
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
func (s *integerAtom) Equals(s2 Script) bool {
	if !s2.Type().IsInteger() {
		return false
	}
	s2Val := ExpectIntegerAtom(s2)
	return s.Integer == s2Val
}
func IsIntegerAtom(s Script) (ok bool) {
	_, ok = s.(*integerAtom)
	return ok
}
func ExpectIntegerAtom(s Script) int {
	if IsIntegerAtom(s) {
		return s.(*integerAtom).Integer
	}
	panic("Expecting integer type, got " + s.Type().Name())
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
	result := []Script{}
	for _, i := range l.List {
		val, err := i.Eval(env)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return LiftList(result), nil
}
func (l *list) Value() (interface{}, error) {
	result := []interface{}{}
	for _, i := range l.List {
		val, err := i.Value()
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}
func (l *list) Type() ValueType {
	return NewType("list")
}
func (s *list) Equals(s2 Script) bool {
	if !s2.Type().IsList() {
		return false
	}
	s2Val := ExpectListAtom(s2)
	if len(s2Val) != len(s.List) {
		return false
	}
	for i := 0; i < len(s.List); i++ {
		if !s.List[i].Equals(s2Val[i]) {
			return false
		}
	}
	return true
}
func IsListAtom(s Script) (ok bool) {
	_, ok = s.(*list)
	return ok
}
func ExpectListAtom(s Script) []Script {
	if IsListAtom(s) {
		return s.(*list).List
	}
	panic("Expecting list type, got " + s.Type().Name())
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
func (s *dict) Equals(s2 Script) bool {
	if !s2.Type().IsMap() {
		return false
	}
	s2Val := ExpectDictAtom(s2)
	if len(s2Val) != len(s.Dict) {
		return false
	}
	for key, value := range s.Dict {
		value2, found := s2Val[key]
		if !found {
			return false
		}
		if !value.Equals(value2) {
			return false
		}
	}
	return true
}
func IsDictAtom(s Script) bool {
	_, ok := s.(*dict)
	return ok
}
func ExpectDictAtom(s Script) map[string]Script {
	if IsDictAtom(s) {
		return s.(*dict).Dict
	}
	panic("Expecting dict type, got " + s.Type().Name())
}
func ExpectDict(s Script) map[string]interface{} {
	result := map[string]interface{}{}
	dict := ExpectDictAtom(s)
	for key, val := range dict {
		var err error
		result[key], err = val.Value()
		if err != nil {
			panic(err.Error())
		}
	}
	return result
}

/*
   Functions
*/
type Function struct {
	Func scriptFuncType
	Doc  string
}

func LiftFunction(f scriptFuncType) Script {
	return &Function{
		Func: f,
	}
}
func (f *Function) Eval(env *ScriptEnvironment) (Script, error) {
	return f, nil
}
func (f *Function) Value() (interface{}, error) {
	return f.Func, nil
}
func (f *Function) Type() ValueType {
	return NewType("func")
}
func (s *Function) Equals(s2 Script) bool {
	return false
}
func IsFunctionAtom(s Script) bool {
	_, ok := s.(*Function)
	return ok
}
func ExpectFunction(s Script) *Function {
	if IsFunctionAtom(s) {
		return s.(*Function)
	}
	panic("Expecting function type, got " + s.Type().Name())
}
func ExpectFunctionAtom(s Script) scriptFuncType {
	if IsFunctionAtom(s) {
		return s.(*Function).Func
	}
	panic("Expecting function type, got " + s.Type().Name())
}

/*
   Lambda
*/
type lambda struct {
	Arguments []string
	Body      Script
}

func NewLambda(args []string, body Script) Script {
	return &lambda{
		Arguments: args,
		Body:      body,
	}
}

func (l *lambda) Eval(env *ScriptEnvironment) (Script, error) {
	return l, nil
}

func (l *lambda) Type() ValueType {
	return NewType("lambda")
}

func (l *lambda) Value() (interface{}, error) {
	return nil, fmt.Errorf("Function application can not be converted to Go value (forgot to eval?)")
}
func (s *lambda) Equals(s2 Script) bool {
	return false
}

func IsLambdaAtom(v Script) bool {
	_, ok := v.(*lambda)
	return ok
}
func ExpectLambdaAtom(v Script) *lambda {
	if IsLambdaAtom(v) {
		return v.(*lambda)
	}
	panic("Expecting lambda function, got " + v.Type().Name())
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
func IsApplyAtom(s Script) bool {
	_, ok := s.(*apply)
	return ok
}
func ExpectApplyAtom(s Script) *apply {
	if IsApplyAtom(s) {
		return s.(*apply)
	}
	panic("Expecting function application, got " + s.Type().Name())
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
	} else if typ.IsLambda() {
		return f.evalLambdaApply(evaledTo, args, env)
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

func (f *apply) evalLambdaApply(script Script, args []Script, env *ScriptEnvironment) (Script, error) {
	lambda := ExpectLambdaAtom(script)
	if len(lambda.Arguments) != len(args) {
		return nil, fmt.Errorf("Argument arity mismatch. Expecting %d arguments, got %d.", len(lambda.Arguments), len(args))
	}
	if env == nil {
		env = NewScriptEnvironment()
	}
	newEnv := map[string]Script{}
	for key, val := range *env {
		newEnv[key] = val
	}
	globals, found := newEnv["$"]
	if !found {
		globals = LiftDict(map[string]Script{})
	}
	globalsDict := ExpectDictAtom(globals)
	newGlobals := map[string]Script{}
	for key, val := range globalsDict {
		newGlobals[key] = val
	}
	for ix, variable := range lambda.Arguments {
		newGlobals[variable] = args[ix]
	}
	newEnv["$"] = LiftDict(newGlobals)
	return lambda.Body.Eval(NewScriptEnvironmentFromMap(newEnv))
}

func (f *apply) evalDictApply(dict Script, args []Script) (Script, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Expecting one argument in dict lookup call, but got '%d'", len(args))
	}
	arg := args[0]
	if !IsStringAtom(arg) {
		return nil, fmt.Errorf("Expecting string argument in dict lookup call, but got '%s'", arg.Type().Name())
	}
	key := ExpectStringAtom(arg)
	d := ExpectDictAtom(dict)
	result, ok := d[key]
	if !ok {
		keys := []string{}
		for k, _ := range d {
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
	if !IsStringAtom(arg) {
		return nil, fmt.Errorf("Expecting string argument in string call, but got '%s'", arg.Type().Name())
	}
	applyTo := ExpectStringAtom(str)
	fun := ExpectStringAtom(arg)
	if fun == "file" {
		result, err := builtinFileStringFunc(applyTo)
		if err != nil {
			return nil, err
		}
		return LiftString(result), nil
	}
	return nil, fmt.Errorf("Calling unknown function '%s' on string", fun)
}

func (f *apply) Value() (interface{}, error) {
	return nil, fmt.Errorf("Function application can not be converted to Go value (forgot to eval?)")
}

func (f *apply) Type() ValueType {
	return f.To.Type() // TODO
}

func (s *apply) Equals(s2 Script) bool {
	return false
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
