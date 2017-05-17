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

package interfaces

import (
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type ProjectState interface {
	ValidateAndFix() error
	GetEnvironmentStateOrMakeNew(string) EnvironmentState
	GetInputs() map[string]interface{}
	IsRemote() bool
	Save() error
	ToJson() string
	GetName() string
	SetName(name string)
}

type EnvironmentState interface {
	ValidateAndFix(string, ProjectState) error
	LookupDeploymentState(deploymentName string) (DeploymentState, error)
	GetDeploymentState(deps []string) (DeploymentState, error)
	GetProjectState() ProjectState
	GetInputs() map[string]interface{}
	IsRemote() bool
	Save() error
	GetDeployments() []DeploymentState
	GetName() string
}

type DeploymentState interface {
	ValidateAndFix(string, EnvironmentState) error

	GetPreStepInputs(stage string) *map[string]interface{}

	GetUserInputs(stage string) *map[string]interface{}
	GetCalculatedInputs(stage string) *map[string]interface{}
	GetCalculatedOutputs(stage string) *map[string]interface{}

	UpdateUserInputs(stage string, v *map[string]interface{}) error
	UpdateInputs(stage string, v *map[string]interface{}) error
	UpdateOutputs(stage string, v *map[string]interface{}) error

	ResolveConsumer(string) (DeploymentState, error)
	ToJson() string

	IsDeployed(stage, version string) bool
	GetVersion(stage string) string
	SetVersion(stage, version string) error
	GetName() string

	GetEnvironmentState() EnvironmentState
	ToScript(metadata *core.ReleaseMetadata, stage string) script.Script
	ToScriptEnvironment(rm map[string]*core.ReleaseMetadata, stage string) (*script.ScriptEnvironment, error)
}
