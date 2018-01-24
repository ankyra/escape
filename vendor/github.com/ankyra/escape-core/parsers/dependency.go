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
	"fmt"
	"strings"
)

type ParsedDependency struct {
	QualifiedReleaseId
	VariableName string
}

const expectFormatError = `Expecting either a "[RELEASE_ID]" or "[RELEASE_ID] as [VARIABLE]".`

func MalformedDependencyStringError(depString, err string) error {
	return fmt.Errorf("Malformed dependency string '%s'. %s", depString, err)
}

func MalformedDependencyStringExpectingError(depString string) error {
	return MalformedDependencyStringError(depString, expectFormatError)
}

func ExpectingAsError(unexpected, in string) error {
	return fmt.Errorf("Unexpected '%s'; expecting 'as' in '%s. %s'", unexpected, in, expectFormatError)
}

func ParseDependency(str string) (*ParsedDependency, error) {
	result := &ParsedDependency{}
	split := strings.Split(str, " ")
	parts := []string{}
	for _, part := range split {
		if strings.TrimSpace(part) != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) != 1 && len(parts) != 3 {
		return nil, MalformedDependencyStringExpectingError(str)
	}
	releaseId, err := ParseQualifiedReleaseId(parts[0])
	if err != nil {
		return nil, err
	}
	if len(parts) == 3 {
		if parts[1] != "as" {
			return nil, ExpectingAsError(parts[1], str)
		}
		id, err := ParseVariableIdent(parts[2])
		if err != nil {
			return nil, MalformedDependencyStringError(str, err.Error())
		}
		result.VariableName = id
	}
	result.QualifiedReleaseId = *releaseId
	return result, nil
}
