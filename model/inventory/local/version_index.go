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

package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/util"
)

type VersionIndex struct {
	Path              string                           `json:"-"`
	Name              string                           `json:"name"`
	EscapeVersion     string                           `json:"escape_version"`
	EscapeCoreVersion string                           `json:"escape_core_version"`
	CoreAPIVersion    int                              `json:"core_api_version"`
	Versions          map[string]*core.ReleaseMetadata `json:"versions"`
}

func NewVersionIndex() *VersionIndex {
	return &VersionIndex{
		EscapeVersion:     util.EscapeVersion,
		EscapeCoreVersion: core.CoreVersion,
		CoreAPIVersion:    core.CurrentApiVersion,
		Versions:          map[string]*core.ReleaseMetadata{},
	}
}

func LoadVersionIndexFromFile(path string) (*VersionIndex, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	index := NewVersionIndex()
	if err := json.Unmarshal(content, index); err != nil {
		return nil, fmt.Errorf("Could not read local version index at '%s': %s", path, err.Error())
	}
	if index.CoreAPIVersion > core.CurrentApiVersion {
		return nil, fmt.Errorf("Could not load application version index, because the index was created by a version of Escape that supports a different (newer) API v%d, instead of v%d, which this version of Escape supports.", index.CoreAPIVersion, core.CurrentApiVersion)
	}
	index.Path = path
	index.EscapeVersion = util.EscapeVersion
	index.EscapeCoreVersion = core.CoreVersion
	index.CoreAPIVersion = core.CurrentApiVersion
	return index, nil
}

func LoadVersionIndexFromFileOrCreateNew(name, path string) (*VersionIndex, error) {
	if util.PathExists(path) {
		return LoadVersionIndexFromFile(path)
	}
	index := NewVersionIndex()
	index.Name = name
	index.Path = path
	return index, index.Save()
}

func (v *VersionIndex) Save() error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(v.Path, bytes, 0644)
}

func (v *VersionIndex) AddRelease(m *core.ReleaseMetadata) error {
	_, exists := v.Versions[m.Version]
	if exists {
		return fmt.Errorf("Could not add release to local inventory at %s. Version %s already exists.", v.Path, m.Version)
	}
	v.Versions[m.Version] = m
	return nil
}

func (v *VersionIndex) GetVersions() []string {
	versions := []string{}
	for version := range v.Versions {
		versions = append(versions, version)
	}
	return versions
}
