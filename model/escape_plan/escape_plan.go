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
	"github.com/ankyra/escape-client/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type EscapePlan struct {
	Build       string                 `yaml:"build,omitempty"`
	Consumes    []string               `yaml:"consumes,omitempty"`
	Depends     []string               `yaml:"depends,omitempty"`
	Deploy      string                 `yaml:"deploy,omitempty"`
	Destroy     string                 `yaml:"destroy,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Extends     []string               `yaml:"extends,omitempty"`
	Errands     map[string]interface{} `yaml:"errands,omitempty"`
	Includes    []string               `yaml:"includes,omitempty"`
	Inputs      []interface{}          `yaml:"inputs,omitempty"`
	Logo        string                 `yaml:"logo,omitempty"`
	Metadata    map[string]string      `yaml:"metadata,omitempty"`
	Name        string                 `yaml:"name"`
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
	Version     string                 `yaml:"version"`
}

func NewEscapePlan() *EscapePlan {
	return &EscapePlan{
		Consumes: []string{},
		Provides: []string{},
		Depends:  []string{},
		Includes: []string{},
		Metadata: map[string]string{},
		Errands:  map[string]interface{}{},
	}
}

func (e *EscapePlan) GetName() string {
	return e.Name
}
func (e *EscapePlan) GetConsumes() []string {
	return e.Consumes
}
func (e *EscapePlan) GetDepends() []string {
	return e.Depends
}
func (e *EscapePlan) GetExtends() []string {
	return e.Extends
}
func (e *EscapePlan) GetDescription() string {
	return e.Description
}
func (e *EscapePlan) GetErrands() map[string]interface{} {
	return e.Errands
}
func (e *EscapePlan) GetIncludes() []string {
	return e.Includes
}
func (e *EscapePlan) GetInputs() []interface{} {
	return e.Inputs
}
func (e *EscapePlan) GetLogo() string {
	return e.Logo
}
func (e *EscapePlan) GetMetadata() map[string]string {
	return e.Metadata
}
func (e *EscapePlan) GetOutputs() []interface{} {
	return e.Outputs
}
func (e *EscapePlan) GetTemplates() []interface{} {
	return e.Templates
}
func (e *EscapePlan) GetPath() string {
	return e.Path
}
func (e *EscapePlan) GetBuild() string {
	return e.Build
}
func (e *EscapePlan) GetDestroy() string {
	return e.Destroy
}
func (e *EscapePlan) GetDeploy() string {
	return e.Deploy
}
func (e *EscapePlan) GetPostBuild() string {
	return e.PostBuild
}
func (e *EscapePlan) GetPostDeploy() string {
	return e.PostDeploy
}
func (e *EscapePlan) GetPostDestroy() string {
	return e.PostDestroy
}
func (e *EscapePlan) GetPreBuild() string {
	return e.PreBuild
}
func (e *EscapePlan) GetPreDeploy() string {
	return e.PreDeploy
}
func (e *EscapePlan) GetPreDestroy() string {
	return e.PreDestroy
}
func (e *EscapePlan) GetProvides() []string {
	return e.Provides
}
func (e *EscapePlan) GetTest() string {
	return e.Test
}
func (e *EscapePlan) GetSmoke() string {
	return e.Smoke
}
func (e *EscapePlan) GetVersion() string {
	return e.Version
}
func (e *EscapePlan) SetName(newValue string) {
	e.Name = newValue
}
func (e *EscapePlan) SetConsumes(newValue []string) {
	e.Consumes = newValue
}
func (e *EscapePlan) SetDepends(newValue []string) {
	e.Depends = newValue
}
func (e *EscapePlan) SetDescription(newValue string) {
	e.Description = newValue
}
func (e *EscapePlan) SetErrands(newValue map[string]interface{}) {
	e.Errands = newValue
}
func (e *EscapePlan) SetIncludes(newValue []string) {
	e.Includes = newValue
}
func (e *EscapePlan) SetInputs(newValue []interface{}) {
	e.Inputs = newValue
}
func (e *EscapePlan) SetLogo(newValue string) {
	e.Logo = newValue
}
func (e *EscapePlan) SetMetadata(newValue map[string]string) {
	e.Metadata = newValue
}
func (e *EscapePlan) SetOutputs(newValue []interface{}) {
	e.Outputs = newValue
}
func (e *EscapePlan) SetPath(newValue string) {
	e.Path = newValue
}
func (e *EscapePlan) SetPostBuild(newValue string) {
	e.PostBuild = newValue
}
func (e *EscapePlan) SetPostDeploy(newValue string) {
	e.PostDeploy = newValue
}
func (e *EscapePlan) SetPostDestroy(newValue string) {
	e.PostDestroy = newValue
}
func (e *EscapePlan) SetPreBuild(newValue string) {
	e.PreBuild = newValue
}
func (e *EscapePlan) SetPreDeploy(newValue string) {
	e.PreDeploy = newValue
}
func (e *EscapePlan) SetPreDestroy(newValue string) {
	e.PreDestroy = newValue
}
func (e *EscapePlan) SetProvides(newValue []string) {
	e.Provides = newValue
}
func (e *EscapePlan) SetTest(newValue string) {
	e.Test = newValue
}
func (e *EscapePlan) SetVersion(newValue string) {
	e.Version = newValue
}
func (e *EscapePlan) GetReleaseId() string {
	return e.Name + "-v" + e.Version
}
func (e *EscapePlan) GetVersionlessReleaseId() string {
	return e.Name
}

func (e *EscapePlan) LoadConfig(cfgFile string) error {
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

func (e *EscapePlan) Init(name string) *EscapePlan {
	e.Name = name
	e.Version = "@"
	return e
}

func (e *EscapePlan) ToYaml() []byte {
	pr := NewPrettyPrinter()
	return pr.Print(e)
}

func (e *EscapePlan) ToMinifiedYaml() []byte {
	pr := NewPrettyPrinter(
		includeEmpty(false),
		includeDocs(false),
		spacing(1),
	)
	return pr.Print(e)
}

func (e *EscapePlan) ToDict() map[string]interface{} {
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
