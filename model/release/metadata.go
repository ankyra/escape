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

package release

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type execStage struct {
	Script string `json:"script"`
}

type releaseMetadata struct {
	ApiVersion  string                `json:"api_version"`
	Branch      string                `json:"branch"`
	Consumes    []string              `json:"consumes"`
	Depends     []string              `json:"depends"`
	Description string                `json:"description"`
	Errands     map[string]*errand    `json:"errands"`
	Files       map[string]string     `json:"files", {}`
	Revision    string                `json:"git_revision"`
	Inputs      []*variable           `json:"inputs"`
	Logo        string                `json:"logo"`
	Metadata    map[string]string     `json:"metadata"`
	Name        string                `json:"name"`
	Outputs     []*variable           `json:"outputs"`
	Path        string                `json:"path"`
	Provides    []string              `json:"provides"`
	Test        string                `json:"test"`
	Type        string                `json:"type"`
	VariableCtx map[string]string     `json:"variable_context"`
	Version     string                `json:"version"`
	Stages      map[string]*execStage `json:"stages"`
}

func NewEmptyReleaseMetadata() ReleaseMetadata {
	return &releaseMetadata{
		ApiVersion:  "1",
		Consumes:    []string{},
		Provides:    []string{},
		Depends:     []string{},
		Files:       map[string]string{},
		Metadata:    map[string]string{},
		Errands:     map[string]*errand{},
		Stages:      map[string]*execStage{},
		Inputs:      []*variable{},
		Outputs:     []*variable{},
		VariableCtx: map[string]string{},
	}
}

func NewReleaseMetadata(typ, name, version string) ReleaseMetadata {
	m := NewEmptyReleaseMetadata()
	m.(*releaseMetadata).Type = typ
	m.(*releaseMetadata).Name = name
	m.(*releaseMetadata).Version = version
	return m
}

func NewReleaseMetadataFromJsonString(content string) (ReleaseMetadata, error) {
	result := NewEmptyReleaseMetadata()
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal JSON release metadata: %s", err.Error())
	}
	if err := validate(result.(*releaseMetadata)); err != nil {
		return nil, err
	}
	return result, nil
}

func NewReleaseMetadataFromFile(metadataFile string) (ReleaseMetadata, error) {
	if !util.PathExists(metadataFile) {
		return nil, errors.New("Release metadata file " + metadataFile + " does not exist")
	}
	content, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}
	return NewReleaseMetadataFromJsonString(string(content))
}

func validate(m *releaseMetadata) error {
	if m.Type == "" {
		return fmt.Errorf("Missing type field in release metadata")
	}
	if m.Name == "" {
		return fmt.Errorf("Missing name field in release metadata")
	}
	if m.Version == "" {
		return fmt.Errorf("Missing version field in release metadata")
	}
	return nil
}

func (m *releaseMetadata) GetStage(stage string) *execStage {
	result, ok := m.Stages[stage]
	if !ok {
		result = &execStage{}
		m.Stages[stage] = result
	}
	return result
}

func (m *releaseMetadata) SetStage(stage, script string) {
	st := m.GetStage(stage)
	st.Script = script
}
func (m *releaseMetadata) GetScript(stage string) string {
	return m.GetStage(stage).Script
}
func (m *releaseMetadata) GetApiVersion() string {
	return m.ApiVersion
}
func (m *releaseMetadata) GetBranch() string {
	return m.Branch
}
func (m *releaseMetadata) SetConsumes(c []string) {
	m.Consumes = c
}
func (m *releaseMetadata) GetConsumes() []string {
	return m.Consumes
}
func (m *releaseMetadata) GetDescription() string {
	return m.Description
}
func (m *releaseMetadata) GetErrands() map[string]Errand {
	result := map[string]Errand{}
	for key, val := range m.Errands {
		result[key] = val
	}
	return result
}
func (m *releaseMetadata) GetFiles() map[string]string {
	return m.Files
}
func (m *releaseMetadata) GetInputs() []Variable {
	result := []Variable{}
	for _, i := range m.Inputs {
		result = append(result, i)
	}
	return result
}
func (m *releaseMetadata) GetRevision() string {
	return m.Revision
}
func (m *releaseMetadata) GetLogo() string {
	return m.Logo
}
func (m *releaseMetadata) GetMetadata() map[string]string {
	return m.Metadata
}
func (m *releaseMetadata) GetName() string {
	return m.Name
}
func (m *releaseMetadata) GetOutputs() []Variable {
	result := []Variable{}
	for _, i := range m.Outputs {
		result = append(result, i)
	}
	return result
}
func (m *releaseMetadata) GetPath() string {
	return m.Path
}
func (m *releaseMetadata) GetProvides() []string {
	return m.Provides
}
func (m *releaseMetadata) GetType() string {
	return m.Type
}
func (m *releaseMetadata) GetVersion() string {
	return m.Version
}
func (m *releaseMetadata) GetDependencies() []string {
	return m.Depends
}
func (m *releaseMetadata) GetVariableContext() map[string]string {
	if m.VariableCtx == nil {
		return map[string]string{}
	}
	return m.VariableCtx
}
func (m *releaseMetadata) SetVariableInContext(v string, ref string) {
	ctx := m.GetVariableContext()
	ctx[v] = ref
	m.VariableCtx = ctx
}
func (m *releaseMetadata) GetReleaseId() string {
	return m.Type + "-" + m.Name + "-v" + m.Version
}

func (m *releaseMetadata) GetVersionlessReleaseId() string {
	return m.Type + "-" + m.Name
}

func (m *releaseMetadata) AddInputVariable(input Variable) {
	m.Inputs = append(m.Inputs, input.(*variable))
}
func (m *releaseMetadata) AddOutputVariable(output Variable) {
	m.Outputs = append(m.Outputs, output.(*variable))
}

func (m *releaseMetadata) ToJson() string {
	str, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (m *releaseMetadata) ToDict() (map[string]interface{}, error) {
	asJson := []byte(m.ToJson())
	result := map[string]interface{}{}
	if err := json.Unmarshal(asJson, &result); err != nil {
		return nil, fmt.Errorf("Couldn't marshal release metadata: %s. This is a bug in Escape", err.Error())
	}
	return result, nil
}

func (m *releaseMetadata) WriteJsonFile(path string) error {
	contents := []byte(m.ToJson())
	return ioutil.WriteFile(path, contents, 0644)
}

func (m *releaseMetadata) AddFileWithDigest(path, hexDigest string) {
	m.Files[path] = hexDigest
}

func (m *releaseMetadata) ToDependency() Dependency {
	return NewDependencyFromMetadata(m)
}

func (m *releaseMetadata) GetDirectories() []string {
	dirs := map[string]bool{}
	for file := range m.Files {
		dir, _ := filepath.Split(file)
		dirs[dir] = true
		root := ""
		for _, d := range strings.Split(dir, "/") {
			if d != "" {
				root += d + "/"
				dirs[root] = true
			}
		}
	}
	result := []string{}
	for d := range dirs {
		if d != "" {
			result = append(result, d)
		}
	}
	return result
}
