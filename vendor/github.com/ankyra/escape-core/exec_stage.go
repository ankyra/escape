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

package core

import (
	"fmt"
	"strings"
)

type ExecStage struct {

	// The command to run. Its arguments, if any, should be defined using the
	// "args" field.
	Cmd string `json:"cmd,omitempty"`
	// Arguments to the command.
	Args []string `json:"args,omitempty"`

	// An inline script, which will be executed using bash. It's an error to
	// specify both the "cmd" and "inline" fields.
	Inline string `json:"inline,omitempty"`

	// Relative path to a script. If the "cmd" field is already populated
	// then this field will be ignored entirely.
	RelativeScript string `json:"script,omitempty"`
}

func ExpectingTypeForExecStageError(typ, field string, val interface{}) error {
	return fmt.Errorf("Expecting %s for exec stage field %s; got '%T'", typ, field, val)
}

func NewExecStageFromDict(values map[interface{}]interface{}) (*ExecStage, error) {
	result := ExecStage{}
	for k, val := range values {
		kStr, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string key in exec stage. Got '%T'", k)
		}
		if kStr == "cmd" {
			valString, ok := val.(string)
			if !ok {
				return nil, ExpectingTypeForExecStageError("string", kStr, val)
			}
			result.Cmd = valString
		} else if kStr == "args" {
			valList, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("Expecting []string in args exec stage, got '%v' (%T)", val, val)
			}
			args := []string{}
			for _, val := range valList {
				kStr, ok := val.(string)
				if !ok {
					return nil, fmt.Errorf("Expecting string in exec stage args, got '%v' (%T)", val, val)
				}
				args = append(args, kStr)
			}
			result.Args = args
		} else if kStr == "inline" {
			valString, ok := val.(string)
			if !ok {
				return nil, ExpectingTypeForExecStageError("string", kStr, val)
			}
			result.Inline = valString
		} else if kStr == "script" {
			valString, ok := val.(string)
			if !ok {
				return nil, ExpectingTypeForExecStageError("string", kStr, val)
			}
			result.RelativeScript = valString
		}
	}
	return &result, nil
}

func NewExecStageForRelativeScript(script string) *ExecStage {
	return &ExecStage{
		RelativeScript: script,
	}
}

func (e *ExecStage) IsEmpty() bool {
	return e.Cmd == "" && e.RelativeScript == "" && e.Inline == ""
}

func (e *ExecStage) GetAsCommand() []string {
	if e.Cmd != "" {
		result := []string{e.Cmd}
		return append(result, e.Args...)
	} else if e.RelativeScript != "" {
		script := "./" + e.RelativeScript + " .escape/outputs.json"
		return []string{"/bin/sh", "-c", script}
	} else if e.Inline != "" {
		panic("not yet supported")
	}
	return []string{}
}

func (e *ExecStage) ValidateAndFix() error {
	fieldsSet := 0
	if e.Cmd != "" {
		fieldsSet += 1
	}
	if e.Inline != "" {
		fieldsSet += 1
	}
	if e.RelativeScript != "" {
		fieldsSet += 1
	}
	if fieldsSet > 1 {
		return fmt.Errorf("More than one field is set. Please specify only one of script, cmd or inline.")
	}
	return nil
}

func (e *ExecStage) String() string {
	if e.Cmd != "" {
		return fmt.Sprintf("%s %s", e.Cmd, strings.Join(e.Args, " "))
	} else if e.RelativeScript != "" {
		return e.RelativeScript
	} else {
		return e.Inline
	}
}
