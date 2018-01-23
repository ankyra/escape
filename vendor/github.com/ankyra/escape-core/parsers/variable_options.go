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

package parsers

import (
	"errors"
	"strings"
)

func ParseOptions(str string) (map[string]interface{}, error) {
	options := map[string]interface{}{}
	for true {
		var setting map[string]interface{}
		str = GreedySpace(str)
		setting, str = parseExpression(str)
		if setting == nil {
			return nil, errors.New("Expecting expression key=value, got: '" + str + "'")
		}
		for k, v := range setting {
			options[k] = v
		}
		str = GreedySpace(str)
		if len(str) == 0 {
			return options, nil
		}
		if str[0] != ',' {
			return nil, errors.New("Expecting ',' got '" + str + "'")
		}
		str = str[1:]
	}
	return options, nil
}

func parseValue(str string) (interface{}, string) {
	val, rest := ParseInteger(str)
	if val == nil {
		return nil, rest
	}
	return *val, rest
}

func parseOperator(str string) (string, string) {
	str = GreedySpace(str)
	if len(str) == 0 {
		return "", str
	}
	if strings.HasPrefix(str, "=") {
		return "=", str[1:]
	}
	return "", str
}

func parseExpression(str string) (map[string]interface{}, string) {
	ident, rest := ParseIdent(str)
	if ident == "" {
		return nil, str
	}

	operator, rest := parseOperator(rest)
	if operator == "" {
		result := map[string]interface{}{ident: true}
		return result, rest
	}

	value, rest := parseValue(rest)
	if value == nil {
		return nil, str
	}
	result := map[string]interface{}{ident: value}
	return result, rest
}
