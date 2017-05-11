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

type stringVarType struct {
}

func NewStringVariableType() VariableType {
	return &stringVarType{}
}

func (s *stringVarType) Validate(value interface{}, options map[string]interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return value.(string), nil
	}
	return "", fmt.Errorf("Expecting 'string' value, but got '%T'", value)
}
