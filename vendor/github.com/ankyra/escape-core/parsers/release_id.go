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
	"regexp"
	"strings"
)

type ReleaseId struct {
	Name    string
	Version string
	Tag     string
}

type QualifiedReleaseId struct {
	*ReleaseId
	Project string
}

func InvalidReleaseFormatError(releaseId string) error {
	return fmt.Errorf("Invalid release format: %s.", releaseId)
}

func InvalidVersionStringInReleaseIdError(releaseId, version string) error {
	return fmt.Errorf("Invalid version string in release ID '%s': '%s' is not a valid version.", releaseId, version)
}

func InvalidReleaseIdError(releaseId, err string) error {
	return fmt.Errorf("Invalid release ID '%s'. %s", releaseId, err)
}

func InvalidVersionError(version string) error {
	return fmt.Errorf("Invalid version string '%s'.", version)
}

func ParseReleaseId(releaseId string) (*ReleaseId, error) {
	colonSplit := strings.Split(releaseId, ":")
	if len(colonSplit) == 1 {
		return parseTagLessReleaseId(releaseId)
	} else if len(colonSplit) != 2 {
		return nil, InvalidReleaseFormatError(releaseId)
	}
	if !IsValidTag(colonSplit[1]) {
		return nil, fmt.Errorf("Invalid tag '%s' in release string '%s'", colonSplit[1], releaseId)
	}
	result := &ReleaseId{}
	result.Name = colonSplit[0]
	result.Tag = colonSplit[1]
	return result, nil
}
func parseTagLessReleaseId(releaseId string) (*ReleaseId, error) {
	split := strings.Split(releaseId, "-")
	if len(split) < 2 { // build-version
		return nil, InvalidReleaseFormatError(releaseId)
	}
	result := &ReleaseId{}
	result.Name = strings.Join(split[:len(split)-1], "-")

	version := split[len(split)-1]
	if version == "latest" || version == "@" || version == "v@" {
		result.Version = "latest"
	} else if strings.HasPrefix(version, "v") {
		result.Version = version[1:]
	} else {
		return nil, InvalidVersionStringInReleaseIdError(releaseId, version)
	}
	if !isValidVersion(result.Version) {
		return nil, InvalidVersionStringInReleaseIdError(releaseId, version)
	}
	return result, nil
}

func ParseQualifiedReleaseId(releaseId string) (*QualifiedReleaseId, error) {
	if releaseId == "" {
		return nil, InvalidReleaseFormatError("''")
	}
	parts := strings.Split(releaseId, "/")
	releaseId = parts[0]
	project := "_"
	if len(parts) > 1 {
		project = parts[0]
		releaseId = strings.Join(parts[1:], "/")
	}
	release, err := ParseReleaseId(releaseId)
	if err != nil {
		return nil, err
	}
	return &QualifiedReleaseId{
		release,
		project,
	}, nil
}

func (r *QualifiedReleaseId) ToString() string {
	return r.Project + "/" + r.ReleaseId.ToString()
}

func isValidVersion(version string) bool {
	if version == "latest" || version == "@" {
		return true
	}
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)*(\.@)?$`)
	return re.Match([]byte(version))
}

func ValidateVersion(version string) error {
	if isValidVersion(version) {
		return nil
	}
	return InvalidVersionError(version)
}

func (r *ReleaseId) ToString() string {
	version := "-" + r.Version
	if version != "-latest" {
		version = "-v" + r.Version
	}
	if r.Tag != "" {
		version = ":" + r.Tag
	}
	return r.Name + version
}

func (r *ReleaseId) NeedsResolving() bool {
	return r.Tag != "" || r.Version == "latest" || strings.HasSuffix(r.Version, ".@")
}
