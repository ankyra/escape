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

package core

import (
	"github.com/ankyra/escape-core/parsers"
)

type Dependency struct {
	Name         string
	Version      string
	VariableName string
}

func NewDependencyFromMetadata(metadata *ReleaseMetadata) *Dependency {
	return &Dependency{
		Name:    metadata.GetName(),
		Version: metadata.GetVersion(),
	}
}

func NewDependencyFromString(str string) (*Dependency, error) {
	parsed, err := parsers.ParseDependency(str)
	if err != nil {
		return nil, err
	}
	return &Dependency{
		Name:         parsed.Name,
		Version:      parsed.Version,
		VariableName: parsed.VariableName,
	}, nil
}

func (d *Dependency) GetBuild() string {
	return d.Name
}
func (d *Dependency) GetVariableName() string {
	return d.VariableName
}
func (d *Dependency) GetVersion() string {
	return d.Version
}

func (d *Dependency) GetReleaseId() string {
	version := "v" + d.Version
	if d.Version == "latest" {
		version = d.Version
	}
	return d.Name + "-" + version
}
func (d *Dependency) GetVersionlessReleaseId() string {
	return d.Name
}
