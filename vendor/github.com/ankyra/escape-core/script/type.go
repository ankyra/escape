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

import ()

type ValueType interface {
	Name() string
	IsFunc() bool
	IsInteger() bool
	IsList() bool
	IsMap() bool
	IsString() bool
	IsLambda() bool
}

/*
   Expression types
*/
type valueType struct {
	Type string
}

func NewType(typ string) ValueType {
	return &valueType{Type: typ}
}
func (typ *valueType) Name() string {
	return typ.Type
}
func (typ *valueType) IsFunc() bool {
	return typ.Type == "func"
}
func (typ *valueType) IsMap() bool {
	return typ.Type == "map"
}
func (typ *valueType) IsString() bool {
	return typ.Type == "string"
}
func (typ *valueType) IsInteger() bool {
	return typ.Type == "integer"
}
func (typ *valueType) IsList() bool {
	return typ.Type == "list"
}
func (typ *valueType) IsLambda() bool {
	return typ.Type == "lambda"
}
