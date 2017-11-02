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
	"io/ioutil"
	"os"

	"github.com/ankyra/escape/util"
	"github.com/ankyra/escape-core"
	"gopkg.in/yaml.v2"
)

type EscapePlan struct {
	Build           string                 `yaml:"build,omitempty"`
	BuildConsumes   []interface{}          `yaml:"build_consumes,omitempty"`
	BuildInputs     []interface{}          `yaml:"build_inputs,omitempty"`
	BuildTemplates  []interface{}          `yaml:"build_templates,omitempty"`
	Consumes        []interface{}          `yaml:"consumes,omitempty"`
	Depends         []interface{}          `yaml:"depends,omitempty"`
	Deploy          string                 `yaml:"deploy,omitempty"`
	DeployInputs    []interface{}          `yaml:"deploy_inputs,omitempty"`
	DeployConsumes  []interface{}          `yaml:"deploy_consumes,omitempty"`
	DeployTemplates []interface{}          `yaml:"deploy_templates,omitempty"`
	Destroy         string                 `yaml:"destroy,omitempty"`
	Description     string                 `yaml:"description,omitempty"`
	Downloads       []*core.DownloadConfig `yaml:"downloads,omitempty"`
	Extends         []string               `yaml:"extends,omitempty"`
	Errands         map[string]interface{} `yaml:"errands,omitempty"`
	Includes        []string               `yaml:"includes,omitempty"`
	Inputs          []interface{}          `yaml:"inputs,omitempty"`
	Logo            string                 `yaml:"logo,omitempty"`
	Metadata        map[string]string      `yaml:"metadata,omitempty"`
	Name            string                 `yaml:"name"`
	Outputs         []interface{}          `yaml:"outputs,omitempty"`
	Path            string                 `yaml:"path,omitempty"`
	PostBuild       string                 `yaml:"post_build,omitempty"`
	PostDeploy      string                 `yaml:"post_deploy,omitempty"`
	PostDestroy     string                 `yaml:"post_destroy,omitempty"`
	PreBuild        string                 `yaml:"pre_build,omitempty"`
	PreDeploy       string                 `yaml:"pre_deploy,omitempty"`
	PreDestroy      string                 `yaml:"pre_destroy,omitempty"`
	Smoke           string                 `yaml:"smoke,omitempty"`
	Provides        []string               `yaml:"provides,omitempty"`
	Templates       []interface{}          `yaml:"templates,omitempty"`
	Test            string                 `yaml:"test,omitempty"`
	Version         string                 `yaml:"version"`
}

func NewEscapePlan() *EscapePlan {
	return &EscapePlan{
		Consumes: []interface{}{},
		Provides: []string{},
		Depends:  []interface{}{},
		Includes: []string{},
		Metadata: map[string]string{},
		Errands:  map[string]interface{}{},
	}
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
	e.Version = "0.0.@"
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

func (e *EscapePlan) GetDependencies() ([]*core.DependencyConfig, error) {
	result := []*core.DependencyConfig{}
	for _, depend := range e.Depends {
		switch depend.(type) {
		case string:
			result = append(result, core.NewDependencyConfig(depend.(string)))
		case map[interface{}]interface{}:
			dep, err := core.NewDependencyConfigFromMap(depend.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			result = append(result, dep)
		default:
			return nil, fmt.Errorf("Invalid dependency format '%v' (expecting dict or string, got '%T')", depend, depend)
		}
	}
	return result, nil
}

func (e *EscapePlan) AddDependency(d *core.DependencyConfig) error {
	bytes, err := yaml.Marshal(d)
	if err != nil {
		return err
	}
	result := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(bytes, &result); err != nil {
		return err
	}
	e.Depends = append(e.Depends, result)
	return nil
}
