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
	"strconv"
)

var integerType = NewUserManagedVariableType("integer", validateInt)

func validateInt(value interface{}, options map[string]interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return value.(int), nil
	case float64:
		return int(value.(float64)), nil
	case string:
		i, err := strconv.Atoi(value.(string))
		if err != nil {
			return nil, fmt.Errorf("Expecting 'integer' value, but got 'string'")
		}
		return i, nil
	}
	return nil, fmt.Errorf("Expecting 'integer' value, but got '%T'", value)
}
