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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ankyra/escape-client/util"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/variables"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ExecStage struct {
	Script string `json:"script"`
}

type ReleaseMetadata struct {
	ApiVersion  string                `json:"api_version"`
	Branch      string                `json:"branch"`
	Consumes    []string              `json:"consumes"`
	Depends     []string              `json:"depends"`
	Extends     []string              `json:"extends"`
	Description string                `json:"description"`
	Errands     map[string]*Errand    `json:"errands"`
	Files       map[string]string     `json:"files", {}`
	Revision    string                `json:"git_revision"`
	Inputs      []*variables.Variable `json:"inputs"`
	Logo        string                `json:"logo"`
	Metadata    map[string]string     `json:"metadata"`
	Name        string                `json:"name"`
	Outputs     []*variables.Variable `json:"outputs"`
	Path        string                `json:"path"`
	Provides    []string              `json:"provides"`
	Templates   []*templates.Template `json:"templates"`
	Test        string                `json:"test"`
	VariableCtx map[string]string     `json:"variable_context"`
	Version     string                `json:"version"`
	Stages      map[string]*ExecStage `json:"stages"`
}

func NewEmptyReleaseMetadata() *ReleaseMetadata {
	return &ReleaseMetadata{
		ApiVersion:  "2",
		Consumes:    []string{},
		Provides:    []string{},
		Depends:     []string{},
		Extends:     []string{},
		Files:       map[string]string{},
		Metadata:    map[string]string{},
		Errands:     map[string]*Errand{},
		Stages:      map[string]*ExecStage{},
		Inputs:      []*variables.Variable{},
		Outputs:     []*variables.Variable{},
		Templates:   []*templates.Template{},
		VariableCtx: map[string]string{},
	}
}

func NewReleaseMetadata(name, version string) *ReleaseMetadata {
	m := NewEmptyReleaseMetadata()
	m.Name = name
	m.Version = version
	return m
}

func NewReleaseMetadataFromJsonString(content string) (*ReleaseMetadata, error) {
	result := NewEmptyReleaseMetadata()
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal JSON release metadata: %s", err.Error())
	}
	if err := validate(result); err != nil {
		return nil, err
	}
	return result, nil
}

func NewReleaseMetadataFromFile(metadataFile string) (*ReleaseMetadata, error) {
	if !util.PathExists(metadataFile) {
		return nil, errors.New("Release metadata file " + metadataFile + " does not exist")
	}
	content, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}
	return NewReleaseMetadataFromJsonString(string(content))
}

func validate(m *ReleaseMetadata) error {
	if m.Name == "" {
		return fmt.Errorf("Missing name field in release metadata")
	}
	if m.Version == "" {
		return fmt.Errorf("Missing version field in release metadata")
	}
	if err := parsers.ValidateVersion(m.Version); err != nil {
		return err
	}
	return nil
}
func (m *ReleaseMetadata) GetExtends() []string {
	return m.Extends
}
func (m *ReleaseMetadata) GetStages() map[string]*ExecStage {
	return m.Stages
}
func (m *ReleaseMetadata) GetStage(stage string) *ExecStage {
	result, ok := m.Stages[stage]
	if !ok {
		result = &ExecStage{}
		m.Stages[stage] = result
	}
	return result
}

func (m *ReleaseMetadata) SetStage(stage, script string) {
	if script == "" {
		return
	}
	st := m.GetStage(stage)
	st.Script = script
}
func (m *ReleaseMetadata) GetScript(stage string) string {
	return m.GetStage(stage).Script
}
func (m *ReleaseMetadata) GetApiVersion() string {
	return m.ApiVersion
}
func (m *ReleaseMetadata) GetBranch() string {
	return m.Branch
}
func (m *ReleaseMetadata) SetConsumes(c []string) {
	m.Consumes = c
}
func (m *ReleaseMetadata) GetConsumes() []string {
	return m.Consumes
}
func (m *ReleaseMetadata) GetDescription() string {
	return m.Description
}
func (m *ReleaseMetadata) GetErrands() map[string]*Errand {
	result := map[string]*Errand{}
	for key, val := range m.Errands {
		result[key] = val
	}
	return result
}
func (m *ReleaseMetadata) GetFiles() map[string]string {
	return m.Files
}
func (m *ReleaseMetadata) GetInputs() []*variables.Variable {
	result := []*variables.Variable{}
	for _, i := range m.Inputs {
		result = append(result, i)
	}
	return result
}
func (m *ReleaseMetadata) GetTemplates() []*templates.Template {
	return m.Templates
}
func (m *ReleaseMetadata) GetRevision() string {
	return m.Revision
}
func (m *ReleaseMetadata) GetLogo() string {
	return m.Logo
}
func (m *ReleaseMetadata) GetMetadata() map[string]string {
	return m.Metadata
}
func (m *ReleaseMetadata) GetName() string {
	return m.Name
}
func (m *ReleaseMetadata) GetOutputs() []*variables.Variable {
	result := []*variables.Variable{}
	for _, i := range m.Outputs {
		result = append(result, i)
	}
	return result
}
func (m *ReleaseMetadata) GetPath() string {
	return m.Path
}
func (m *ReleaseMetadata) GetProvides() []string {
	return m.Provides
}
func (m *ReleaseMetadata) GetVersion() string {
	return m.Version
}
func (m *ReleaseMetadata) GetDependencies() []string {
	return m.Depends
}
func (m *ReleaseMetadata) GetVariableContext() map[string]string {
	if m.VariableCtx == nil {
		return map[string]string{}
	}
	return m.VariableCtx
}
func (m *ReleaseMetadata) SetVariableInContext(v string, ref string) {
	ctx := m.GetVariableContext()
	ctx[v] = ref
	m.VariableCtx = ctx
}
func (m *ReleaseMetadata) GetReleaseId() string {
	return m.Name + "-v" + m.Version
}

func (m *ReleaseMetadata) GetVersionlessReleaseId() string {
	return m.Name
}

func (m *ReleaseMetadata) AddInputVariable(input *variables.Variable) {
	m.Inputs = append(m.Inputs, input)
}
func (m *ReleaseMetadata) AddOutputVariable(output *variables.Variable) {
	m.Outputs = append(m.Outputs, output)
}

func (m *ReleaseMetadata) ToJson() string {
	str, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (m *ReleaseMetadata) ToDict() (map[string]interface{}, error) {
	asJson := []byte(m.ToJson())
	result := map[string]interface{}{}
	if err := json.Unmarshal(asJson, &result); err != nil {
		return nil, fmt.Errorf("Couldn't marshal release metadata: %s. This is a bug in Escape", err.Error())
	}
	return result, nil
}

func (m *ReleaseMetadata) WriteJsonFile(path string) error {
	contents := []byte(m.ToJson())
	return ioutil.WriteFile(path, contents, 0644)
}

func (m *ReleaseMetadata) AddFileWithDigest(path, hexDigest string) {
	m.Files[path] = hexDigest
}

func (m *ReleaseMetadata) ToDependency() *Dependency {
	return NewDependencyFromMetadata(m)
}

func (m *ReleaseMetadata) GetDirectories() []string {
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

func (m *ReleaseMetadata) ToScript() script.Script {
	return script.LiftDict(m.ToScriptMap())
}

func (m *ReleaseMetadata) ToScriptMap() map[string]script.Script {
	metadataDict := map[string]script.Script{}
	for key, val := range m.GetMetadata() {
		metadataDict[key] = script.LiftString(val)
	}
	return map[string]script.Script{
		"metadata": script.LiftDict(metadataDict),

		"branch":      script.LiftString(m.GetBranch()),
		"description": script.LiftString(m.GetDescription()),
		"logo":        script.LiftString(m.GetLogo()),
		"build":       script.LiftString(m.GetName()),
		"revision":    script.LiftString(m.GetRevision()),
		"id":          script.LiftString(m.GetReleaseId()),
		"version":     script.LiftString(m.GetVersion()),
	}
}
