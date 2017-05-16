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

package variable_types

import (
	"fmt"
)

var versionType = NewMagicVariable("version", "$this.version")
var clientType = NewMagicVariable("client", "$this.client")
var projectType = NewMagicVariable("project", "$this.project")
var deploymentType = NewMagicVariable("deployment", "$this.deployment")
var environmenType = NewMagicVariable("environment", "$this.environment")

var knownTypes = []*VariableType{stringType, boolType, integerType, listType,
	versionType, clientType, projectType, deploymentType, environmenType}

type Validator func(value interface{}, options map[string]interface{}) (interface{}, error)

type VariableType struct {
	Type            string
	UserCanOverride bool
	Script          string
	Validate        Validator
}

func NewUserManagedVariableType(typ string, validate Validator) *VariableType {
	return &VariableType{
		Type:            typ,
		UserCanOverride: true,
		Validate:        validate,
	}
}

func NewMagicVariable(typ string, script string) *VariableType {
	return &VariableType{
		Type:   typ,
		Script: script,
	}
}

func GetVariableType(typ string) (*VariableType, error) {
	for _, varType := range knownTypes {
		if varType.Type == typ {
			return varType, nil
		}
	}
	return nil, fmt.Errorf("Unknown variable type '%s'", typ)
}

func VariableIdIsReservedType(typ string) bool {
	for _, varType := range knownTypes {
		if varType.Type == typ {
			return !varType.UserCanOverride
		}
	}
	return false
}

func GetSupportedTypes() []string {
	result := make([]string, len(knownTypes))
	for i, varType := range knownTypes {
		result[i] = varType.Type
	}
	return result
}
