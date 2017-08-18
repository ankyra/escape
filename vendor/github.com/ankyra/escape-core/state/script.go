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

package state

import (
	"fmt"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type DeploymentResolver interface {
	GetDependencyMetadata(depend *core.DependencyConfig) (*core.ReleaseMetadata, error)
}

type deploymentResolver struct {
	resolver func(*core.DependencyConfig) (*core.ReleaseMetadata, error)
}

func (d *deploymentResolver) GetDependencyMetadata(depend *core.DependencyConfig) (*core.ReleaseMetadata, error) {
	return d.resolver(depend)
}

func newResolverFromMap(metaMap map[string]*core.ReleaseMetadata) DeploymentResolver {
	return &deploymentResolver{
		resolver: func(depend *core.DependencyConfig) (*core.ReleaseMetadata, error) {
			m, ok := metaMap[depend.ReleaseId]
			if !ok {
				return nil, fmt.Errorf("Metadata for '%s' not found", depend.ReleaseId)
			}
			return m, nil
		},
	}
}

func ToScriptEnvironment(d *DeploymentState, metadata *core.ReleaseMetadata, stage string, context DeploymentResolver) (*script.ScriptEnvironment, error) {
	if d == nil {
		return nil, fmt.Errorf("Missing deployment state. This is a bug in Escape.")
	}
	result, err := NewStateCompiler(context).Compile(d, metadata, stage)
	if err != nil {
		return nil, err
	}
	return script.NewScriptEnvironmentWithGlobals(script.ExpectDictAtom(result)), nil
}

func ToScriptEnvironmentForDependencyStep(d *DeploymentState, metadata *core.ReleaseMetadata, stage string, context DeploymentResolver) (*script.ScriptEnvironment, error) {
	if d == nil {
		return nil, fmt.Errorf("Missing deployment state. This is a bug in Escape.")
	}
	st := NewStateCompiler(context)
	st.DependencyInputsAreAvailable = false
	result, err := st.Compile(d, metadata, stage)
	if err != nil {
		return nil, err
	}
	return script.NewScriptEnvironmentWithGlobals(script.ExpectDictAtom(result)), nil
}

type StateCompiler struct {
	Result                       map[string]script.Script
	Resolver                     DeploymentResolver
	DependencyInputsAreAvailable bool
}

func NewStateCompiler(context DeploymentResolver) *StateCompiler {
	return &StateCompiler{
		Result:                       map[string]script.Script{},
		Resolver:                     context,
		DependencyInputsAreAvailable: true,
	}
}

func (s *StateCompiler) CompileDependencies(d *DeploymentState, metadata *core.ReleaseMetadata, stage string) error {
	for _, depend := range metadata.Depends {
		depMetadata, err := s.Resolver.GetDependencyMetadata(depend)
		if err != nil {
			return err
		}
		depState := d.GetDeploymentOrMakeNew(stage, depMetadata.GetVersionlessReleaseId())
		s.Result[depend.ReleaseId] = s.CompileState(depState, depMetadata, "deploy", s.DependencyInputsAreAvailable)
	}
	return nil
}

func (s *StateCompiler) CompileProviders(d *DeploymentState, metadata *core.ReleaseMetadata, stage string) error {
	providers := d.GetProviders(stage)
	for _, consumes := range metadata.GetConsumes(stage) {
		deplName, ok := providers[consumes]
		if !ok {
			return fmt.Errorf("No provider of type '%s' was configured in the deployment state.", consumes)
		}
		deplState, err := d.GetEnvironmentState().LookupDeploymentState(deplName)
		if err != nil {
			return err
		}
		depMetadata, err := s.Resolver.GetDependencyMetadata(core.NewDependencyConfig(deplState.GetReleaseId("deploy")))
		if err != nil {
			return err
		}
		s.Result[consumes] = s.CompileState(deplState, depMetadata, "deploy", true)
	}
	return nil
}

func (s *StateCompiler) CompileVariableCtx(metadata *core.ReleaseMetadata) {
	for key, ref := range metadata.GetVariableContext() {
		script, ok := s.Result[ref]
		if !ok {
			continue
		}
		s.Result[key] = script
	}
}

func (s *StateCompiler) Compile(d *DeploymentState, metadata *core.ReleaseMetadata, stage string) (script.Script, error) {
	s.Result["this"] = s.CompileState(d, metadata, stage, s.DependencyInputsAreAvailable)
	if err := s.CompileDependencies(d, metadata, stage); err != nil {
		return nil, err
	}
	if err := s.CompileProviders(d, metadata, stage); err != nil {
		return nil, err
	}
	s.CompileVariableCtx(metadata)
	return script.LiftDict(s.Result), nil
}

func (s *StateCompiler) CompileState(d *DeploymentState, metadata *core.ReleaseMetadata, stage string, includeVariables bool) script.Script {
	result := map[string]script.Script{}
	if metadata != nil {
		result = metadata.ToScriptMap()
	}
	if d == nil {
		return script.LiftDict(result)
	}
	if includeVariables {
		inputs := map[string]interface{}{}
		outputs := map[string]interface{}{}
		for key, val := range d.GetCalculatedInputs(stage) {
			for _, defined := range metadata.GetInputs(stage) {
				if key == defined.Id {
					inputs[key] = val
				}
			}
		}
		for key, val := range d.GetCalculatedOutputs(stage) {
			for _, defined := range metadata.GetOutputs(stage) {
				if key == defined.Id {
					outputs[key] = val
				}
			}
		}
		result["inputs"] = script.ShouldLift(inputs)
		result["outputs"] = script.ShouldLift(outputs)
	}
	env := d.GetEnvironmentState()
	result["project"] = script.LiftString(env.GetProjectName())
	result["environment"] = script.LiftString(env.Name)
	result["deployment"] = script.LiftString(d.Name)
	return script.LiftDict(result)
}
