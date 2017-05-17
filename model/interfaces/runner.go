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
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type Runner interface {
	Run(RunnerContext) error
}

type RunnerContext interface {
	GetEnvironmentState() EnvironmentState
	GetDeploymentState() DeploymentState
	SetDeploymentState(DeploymentState)
	GetReleaseMetadata() *core.ReleaseMetadata
	SetReleaseMetadata(*core.ReleaseMetadata)
	Logger() util.Logger
	GetPath() Paths
	GetDepends() []string
	GetBuildInputs() *map[string]interface{}
	SetBuildInputs(*map[string]interface{})
	GetBuildOutputs() *map[string]interface{}
	SetBuildOutputs(*map[string]interface{})
	NewContextForDependency(*core.ReleaseMetadata) RunnerContext
	GetScriptEnvironment(stage string) (*script.ScriptEnvironment, error)
}
