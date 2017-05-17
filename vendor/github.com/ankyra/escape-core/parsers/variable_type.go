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

package parsers

import (
	"errors"
	"github.com/ankyra/escape-core/variables/variable_types"
	"strings"
)

type ParsedVariableType struct {
	Type    string
	Options map[string]interface{}
}

func ParseVariableType(str string) (*ParsedVariableType, error) {
	result := &ParsedVariableType{}
	str = strings.TrimSpace(str)
	if strings.Contains(str, "[") && strings.HasSuffix(str, "]") {
		parts := strings.Split(str, "[")
		result.Type = parts[0]
		rest := strings.Join(parts[1:], "[")
		rest = strings.TrimSuffix(rest, "]")
		options, err := ParseOptions(rest)
		if err != nil {
			return nil, err
		}
		result.Options = options
	} else {
		result.Type = str
	}
	if !isValidVariableType(result.Type) {
		return nil, errors.New("Unknown variable type: " + result.Type)
	}
	return result, nil
}

func isValidVariableType(typ string) bool {
	for _, supported := range variable_types.GetSupportedTypes() {
		if supported == typ {
			return true
		}
	}
	return false
}
