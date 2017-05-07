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
	"strings"
)

type ReleaseId struct {
	Type    string
	Build   string
	Version string
}

func ParseReleaseId(releaseId string) (*ReleaseId, error) {
	split := strings.Split(releaseId, "-")
	if len(split) < 3 { // type-build-version
		return nil, errors.New("Invalid release format: " + releaseId)
	}
	result := &ReleaseId{}
	result.Type = split[0]
	result.Build = strings.Join(split[1:len(split)-1], "-")

	version := split[len(split)-1]
	if version == "latest" || version == "@" || version == "v@" {
		result.Version = "latest"
	} else if strings.HasPrefix(version, "v") {
		result.Version = version[1:]
	} else {
		return nil, errors.New("Invalid version string in release ID '" + releaseId + "': " + version)
	}
	return result, nil
}
