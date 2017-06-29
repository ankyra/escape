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
	"github.com/ankyra/escape-core/parsers"
	"strings"
	"unicode"
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
	var result *parseResult
	if str == "" {
		return parseError(fmt.Errorf("Expecting expression starting with '$' or '\"' or '[0-9\\-]', got empty string"))
	} else if strings.HasPrefix(str, "$") {
		result = parseEnvLookup(str)
	} else if strings.HasPrefix(str, "\"") {
		result = parseString(str)
	} else if unicode.IsDigit(rune(str[0])) || str[0:1] == "-" {
		result = parseInteger(str)
	}
	if result == nil {
		return parseError(fmt.Errorf("Expecting expression starting with '$' or '\"' or '[0-9\\-]', got: '%s'", str))
	}
	if result.Error != nil {
		return result
	}
	if strings.HasPrefix(result.Rest, ".") {
		return parseApply(result.Result, result.Rest)
	}
	if strings.HasPrefix(result.Rest, "[") {
		return parseListIndex(result.Result, result.Rest)
	}
	return result
}

func parseInteger(str string) *parseResult {
	if str == "" {
		return parseError(fmt.Errorf("Expecting digit"))
	} else if !unicode.IsDigit(rune(str[0])) && str[0:1] != "-" {
		return parseError(fmt.Errorf("Expecting digit"))
	}
	integer, rest := parsers.ParseInteger(str)
	if integer == nil {
		return parseError(fmt.Errorf("Expecting digit"))
	}
	return parseSuccess(LiftInteger(*integer), rest)
}

func parseString(str string) *parseResult {
	if !strings.HasPrefix(str, `"`) {
		return parseError(fmt.Errorf("Expecting '\"'"))
	}
	str = str[1:]
	result := []byte{}
	escaping := false
	for {
		if str == "" {
			return parseError(fmt.Errorf("Expecting '\"'"))
		}
		if strings.HasPrefix(str, "\"") && !escaping {
			break
		}
		if strings.HasPrefix(str, "\\") {
			if !escaping {
				str = str[1:]
				escaping = true
				continue
			}
		}
		if escaping {
			if str[0] == 'n' {
				result = append(result, '\n')
			} else if str[0] == '"' {
				result = append(result, '"')
			} else if str[0] == 't' {
				result = append(result, '\t')
			} else if str[0] == '\\' {
				result = append(result, '\\')
			} else {
				return parseError(fmt.Errorf("Unexpected escape character '%s' in '%s'", str[0], str))
			}
		} else {
			result = append(result, str[0])
		}
		escaping = false
		str = str[1:]
	}
	return parseSuccess(LiftString(string(result)), str[1:])
}

func parseEnvLookup(str string) *parseResult {
	if !strings.HasPrefix(str, "$") {
		return parseError(fmt.Errorf("Expecting '$'"))
	}
	str = str[1:]
	result, rest := parsers.ParseIdent(str)
	if result == "" {
		if strings.HasPrefix(str, "__") {
			return parseEnvFuncCall(str)
		}
		return parseError(fmt.Errorf("Expecting indentifier, got '%s'", str))
	}
	envLookup := LiftFunction(builtinEnvLookup)
	key := LiftString(result)
	apply2 := NewApply(envLookup, []Script{LiftString("$")})
	apply1 := NewApply(apply2, []Script{key})
	return parseSuccess(apply1, rest)
}

func parseArguments(str string) *parseResult {
	if !strings.HasPrefix(str, "(") {
		return parseError(fmt.Errorf("Expecting '(', got '%s'", str))
	}
	result := []Script{}
	orig := str
	str = strings.TrimSpace(str[1:])

	for {
		if str == "" {
			return parseError(fmt.Errorf("Expecting ')', got EOF in %s", orig))
		}
		if strings.HasPrefix(str, ")") {
			break
		}

		arg := parseExpression(str)
		if arg.Error != nil {
			return parseError(fmt.Errorf("Couldn't parse function argument: %s", arg.Error.Error()))
		}
		result = append(result, arg.Result)

		str = strings.TrimSpace(arg.Rest)
		if strings.HasPrefix(str, ")") {
			break
		}
		if !strings.HasPrefix(str, ",") {
			return parseError(fmt.Errorf("Expecting ',' or ')', but got: \"%s\" in \"%s\"", str, orig))
		}
		str = strings.TrimSpace(str[1:])
	}
	return parseSuccess(LiftList(result), str[1:])
}

func parseEnvFuncCall(str string) *parseResult {
	if !strings.HasPrefix(str, "__") {
		return parseError(fmt.Errorf("Expecting '__', got: '%s'", str))
	}
	funcName, rest := parsers.ParseIdent(str[2:])
	if funcName == "" {
		return parseError(fmt.Errorf("Expecting __indentifier, got '%s'", str))
	}
	if !strings.HasPrefix(rest, "(") {
		return parseError(fmt.Errorf("Expecting '(', got '%s'", rest))
	}
	parseArgsResult := parseArguments(rest)
	if parseArgsResult.Error != nil {
		return parseError(fmt.Errorf("Failed to parse function call to __%s: %s", funcName, parseArgsResult.Error.Error()))
	}
	envLookup := LiftFunction(builtinEnvLookup)
	key := LiftString("__" + funcName)
	apply2 := NewApply(envLookup, []Script{LiftString("$")})
	apply1 := NewApply(apply2, []Script{key})
	apply := NewApply(apply1, ExpectListAtom(parseArgsResult.Result))
	return parseSuccess(apply, parseArgsResult.Rest)
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
	var apply Script
	if strings.HasPrefix(rest, "(") {
		parseArgsResult := parseArguments(rest)
		if parseArgsResult.Error != nil {
			return parseError(fmt.Errorf("Failed to parse function call to __%s: %s", result, parseArgsResult.Error.Error()))
		}
		envLookup := LiftFunction(builtinEnvLookup)
		funcName := LiftString("__" + result)
		apply2 := NewApply(envLookup, []Script{LiftString("$")})
		apply1 := NewApply(apply2, []Script{funcName})
		args := []Script{to}
		for _, arg := range ExpectListAtom(parseArgsResult.Result) {
			args = append(args, arg)
		}
		apply = NewApply(apply1, args)
		rest = parseArgsResult.Rest
	} else {
		apply = NewApply(to, []Script{LiftString(result)})
	}
	if strings.HasPrefix(rest, ".") {
		return parseApply(apply, rest)
	}
	if strings.HasPrefix(rest, "[") {
		return parseListIndex(apply, rest)
	}
	return parseSuccess(apply, rest)
}

func parseListIndex(lst Script, str string) *parseResult {
	if !strings.HasPrefix(str, "[") {
		return parseError(fmt.Errorf("Expecting '[', got: '%s'", str))
	}
	isBeginSlice := false
	str = strings.TrimSpace(str[1:])
	if strings.HasPrefix(str, ":") {
		isBeginSlice = true
		str = str[1:]
	}
	intResult := parseInteger(str)
	if intResult.Error != nil {
		return parseError(fmt.Errorf("Couldn't parse '%s': %s", str, intResult.Error.Error()))
	}
	rest := strings.TrimSpace(intResult.Rest)
	if rest == "" || rest[0:1] != "]" && rest[0:1] != ":" {
		return parseError(fmt.Errorf("Expecting ']' or ':', got: '%s'", rest))
	}
	envLookup := LiftFunction(builtinEnvLookup)
	apply2 := NewApply(envLookup, []Script{LiftString("$")})
	var apply Script

	if rest[0:1] == ":" {
		rest = strings.TrimSpace(rest[1:])
		apply1 := NewApply(apply2, []Script{LiftString("__list_slice")})
		if strings.HasPrefix(rest, "]") {
			apply = NewApply(apply1, []Script{lst, intResult.Result})
		} else {
			endSlice := parseInteger(rest)
			if endSlice.Error != nil {
				return parseError(fmt.Errorf("Couldn't parse '%s': %s", str, endSlice.Error.Error()))
			}
			rest = strings.TrimSpace(endSlice.Rest)
			if !strings.HasPrefix(rest, "]") {
				return parseError(fmt.Errorf("Expecting ']', got: '%s'", rest))
			}
			apply = NewApply(apply1, []Script{lst, intResult.Result, endSlice.Result})
		}
	} else {
		if !isBeginSlice {
			apply1 := NewApply(apply2, []Script{LiftString("__list_index")})
			apply = NewApply(apply1, []Script{lst, intResult.Result})
		} else {
			apply1 := NewApply(apply2, []Script{LiftString("__list_slice")})
			apply = NewApply(apply1, []Script{lst, LiftInteger(0), intResult.Result})
		}
	}
	rest = rest[1:]
	if strings.HasPrefix(rest, ".") {
		return parseApply(apply, rest)
	}
	return parseSuccess(apply, rest)
}
