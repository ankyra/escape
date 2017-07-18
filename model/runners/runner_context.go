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

package runners

import (
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/model/state"
	"github.com/ankyra/escape-client/model/state/types"
	"github.com/ankyra/escape-client/util"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type RunnerContext interface {
	GetEnvironmentState() *types.EnvironmentState
	GetDeploymentState() *types.DeploymentState
	SetDeploymentState(*types.DeploymentState)
	GetReleaseMetadata() *core.ReleaseMetadata
	SetReleaseMetadata(*core.ReleaseMetadata)
	Logger() util.Logger
	GetPath() *paths.Path
	GetBuildInputs() map[string]interface{}
	SetBuildInputs(map[string]interface{})
	GetBuildOutputs() map[string]interface{}
	SetBuildOutputs(map[string]interface{})
	NewContextForDependency(*core.ReleaseMetadata) RunnerContext
	GetScriptEnvironment(stage string) (*script.ScriptEnvironment, error)
	GetScriptEnvironmentForPreDependencyStep(stage string) (*script.ScriptEnvironment, error)
	GetRootDeploymentName() string
}

type runnerContext struct {
	environmentState *types.EnvironmentState
	deploymentState  *types.DeploymentState
	releaseMetadata  *core.ReleaseMetadata
	path             *paths.Path
	inputs           map[string]interface{}
	outputs          map[string]interface{}
	logger           util.Logger
	context          Context
	stage            string
}

func NewRunnerContext(context Context, rootStage string) (RunnerContext, error) {
	metadata := context.GetReleaseMetadata()
	if metadata == nil {
		return nil, fmt.Errorf("Missing metadata in context. This is a bug in Escape.")
	}
	deplState := context.GetEnvironmentState().GetOrCreateDeploymentState(context.GetRootDeploymentName())
	deplState.Release = metadata.GetVersionlessReleaseId()
	return &runnerContext{
		path:             paths.NewPath(),
		environmentState: context.GetEnvironmentState(),
		deploymentState:  deplState,
		releaseMetadata:  metadata,
		logger:           context.GetLogger(),
		context:          context,
		stage:            rootStage,
	}, nil
}

func (r *runnerContext) GetPath() *paths.Path {
	return r.path
}

func (r *runnerContext) GetEnvironmentState() *types.EnvironmentState {
	return r.environmentState
}

func (r *runnerContext) GetDeploymentState() *types.DeploymentState {
	return r.deploymentState
}

func (r *runnerContext) GetRootDeploymentName() string {
	return r.context.GetRootDeploymentName()
}

func (r *runnerContext) SetDeploymentState(d *types.DeploymentState) {
	r.deploymentState = d
}

func (r *runnerContext) Logger() util.Logger {
	return r.logger
}

func (r *runnerContext) GetReleaseMetadata() *core.ReleaseMetadata {
	return r.releaseMetadata
}

func (r *runnerContext) SetReleaseMetadata(m *core.ReleaseMetadata) {
	r.releaseMetadata = m
}

func (r *runnerContext) GetBuildInputs() map[string]interface{} {
	return r.inputs
}

func (r *runnerContext) SetBuildInputs(inputs map[string]interface{}) {
	r.inputs = inputs
}

func (r *runnerContext) GetBuildOutputs() map[string]interface{} {
	return r.outputs
}

func (r *runnerContext) SetBuildOutputs(outputs map[string]interface{}) {
	r.outputs = outputs
}

func (r *runnerContext) GetScriptEnvironment(stage string) (*script.ScriptEnvironment, error) {
	return state.ToScriptEnvironment(r.GetDeploymentState(), r.GetReleaseMetadata(), stage, r.context)
}

func (r *runnerContext) GetScriptEnvironmentForPreDependencyStep(stage string) (*script.ScriptEnvironment, error) {
	// should only contain metadata, parent inputs and providers
	return state.ToScriptEnvironmentForDependencyStep(r.GetDeploymentState(), r.GetReleaseMetadata(), stage, r.context)
}

func (r *runnerContext) NewContextForDependency(metadata *core.ReleaseMetadata) RunnerContext {
	return &runnerContext{
		environmentState: r.environmentState,
		deploymentState:  r.deploymentState.GetDeployment(r.stage, metadata.GetVersionlessReleaseId()),
		path:             r.path.NewPathForDependency(metadata),
		releaseMetadata:  metadata,
		logger:           r.logger,
		inputs:           r.inputs,
		outputs:          r.outputs,
		context:          r.context,
		stage:            "deploy",
	}
}
