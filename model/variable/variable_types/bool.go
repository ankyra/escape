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

var boolType = NewUserManagedVariableType("bool", validateBool)

func validateBool(value interface{}, options map[string]interface{}) (interface{}, error) {
	switch value.(type) {
	case bool:
		return value.(bool), nil
	case int:
		return value.(int) == 1, nil
	case string:
		return value.(string) == "1", nil
	}
	return nil, fmt.Errorf("Expecting 'bool' value, but got '%T'", value)
}
