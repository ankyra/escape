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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/util"
	"github.com/ankyra/escape-core/variables"
)

type ExecStage struct {
	Script string `json:"script"`
}

type ProviderConfig struct {
	Name string `json:"name"`
}

func NewProviderConfig(name string) *ProviderConfig {
	return &ProviderConfig{name}
}

type ReleaseMetadata struct {
	ApiVersion             int               `json:"api_version"`
	BuiltWithCoreVersion   string            `json:"built_with_core_version"`
	BuiltWithEscapeVersion string            `json:"built_with_escape_version"`
	Description            string            `json:"description"`
	Files                  map[string]string `json:"files", {}`
	Logo                   string            `json:"logo"`
	Name                   string            `json:"name"`
	Repository             string            `json:"repository"`
	Branch                 string            `json:"branch"`
	Revision               string            `json:"git_revision"`
	RevisionMessage        string            `json:"revision_message"`
	RevisionAuthor         string            `json:"revision_author"`
	Metadata               map[string]string `json:"metadata"`
	Version                string            `json:"version"`

	Consumes  []*ConsumerConfig     `json:"consumes"`
	Downloads []*DownloadConfig     `json:"downloads"`
	Depends   []*DependencyConfig   `json:"depends"`
	Errands   map[string]*Errand    `json:"errands"`
	Extends   []*ExtensionConfig    `json:"extends"`
	Inputs    []*variables.Variable `json:"inputs"`
	Outputs   []*variables.Variable `json:"outputs"`
	Project   string                `json:"project"`
	Provides  []*ProviderConfig     `json:"provides"`
	Stages    map[string]*ExecStage `json:"stages"`
	Templates []*templates.Template `json:"templates"`

	// The VariableCtx is deprecrated and used to support packages that were
	// compiled using an Escape version below v0.23.0. Pay no heed.  This field
	// was superseded by tracking VariableNames on the DependencyConfig and
	// ConsumerConfig instead.
	VariableCtx map[string]string `json:"variable_context"`
}

func NewEmptyReleaseMetadata() *ReleaseMetadata {
	return &ReleaseMetadata{
		ApiVersion:           CurrentApiVersion,
		BuiltWithCoreVersion: CoreVersion,
		Files:                map[string]string{},
		Metadata:             map[string]string{},

		Consumes:    []*ConsumerConfig{},
		Depends:     []*DependencyConfig{},
		Errands:     map[string]*Errand{},
		Extends:     []*ExtensionConfig{},
		Inputs:      []*variables.Variable{},
		Outputs:     []*variables.Variable{},
		Provides:    []*ProviderConfig{},
		Stages:      map[string]*ExecStage{},
		Templates:   []*templates.Template{},
		VariableCtx: map[string]string{},
	}
}

func NewReleaseMetadata(name, version string) *ReleaseMetadata {
	m := NewEmptyReleaseMetadata()
	m.Name = name
	m.Version = version
	m.Project = "_"
	return m
}

func NewReleaseMetadataFromJsonString(content string) (*ReleaseMetadata, error) {
	result := NewEmptyReleaseMetadata()
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal JSON release metadata: %s", err.Error())
	}
	return result, validate(result)
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

func (m *ReleaseMetadata) Validate() error {
	return validate(m)
}

func validate(m *ReleaseMetadata) error {
	if m == nil {
		return fmt.Errorf("Missing release metadata")
	}
	if m.Name == "" {
		return fmt.Errorf("Missing name field in release metadata")
	}
	if err := validateName(m.Name); err != nil {
		return err
	}
	if m.Version == "" {
		return fmt.Errorf("Missing version field in release metadata")
	}
	if m.Project == "" {
		m.Project = "_"
	}
	if err := ValidateProjectName(m.Project); err != nil {
		return err
	}
	if m.ApiVersion <= 0 || m.ApiVersion > CurrentApiVersion {
		return fmt.Errorf("The release metadata is compiled with a version of Escape targetting API version v%d, but this build supports up to v%d", m.ApiVersion, CurrentApiVersion)
	}
	if err := parsers.ValidateVersion(m.Version); err != nil {
		return err
	}
	for _, i := range m.Inputs {
		if err := i.Validate(); err != nil {
			return err
		}
	}
	for _, i := range m.Outputs {
		if err := i.Validate(); err != nil {
			return err
		}
	}
	for _, d := range m.Depends {
		if err := d.Validate(m); err != nil {
			return err
		}
	}
	for _, c := range m.Consumes {
		if err := c.ValidateAndFix(); err != nil {
			return err
		}
	}
	for _, d := range m.Downloads {
		if err := d.ValidateAndFix(); err != nil {
			return err
		}
	}
	return nil
}

func validateName(name string) error {
	re := regexp.MustCompile("^[a-z]+[a-z0-9-_]+$")
	if !re.MatchString(name) {
		return fmt.Errorf("Invalid name '%s'", name)
	}
	protectedNames := map[string]bool{
		"this":    false,
		"string":  false,
		"integer": false,
		"list":    false,
		"dict":    false,
		"func":    false,
	}
	if _, found := protectedNames[name]; found {
		return fmt.Errorf("The name '%s' is a protected variable.", name)
	}
	return nil
}

func ValidateProjectName(name string) error {
	if name == "_" {
		return nil
	}
	return validateName(name)
}

func (m *ReleaseMetadata) AddExtension(releaseId string) {
	for _, e := range m.Extends {
		if e.ReleaseId == releaseId {
			return
		}
	}
	m.Extends = append(m.Extends, NewExtensionConfig(releaseId))
}

func (m *ReleaseMetadata) GetExtensions() []string {
	result := []string{}
	for _, ext := range m.Extends {
		result = append(result, ext.ReleaseId)
	}
	return result
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

func (m *ReleaseMetadata) AddInputVariable(input *variables.Variable) {
	for _, i := range m.Inputs {
		if i.Id == input.Id {
			i.Default = input.Default
			if len(i.Scopes) < len(input.Scopes) {
				i.Scopes = input.Scopes
			}
			return
		}
	}
	m.Inputs = append(m.Inputs, input)
}

func (m *ReleaseMetadata) AddOutputVariable(output *variables.Variable) {
	for _, i := range m.Outputs {
		if i.Id == output.Id {
			if len(i.Scopes) < len(output.Scopes) {
				i.Scopes = output.Scopes
			}
			return
		}
	}
	m.Outputs = append(m.Outputs, output)
}

func (m *ReleaseMetadata) AddConsumes(c *ConsumerConfig) {
	for _, consumer := range m.Consumes {
		if consumer.Name == c.Name && consumer.VariableName == c.VariableName {
			if len(consumer.Scopes) < len(c.Scopes) {
				consumer.Scopes = c.Scopes
			}
			return
		}
	}
	m.Consumes = append(m.Consumes, c)
}

func (m *ReleaseMetadata) SetConsumes(c []string) {
	for _, consumer := range c {
		m.AddConsumes(NewConsumerConfig(consumer))
	}
}

func (m *ReleaseMetadata) GetConsumes(stage string) []string {
	result := []string{}
	for _, c := range m.GetConsumerConfig(stage) {
		result = append(result, c.Name)
	}
	return result
}

func (m *ReleaseMetadata) GetConsumerConfig(stage string) []*ConsumerConfig {
	result := []*ConsumerConfig{}
	for _, c := range m.Consumes {
		if c.InScope(stage) {
			result = append(result, c)
		}
	}
	return result
}

func (m *ReleaseMetadata) GetDownloads(stage string) []*DownloadConfig {
	result := []*DownloadConfig{}
	for _, d := range m.Downloads {
		if d.InScope(stage) {
			result = append(result, d)
		}
	}
	return result
}

func (m *ReleaseMetadata) GetErrands() map[string]*Errand {
	result := map[string]*Errand{}
	for key, val := range m.Errands {
		result[key] = val
	}
	return result
}

func (m *ReleaseMetadata) GetInputs(stage string) []*variables.Variable {
	result := []*variables.Variable{}
	for _, i := range m.Inputs {
		if i.InScope(stage) {
			result = append(result, i)
		}
	}
	return result
}

func (m *ReleaseMetadata) GetOutputs(stage string) []*variables.Variable {
	result := []*variables.Variable{}
	for _, i := range m.Outputs {
		if i.InScope(stage) {
			result = append(result, i)
		}
	}
	return result
}

func (m *ReleaseMetadata) GetTemplates(stage string) []*templates.Template {
	result := []*templates.Template{}
	for _, t := range m.Templates {
		if t.InScope(stage) {
			result = append(result, t)
		}
	}
	return result
}

func (m *ReleaseMetadata) AddProvides(p string) {
	for _, provider := range m.Provides {
		if provider.Name == p {
			return
		}
	}
	m.Provides = append(m.Provides, NewProviderConfig(p))
}

func (m *ReleaseMetadata) GetProvides() []string {
	result := []string{}
	for _, c := range m.Provides {
		result = append(result, c.Name)
	}
	return result
}

func (m *ReleaseMetadata) SetProvides(p []string) {
	for _, provider := range p {
		m.AddProvides(provider)
	}
}

func (m *ReleaseMetadata) AddDependency(dep *DependencyConfig) {
	m.Depends = append(m.Depends, dep)
}

func (m *ReleaseMetadata) AddDependencyFromString(dep string) {
	cfg := NewDependencyConfig(dep)
	if err := cfg.Validate(m); err != nil {
		panic(err)
	}
	m.Depends = append(m.Depends, cfg)
}

func (m *ReleaseMetadata) SetDependencies(deps []string) {
	m.Depends = []*DependencyConfig{}
	for _, d := range deps {
		m.AddDependencyFromString(d)
	}
}

func (m *ReleaseMetadata) GetReleaseId() string {
	return m.Name + "-v" + m.Version
}

func (m *ReleaseMetadata) GetQualifiedReleaseId() string {
	return m.GetProject() + "/" + m.Name + "-v" + m.Version
}

func (m *ReleaseMetadata) GetProject() string {
	if m.Project == "" {
		return "_"
	}
	return m.Project
}

func (m *ReleaseMetadata) GetVersionlessReleaseId() string {
	return m.GetProject() + "/" + m.Name
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
	for key, val := range m.Metadata {
		metadataDict[key] = script.LiftString(val)
	}
	return map[string]script.Script{
		"metadata": script.LiftDict(metadataDict),

		"branch":              script.LiftString(m.Branch),
		"description":         script.LiftString(m.Description),
		"logo":                script.LiftString(m.Logo),
		"name":                script.LiftString(m.Name),
		"revision":            script.LiftString(m.Revision),
		"repository":          script.LiftString(m.Repository),
		"version":             script.LiftString(m.Version),
		"release":             script.LiftString(m.GetReleaseId()),
		"versionless_release": script.LiftString(m.GetVersionlessReleaseId()),
		"id": script.LiftString(m.GetQualifiedReleaseId()),
	}
}
