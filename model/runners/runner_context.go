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

package runners

import (
	"fmt"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/state"
	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape/util"
)

type RunnerContext struct {
	environmentState *state.EnvironmentState
	deploymentState  *state.DeploymentState
	releaseMetadata  *core.ReleaseMetadata
	path             *paths.Path
	inputs           map[string]interface{}
	outputs          map[string]interface{}
	logger           util.Logger
	context          Context
	stage            string

	toScriptEnvironment func(d *state.DeploymentState, metadata *core.ReleaseMetadata, stage string, context state.DeploymentResolver) (*script.ScriptEnvironment, error)
}

func NewRunnerContext(context Context, rootStage string) (*RunnerContext, error) {
	metadata := context.GetReleaseMetadata()
	if metadata == nil {
		return nil, fmt.Errorf("Missing metadata in context. This is a bug in Escape.")
	}
	deplState, err := context.GetEnvironmentState().GetOrCreateDeploymentState(context.GetRootDeploymentName())
	if err != nil {
		return nil, err
	}
	deplState.Release = metadata.GetVersionlessReleaseId()
	return &RunnerContext{
		path:                paths.NewPath(),
		environmentState:    context.GetEnvironmentState(),
		deploymentState:     deplState,
		releaseMetadata:     metadata,
		logger:              context.GetLogger(),
		context:             context,
		stage:               rootStage,
		toScriptEnvironment: state.ToScriptEnvironment,
	}, nil
}

func (r *RunnerContext) GetPath() *paths.Path {
	return r.path
}

func (r *RunnerContext) GetEnvironmentState() *state.EnvironmentState {
	return r.environmentState
}

func (r *RunnerContext) GetDeploymentState() *state.DeploymentState {
	return r.deploymentState
}

func (r *RunnerContext) GetRootDeploymentName() string {
	return r.context.GetRootDeploymentName()
}

func (r *RunnerContext) SetDeploymentState(d *state.DeploymentState) {
	r.deploymentState = d
}

func (r *RunnerContext) Logger() util.Logger {
	return r.logger
}

func (r *RunnerContext) GetReleaseMetadata() *core.ReleaseMetadata {
	return r.releaseMetadata
}

func (r *RunnerContext) SetReleaseMetadata(m *core.ReleaseMetadata) {
	r.releaseMetadata = m
}

func (r *RunnerContext) GetBuildInputs() map[string]interface{} {
	return r.inputs
}

func (r *RunnerContext) SetBuildInputs(inputs map[string]interface{}) {
	r.inputs = inputs
}

func (r *RunnerContext) GetBuildOutputs() map[string]interface{} {
	return r.outputs
}

func (r *RunnerContext) SetBuildOutputs(outputs map[string]interface{}) {
	r.outputs = outputs
}

func (r *RunnerContext) GetScriptEnvironment(stage string) (*script.ScriptEnvironment, error) {
	return r.toScriptEnvironment(r.GetDeploymentState(), r.GetReleaseMetadata(), stage, r.context)
}

func (r *RunnerContext) GetScriptEnvironmentForPreDependencyStep(stage string) (*script.ScriptEnvironment, error) {
	// should only contain metadata, parent inputs and providers
	return state.ToScriptEnvironmentForDependencyStep(r.GetDeploymentState(), r.GetReleaseMetadata(), stage, r.context)
}

func (r *RunnerContext) NewContextForDependency(deploymentName string, metadata *core.ReleaseMetadata, consumerMapping map[string]string) (*RunnerContext, error) {
	depl, err := r.deploymentState.GetDeploymentOrMakeNew(r.stage, deploymentName)
	if err != nil {
		return nil, err
	}
	depl.Release = metadata.GetVersionlessReleaseId()

	scriptEnv, err := r.GetScriptEnvironment(r.stage)
	if err != nil {
		return nil, err
	}
	for iface, providerDepl := range consumerMapping {
		val, err := script.ParseAndEvalToGoValue(providerDepl, scriptEnv)
		if err != nil {
			return nil, err
		}
		valStr, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string for provider mapping '%s', but got '%v'", iface, val)
		}
		consumerMapping[iface] = valStr
	}
	return &RunnerContext{
		environmentState:    r.environmentState,
		deploymentState:     depl,
		path:                r.path.NewPathForDependency(metadata),
		releaseMetadata:     metadata,
		logger:              r.logger,
		inputs:              r.inputs,
		outputs:             r.outputs,
		context:             r.context,
		stage:               "deploy",
		toScriptEnvironment: r.toScriptEnvironment,
	}, depl.ConfigureProviders(metadata, "deploy", consumerMapping)
}
