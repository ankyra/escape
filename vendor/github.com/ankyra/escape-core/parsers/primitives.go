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
	"strconv"
)

func GreedySpace(str string) string {
	for len(str) > 0 && str[0] == ' ' {
		str = str[1:]
	}
	return str
}

func ParseIdent(str string) (string, string) {
	str = GreedySpace(str)
	result := []rune{}
	inWord := false
	for _, c := range str {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			result = append(result, c)
			inWord = true
		} else if inWord && ((c >= '0' && c <= '9') || c == '-' || c == '_') {
			result = append(result, c)
		} else {
			break
		}
	}
	if len(result) == 0 {
		return "", str
	}
	return string(result), str[len(result):]
}

func ParseInteger(str string) (*int, string) {
	str = GreedySpace(str)
	result := []rune{}
	for _, c := range str {
		if c >= '0' && c <= '9' {
			result = append(result, c)
		} else {
			break
		}
	}
	if len(result) == 0 {
		return nil, str
	}
	i, err := strconv.Atoi(string(result))
	if err != nil {
		// TODO should result in error
		return nil, str
	}
	return &i, str[len(result):]
}
