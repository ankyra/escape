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

package escape_plan

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape/util"
	"gopkg.in/yaml.v2"
)

const EscapePlanInitTemplate = `name: %s
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: 
deploy:

`

// Everything starts with a plan. An Escape plan.
//
// The Escape plan gets compiled into release metadata at build time.
//
type EscapePlan struct {
	// The package name is a required field. The name can be qualified by a
	// project name, but if no project is specified then the default project `_`
	// will be used.
	//
	// Format: `/([a-za-z]+[a-za-z0-9-]*\/)?[a-za-z]+[a-za-z0-9-]*/`
	//
	// Examples:
	//
	// * Fully qualified: `name: my-project/my-package`
	//
	// * Default project: `name: my-package`
	//
	Name string `yaml:"name"`

	// The version is a required field. Escape uses semantic versioning to
	// version packages.  Either specify the full version or use the '@' symbol
	// to let Escape pick the next version at build time. See
	// [here](/docs/guides/versioning/) for more versioning approaches.
	//
	// Format: `/[0-9]+(\.[0-9]+)*(\.@)?/`
	//
	// Examples:
	//
	// * Build version 1.5: `version: 1.5`
	//
	// * Build the next minor release in the 1.* series: `version: 1.@`
	//
	// * Build the next patch release in the 1.1.* series: `version: 1.1.@`
	//
	Version string `yaml:"version"`

	// A description for this package. Only used for presentation purposes.
	Description string `yaml:"description,omitempty"`

	// A path to an image. Only used for presentation purposes.
	Logo string `yaml:"logo,omitempty"`

	// The license. For example `Apache Software License`, `BSD License`, `GPLv3`, etc.
	// Currently no input validation is performed on this field.
	License string `yaml:"license,omitempty"`

	// Metadata key value pairs.
	//
	// [Escape Script](/docs/scripting-language/) can be used to
	// programmatically set values using the [default
	// context](/docs/scripting-language/#context).
	//
	// Example:
	//
	//   metadata:
	//     author: Fictional Character
	//     co_author: $dependency.metadata.author
	//
	Metadata map[string]string `yaml:"metadata,omitempty"`

	// Reference depedencies by their full ID or use the `@` symbol to resolve
	// versions at build time.
	Depends []interface{} `yaml:"depends,omitempty"`

	Extends []string `yaml:"extends,omitempty"`

	// The files to includes in this release. The files don't have to exist and can
	// be produced during build time. Globbing patterns are supported. Directories
	// are added recursively.
	//
	Includes []string `yaml:"includes,omitempty"`

	// Files that are generated during the build phase. Globbing patterns are
	// supported.  Directories are added recursively. The main reason to use
	// this over `includes` is that the `generates` field is copied to the
	// parent release, when a release gets extended, but `includes` aren't.
	Generates []string `yaml:"generates,omitempty"`

	// The release can declare zero or more providers so that consumers
	// can loosely depend on it at deploy time.
	Provides []string `yaml:"provides,omitempty"`

	// At deploy time a package can consume zero or more providers from the
	// target environment.
	Consumes []interface{} `yaml:"consumes,omitempty"`

	// Same as `consumes`, but scoped to the build stage (ie. the consumer is
	// not required/available at deploy time).
	BuildConsumes []interface{} `yaml:"build_consumes,omitempty"`

	// Same as `consumes`, but scoped to the deploy stage (ie. the consumer is
	// not required/available at build time).
	DeployConsumes []interface{} `yaml:"deploy_consumes,omitempty"`

	// Input variables.
	Inputs []interface{} `yaml:"inputs,omitempty"`

	// Same as `inputs`, but all variables are scoped to the build phase (ie. the
	// variables won't be required/available at deploy time).
	BuildInputs []interface{} `yaml:"build_inputs,omitempty"`

	// Same as `inputs`, but all variables are scoped to the deployment phase (ie. the
	// variables won't be required/available at build time).
	DeployInputs []interface{} `yaml:"deploy_inputs,omitempty"`

	// Output variables.
	Outputs []interface{} `yaml:"outputs,omitempty"`

	// Build script.
	Build interface{} `yaml:"build,omitempty"`

	// Pre-build script. The script has access to all the build scoped input
	// variables.
	PreBuild interface{} `yaml:"pre_build,omitempty"`

	// Post-build script. The script has access to all the build scoped input
	// and output variables.
	PostBuild interface{} `yaml:"post_build,omitempty"`

	// Test script.  Generally run after a build as part of the release
	// process, but can be triggered separately using `escape run test`.  The
	// script has access to all the build scoped input and output variables.
	Test interface{} `yaml:"test,omitempty"`

	// Deploy script. The script has access to the deployment input variables,
	// and can define outputs by writing a JSON object to .escape/outputs.json.
	Deploy interface{} `yaml:"deploy,omitempty"`

	// Pre-deploy script. The script has access to all the deploy scoped input
	// variables.
	PreDeploy interface{} `yaml:"pre_deploy,omitempty"`

	// Post-deploy script. The script has access to all the deploy scoped input
	// and output variables.
	PostDeploy interface{} `yaml:"post_deploy,omitempty"`

	// Activate provider script. This script is run when this release is being
	// consumed as a provider by another release during a build or deployment.
	// The script has access to all the deploy scoped input and output
	// variables.
	ActivateProvider interface{} `yaml:"activate_provider,omitempty"`

	// Deactive provider script. This script is run when this release is being
	// done being consumed by another release using it as a provider. The
	// script has access to all the deploy scoped input and output variables.
	DeactivateProvider interface{} `yaml:"deactivate_provider,omitempty"`

	// Smoke test script.
	Smoke interface{} `yaml:"smoke,omitempty"`

	// Destroy script.
	Destroy interface{} `yaml:"destroy,omitempty"`

	// Pre-destroy script.
	PreDestroy interface{} `yaml:"pre_destroy,omitempty"`

	// Post-destroy script.
	PostDestroy interface{} `yaml:"post_destroy,omitempty"`

	// Errands are scripts that can be run against the deployment of this release.
	// The scripts receive the deployment's inputs and outputs as environment
	// variables.
	Errands map[string]interface{} `yaml:"errands,omitempty"`

	// Templates.
	Templates []interface{} `yaml:"templates,omitempty"`

	// Same as `templates`, but all the templates are scoped to the build stage
	// (ie. templates won't be rendered at deploy time).
	BuildTemplates []interface{} `yaml:"build_templates,omitempty"`

	// Same as `templates`, but all the templates are scoped to the deploy stage
	// (ie. templates won't be rendered at deploy time).
	DeployTemplates []interface{} `yaml:"deploy_templates,omitempty"`

	// Downloads.
	Downloads []*core.DownloadConfig `yaml:"downloads,omitempty"`
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
		return fmt.Errorf("Escape plan '%s' was not found. Use 'escape plan init' to create it.", cfgFile)
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

func (e *EscapePlan) ToInitTemplate() []byte {
	return []byte(fmt.Sprintf(EscapePlanInitTemplate, e.Name))
}

func (e *EscapePlan) ToYaml() []byte {
	pr := NewPrettyPrinter()
	return pr.Print(e)
}

func (e *EscapePlan) ToMinifiedYaml() []byte {
	pr := NewPrettyPrinter(
		includeEmpty(false),
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
