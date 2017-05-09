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

package escape_plan

import (
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func indent(s string) string {
	parts := []string{}
	for _, part := range strings.Split(s, "\n") {
		if part != "" {
			parts = append(parts, "  "+part)
		}
	}
	return strings.Join(parts, "\n")
}

type escapePlan struct {
	Build       string                 `yaml:"build"`
	Consumes    []string               `yaml:"consumes,omitempty"`
	Depends     []string               `yaml:"depends,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Errands     map[string]interface{} `yaml:"errands,omitempty"`
	Includes    []string               `yaml:"includes,omitempty"`
	Inputs      []interface{}          `yaml:"inputs,omitempty"`
	Logo        string                 `yaml:"logo,omitempty"`
	Metadata    map[string]string      `yaml:"metadata,omitempty"`
	Outputs     []interface{}          `yaml:"outputs,omitempty"`
	Path        string                 `yaml:"path,omitempty"`
	PostBuild   string                 `yaml:"post_build,omitempty"`
	PostDeploy  string                 `yaml:"post_deploy,omitempty"`
	PostDestroy string                 `yaml:"post_destroy,omitempty"`
	PreBuild    string                 `yaml:"pre_build,omitempty"`
	PreDeploy   string                 `yaml:"pre_deploy,omitempty"`
	PreDestroy  string                 `yaml:"pre_destroy,omitempty"`
	Smoke       string                 `yaml:"smoke,omitempty"`
	Provides    []string               `yaml:"provides,omitempty"`
	Templates   []interface{}          `yaml:"templates,omitempty"`
	Test        string                 `yaml:"test,omitempty"`
	Type        string                 `yaml:"type"`
	Version     string                 `yaml:"version"`
}

func (e *escapePlan) GetBuild() string {
	return e.Build
}
func (e *escapePlan) GetConsumes() []string {
	return e.Consumes
}
func (e *escapePlan) GetDepends() []string {
	return e.Depends
}
func (e *escapePlan) GetDescription() string {
	return e.Description
}
func (e *escapePlan) GetErrands() map[string]interface{} {
	return e.Errands
}
func (e *escapePlan) GetIncludes() []string {
	return e.Includes
}
func (e *escapePlan) GetInputs() []interface{} {
	return e.Inputs
}
func (e *escapePlan) GetLogo() string {
	return e.Logo
}
func (e *escapePlan) GetMetadata() map[string]string {
	return e.Metadata
}
func (e *escapePlan) GetOutputs() []interface{} {
	return e.Outputs
}
func (e *escapePlan) GetTemplates() []interface{} {
	return e.Templates
}
func (e *escapePlan) GetPath() string {
	return e.Path
}
func (e *escapePlan) GetPostBuild() string {
	return e.PostBuild
}
func (e *escapePlan) GetPostDeploy() string {
	return e.PostDeploy
}
func (e *escapePlan) GetPostDestroy() string {
	return e.PostDestroy
}
func (e *escapePlan) GetPreBuild() string {
	return e.PreBuild
}
func (e *escapePlan) GetPreDeploy() string {
	return e.PreDeploy
}
func (e *escapePlan) GetPreDestroy() string {
	return e.PreDestroy
}
func (e *escapePlan) GetProvides() []string {
	return e.Provides
}
func (e *escapePlan) GetTest() string {
	return e.Test
}
func (e *escapePlan) GetSmoke() string {
	return e.Smoke
}
func (e *escapePlan) GetType() string {
	return e.Type
}
func (e *escapePlan) GetVersion() string {
	return e.Version
}
func (e *escapePlan) SetBuild(newValue string) {
	e.Build = newValue
}
func (e *escapePlan) SetConsumes(newValue []string) {
	e.Consumes = newValue
}
func (e *escapePlan) SetDepends(newValue []string) {
	e.Depends = newValue
}
func (e *escapePlan) SetDescription(newValue string) {
	e.Description = newValue
}
func (e *escapePlan) SetErrands(newValue map[string]interface{}) {
	e.Errands = newValue
}
func (e *escapePlan) SetIncludes(newValue []string) {
	e.Includes = newValue
}
func (e *escapePlan) SetInputs(newValue []interface{}) {
	e.Inputs = newValue
}
func (e *escapePlan) SetLogo(newValue string) {
	e.Logo = newValue
}
func (e *escapePlan) SetMetadata(newValue map[string]string) {
	e.Metadata = newValue
}
func (e *escapePlan) SetOutputs(newValue []interface{}) {
	e.Outputs = newValue
}
func (e *escapePlan) SetPath(newValue string) {
	e.Path = newValue
}
func (e *escapePlan) SetPostBuild(newValue string) {
	e.PostBuild = newValue
}
func (e *escapePlan) SetPostDeploy(newValue string) {
	e.PostDeploy = newValue
}
func (e *escapePlan) SetPostDestroy(newValue string) {
	e.PostDestroy = newValue
}
func (e *escapePlan) SetPreBuild(newValue string) {
	e.PreBuild = newValue
}
func (e *escapePlan) SetPreDeploy(newValue string) {
	e.PreDeploy = newValue
}
func (e *escapePlan) SetPreDestroy(newValue string) {
	e.PreDestroy = newValue
}
func (e *escapePlan) SetProvides(newValue []string) {
	e.Provides = newValue
}
func (e *escapePlan) SetTest(newValue string) {
	e.Test = newValue
}
func (e *escapePlan) SetType(newValue string) {
	e.Type = newValue
}
func (e *escapePlan) SetVersion(newValue string) {
	e.Version = newValue
}

func NewEscapePlan() EscapePlan {
	return &escapePlan{
		Consumes: []string{},
		Provides: []string{},
		Depends:  []string{},
		Includes: []string{},
		Metadata: map[string]string{},
		Errands:  map[string]interface{}{},
	}
}

func (e *escapePlan) GetReleaseId() string {
	return e.Type + "-" + e.Build + "-v" + e.Version
}
func (e *escapePlan) GetVersionlessReleaseId() string {
	return e.Type + "-" + e.Build
}

func (e *escapePlan) LoadConfig(cfgFile string) error {
	if !util.PathExists(cfgFile) {
		return fmt.Errorf("Escape plan '%s' was not found. Use 'escape plan init' to create it", cfgFile)
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		errorString := err.Error()
		osErr, ok := err.(*os.PathError)
		if ok {
			errorString = osErr.Err.Error()
		}
		return fmt.Errorf("Couldn't read Escape plan '%s': %s", cfgFile, errorString)
	}
	err = yaml.Unmarshal(data, e)
	if err != nil {
		return fmt.Errorf("Couldn't parse Escape plan '%s': %s", cfgFile, err.Error())
	}
	return nil
}

func (e *escapePlan) Init(typ, buildId string) EscapePlan {
	e.Build = buildId
	e.Type = typ
	e.Version = "@"
	return e
}

func (e *escapePlan) ToYaml() []byte {
	pr := NewPrettyPrinter()
	return pr.Print(e)
}

func (e *escapePlan) ToDict() map[string]interface{} {
	str, err := yaml.Marshal(e)
	if err != nil {
		panic(err)
	}
	yamlMap := map[string]interface{}{}
	if err := yaml.Unmarshal(str, &yamlMap); err != nil {
		panic(err)
	}
	return yamlMap
}
