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
	"fmt"
	"regexp"
	"strings"
)

type ReleaseId struct {
	Name    string
	Version string
}

func ParseReleaseId(releaseId string) (*ReleaseId, error) {
	split := strings.Split(releaseId, "-")
	if len(split) < 2 { // build-version
		return nil, fmt.Errorf("Invalid release format: %s", releaseId)
	}
	result := &ReleaseId{}
	result.Name = strings.Join(split[:len(split)-1], "-")

	version := split[len(split)-1]
	if version == "latest" || version == "@" || version == "v@" {
		result.Version = "latest"
	} else if strings.HasPrefix(version, "v") {
		result.Version = version[1:]
	} else {
		return nil, fmt.Errorf("Invalid version string in release ID '%s': %s", releaseId, version)
	}

	if err := result.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid release ID '%s': %s", releaseId, err.Error())
	}
	return result, nil
}

func (r *ReleaseId) Validate() error {
	return ValidateVersion(r.Version)
}

func ValidateVersion(version string) error {
	if version == "latest" || version == "@" {
		return nil
	}
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)*(\.@)?$`)
	matches := re.Match([]byte(version))
	if !matches {
		return fmt.Errorf("Invalid version format: %s", version)
	}
	return nil
}

func (r *ReleaseId) ToString() string {
	version := r.Version
	if version != "latest" {
		version = "v" + version
	}
	return r.Name + "-" + version
}

func (r *ReleaseId) NeedsResolving() bool {
	return r.Version == "latest" || strings.HasSuffix(r.Version, ".@")
}
