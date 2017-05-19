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

type EnvironmentState interface {
	GetProjectName() string
	GetName() string
	GetDeployments() []DeploymentState
	GetDeploymentState(deps []string) (DeploymentState, error)
	LookupDeploymentState(name string) (DeploymentState, error)
}

type DeploymentState interface {
	GetName() string
	GetVersion(stage string) string
	GetRelease() string
	GetEnvironmentState() EnvironmentState
	GetDeployments() []DeploymentState
	GetProviders() map[string]string
	IsDeployed(stage, version string) bool

	GetPreStepInputs(stage string) map[string]interface{}
	GetUserInputs(stage string) map[string]interface{}
	GetCalculatedInputs(stage string) map[string]interface{}
	GetCalculatedOutputs(stage string) map[string]interface{}

	UpdateUserInputs(stage string, v map[string]interface{}) error
	UpdateInputs(stage string, v map[string]interface{}) error
	UpdateOutputs(stage string, v map[string]interface{}) error
	SetVersion(stage, version string) error
	Save() error

	ToJson() string
}
