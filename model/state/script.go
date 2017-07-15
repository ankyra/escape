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
	. "github.com/ankyra/escape-client/model/state/types"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type DeploymentResolver interface {
	GetDependencyMetadata(depend string) (*core.ReleaseMetadata, error)
}

type deploymentResolver struct {
	resolver func(string) (*core.ReleaseMetadata, error)
}

func (d *deploymentResolver) GetDependencyMetadata(depend string) (*core.ReleaseMetadata, error) {
	return d.resolver(depend)
}

func newResolverFromMap(metaMap map[string]*core.ReleaseMetadata) DeploymentResolver {
	return &deploymentResolver{
		resolver: func(depend string) (*core.ReleaseMetadata, error) {
			m, ok := metaMap[depend]
			if !ok {
				return nil, fmt.Errorf("Metadata for '%s' not found", depend)
			}
			return m, nil
		},
	}
}

func ToScriptEnvironment(d *DeploymentState, metadata *core.ReleaseMetadata, stage string, context DeploymentResolver) (*script.ScriptEnvironment, error) {
	if d == nil {
		return nil, fmt.Errorf("Missing deployment state. This is a bug in Escape.")
	}
	result, err := ToScript(d, metadata, stage, context)
	if err != nil {
		return nil, err
	}
	return script.NewScriptEnvironmentWithGlobals(script.ExpectDictAtom(result)), nil
}

func ToScript(d *DeploymentState, metadata *core.ReleaseMetadata, stage string, context DeploymentResolver) (script.Script, error) {
	result := map[string]script.Script{}
	result["this"] = toScript(d, metadata, stage)

	for _, depend := range metadata.Depends {
		depMetadata, err := context.GetDependencyMetadata(depend.ReleaseId)
		if err != nil {
			return nil, err
		}
		depState := d.GetDeployment(stage, depMetadata.GetVersionlessReleaseId())
		result[depend.ReleaseId] = toScript(depState, depMetadata, "deploy")
	}

	providers := d.GetProviders(stage)
	for _, consumes := range metadata.GetConsumes() {
		deplName, ok := providers[consumes]
		if !ok {
			return nil, fmt.Errorf("No provider of type '%s' was configured in the deployment state.", consumes)
		}
		deplState, err := d.GetEnvironmentState().LookupDeploymentState(deplName)
		if err != nil {
			return nil, err
		}
		depMetadata, err := context.GetDependencyMetadata(deplState.GetReleaseId("deploy"))
		if err != nil {
			return nil, err
		}
		result[consumes] = toScript(deplState, depMetadata, "deploy")
	}

	for key, ref := range metadata.GetVariableContext() {
		script, ok := result[ref]
		if !ok {
			continue
		}
		result[key] = script
	}
	return script.LiftDict(result), nil
}

func toScript(d *DeploymentState, metadata *core.ReleaseMetadata, stage string) script.Script {
	result := map[string]script.Script{}
	if metadata != nil {
		result = metadata.ToScriptMap()
	}
	if d == nil {
		return script.LiftDict(result)
	}
	inputs := map[string]interface{}{}
	outputs := map[string]interface{}{}
	for key, val := range d.GetCalculatedInputs(stage) {
		for _, defined := range metadata.GetInputs() {
			if key == defined.GetId() {
				inputs[key] = val
			}
		}
	}
	for key, val := range d.GetCalculatedOutputs(stage) {
		for _, defined := range metadata.GetOutputs() {
			if key == defined.GetId() {
				outputs[key] = val
			}
		}
	}
	result["inputs"] = script.ShouldLift(inputs)
	result["outputs"] = script.ShouldLift(outputs)
	env := d.GetEnvironmentState()
	result["project"] = script.LiftString(env.GetProjectName())
	result["environment"] = script.LiftString(env.GetName())
	result["deployment"] = script.LiftString(d.GetName())
	return script.LiftDict(result)
}
