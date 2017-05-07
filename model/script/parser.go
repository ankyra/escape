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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/parsers"
	"strings"
)

type parseResult struct {
	Result Script
	Rest   string
	Error  error
}

func parseSuccess(script Script, rest string) *parseResult {
	return &parseResult{
		Result: script,
		Rest:   rest,
	}
}
func parseError(err error) *parseResult {
	return &parseResult{
		Error: err,
	}
}

func ParseScript(str string) (Script, error) {
	var result *parseResult
	if strings.HasPrefix(str, "$$") {
		return LiftString(str), nil
	} else if strings.HasPrefix(str, "$") {
		result = parseExpression(str)
	} else if strings.Contains(str, "{{") {
		result = parseExpressionInString(str)
	} else {
		return LiftString(str), nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("Couldn't parse expression '%s': %s", str, result.Error.Error())
	}
	if result.Rest != "" {
		return nil, fmt.Errorf("Invalid expression, unexpected '%s'", result.Rest)
	}
	return result.Result, nil
}

// TODO
// "this should concat {{ $gcp.test }} yall"

func parseExpressionInString(str string) *parseResult {
	return parseSuccess(LiftString(str), "")
}

func parseExpression(str string) *parseResult {
	envLookup := parseEnvLookup(str)
	if envLookup.Error != nil {
		return envLookup
	}
	if envLookup.Rest == "" {
		return envLookup
	}
	return parseApply(envLookup.Result, envLookup.Rest)
}

func parseEnvLookup(str string) *parseResult {
	if !strings.HasPrefix(str, "$") {
		return parseError(fmt.Errorf("Expecting '$'"))
	}
	str = str[1:]
	result, rest := parsers.ParseIdent(str)
	if result == "" {
		return parseError(fmt.Errorf("Expecting indentifier, got '%s'", str))
	}
	envLookup := LiftFunction(builtinEnvLookup)
	key := LiftString(result)
	apply2 := NewApply(envLookup, []Script{LiftString("$")})
	apply1 := NewApply(apply2, []Script{key})
	return parseSuccess(apply1, rest)
}

func parseApply(to Script, str string) *parseResult {
	if !strings.HasPrefix(str, ".") {
		return parseError(fmt.Errorf("Expecting '.', got: '%s'", str))
	}
	str = str[1:]
	result, rest := parsers.ParseIdent(str)
	if result == "" {
		return parseError(fmt.Errorf("Expecting indentifier, got '%s'", str))
	}
	apply := NewApply(to, []Script{LiftString(result)})
	if rest == "" {
		return parseSuccess(apply, rest)
	}
	return parseApply(apply, rest)
}
