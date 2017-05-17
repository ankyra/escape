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
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type runnerContext struct {
	environmentState EnvironmentState
	deploymentState  DeploymentState
	releaseMetadata  *core.ReleaseMetadata
	path             Paths
	inputs           *map[string]interface{}
	outputs          *map[string]interface{}
	depends          []*core.ReleaseMetadata
	logger           util.Logger
	context          Context
}

func NewRunnerContext(context Context) (RunnerContext, error) {
	metadata := context.GetReleaseMetadata()
	if metadata == nil {
		return nil, fmt.Errorf("Missing metadata in context. This is a bug in Escape.")
	}
	for _, consumer := range metadata.GetConsumes() {
		envState := context.GetEnvironmentState()
		deplState, err := envState.GetDeploymentState([]string{metadata.GetVersionlessReleaseId()})
		if err != nil {
			return nil, err
		}
		deplState, err = deplState.ResolveConsumer(consumer)
		if err != nil {
			return nil, err
		}
	}
	return &runnerContext{
		path:             paths.NewPath(),
		environmentState: context.GetEnvironmentState(),
		releaseMetadata:  context.GetReleaseMetadata(),
		logger:           context.GetLogger(),
		depends:          []*core.ReleaseMetadata{context.GetReleaseMetadata()},
		context:          context,
	}, nil
}

func (r *runnerContext) GetPath() Paths {
	return r.path
}
func (r *runnerContext) GetEnvironmentState() EnvironmentState {
	return r.environmentState
}
func (r *runnerContext) GetDeploymentState() DeploymentState {
	return r.deploymentState
}
func (r *runnerContext) GetDepends() []string {
	deps := []string{}
	for _, d := range r.depends {
		deps = append(deps, d.GetVersionlessReleaseId())
	}
	return deps
}
func (r *runnerContext) SetDeploymentState(d DeploymentState) {
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
func (r *runnerContext) GetBuildInputs() *map[string]interface{} {
	return r.inputs
}
func (r *runnerContext) SetBuildInputs(inputs *map[string]interface{}) {
	r.inputs = inputs
}
func (r *runnerContext) GetBuildOutputs() *map[string]interface{} {
	return r.outputs
}
func (r *runnerContext) SetBuildOutputs(outputs *map[string]interface{}) {
	r.outputs = outputs
}

func (r *runnerContext) GetScriptEnvironment(stage string) (*script.ScriptEnvironment, error) {
	if r.GetDeploymentState() == nil {
		return nil, fmt.Errorf("Missing deployment state in context. This is a bug in Escape.")
	}
	metadataCtx := map[string]*core.ReleaseMetadata{}
	for _, depend := range r.GetReleaseMetadata().GetDependencies() {
		metadata, err := r.context.GetDependencyMetadata(depend)
		if err != nil {
			return nil, err
		}
		metadataCtx[depend] = metadata
	}
	metadataCtx["this"] = r.GetReleaseMetadata()
	for key, ref := range r.GetReleaseMetadata().GetVariableContext() {
		metadata, err := r.context.GetDependencyMetadata(ref)
		if err != nil {
			return nil, err
		}
		previous, ok := metadataCtx[metadata.GetReleaseId()]
		if !ok {
			metadataCtx[key] = metadata
		} else {
			metadataCtx[key] = previous
		}
		metadataCtx[metadata.GetVersionlessReleaseId()] = metadataCtx[key]
	}
	return r.GetDeploymentState().ToScriptEnvironment(metadataCtx, stage)
}

func (r *runnerContext) NewContextForDependency(metadata *core.ReleaseMetadata) RunnerContext {
	return &runnerContext{
		environmentState: r.environmentState,
		deploymentState:  r.deploymentState,
		path:             r.path.NewPathForDependency(metadata),
		depends:          append(r.depends, metadata),
		releaseMetadata:  metadata,
		logger:           r.logger,
		inputs:           r.inputs,
		outputs:          r.outputs,
		context:          r.context,
	}
}
