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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

func ToScriptEnvironment(d DeploymentState, metadataMap map[string]*core.ReleaseMetadata, stage string) (*script.ScriptEnvironment, error) {
	result := map[string]script.Script{}
	result["this"] = toScript(d, metadataMap["this"], stage)
	providers := d.GetProviders()
	for _, consumes := range metadataMap["this"].GetConsumes() {
		deplName, ok := providers[consumes]
		if !ok {
			return nil, fmt.Errorf("No provider of type '%s' was configured in the deployment state.", consumes)
		}
		deplState, err := d.GetEnvironmentState().LookupDeploymentState(deplName)
		if err != nil {
			return nil, err
		}
		result[consumes] = toScript(deplState, metadataMap[consumes], "deploy")
	}
	// TODO: only add deployments for dependencies that are found in
	// release metadata
	for _, deplState := range d.GetDeployments() {
		key := deplState.GetRelease()
		metadata, ok := metadataMap[key]
		if !ok {
			return nil, fmt.Errorf("Couldn't find metadata for '%s'. This is a bug in Escape", key)
		}
		version := metadata.GetVersion()
		if deplState.IsDeployed("deploy", version) {
			result[key+"-v"+version] = toScript(deplState, metadataMap[key], "deploy")
		}
	}
	for key, metadata := range metadataMap {
		reference, exists := result[metadata.GetReleaseId()]
		if exists {
			result[key] = reference
		}
	}
	return script.NewScriptEnvironmentWithGlobals(result), nil

}

func toScript(d DeploymentState, metadata *core.ReleaseMetadata, stage string) script.Script {
	result := map[string]script.Script{}
	if metadata != nil {
		result = metadata.ToScriptMap()
	}
	inputs := map[string]interface{}{}
	for key, val := range d.GetCalculatedInputs(stage) {
		for _, defined := range metadata.GetInputs() {
			if key == defined.GetId() {
				inputs[key] = val
			}
		}
	}
	outputs := map[string]interface{}{}
	for key, val := range d.GetCalculatedOutputs(stage) {
		for _, defined := range metadata.GetOutputs() {
			if key == defined.GetId() {
				outputs[key] = val
			}
		}
	}
	result["inputs"] = script.LiftDict(liftScriptValues(inputs))
	result["outputs"] = script.LiftDict(liftScriptValues(outputs))
	env := d.GetEnvironmentState()
	result["project"] = script.LiftString(env.GetProjectName())
	result["environment"] = script.LiftString(env.GetName())
	result["deployment"] = script.LiftString(d.GetName())
	return script.LiftDict(result)
}

func liftScriptValues(values map[string]interface{}) map[string]script.Script {
	result := map[string]script.Script{}
	if values != nil {
		for key, val := range values {
			v, err := script.Lift(val)
			if err != nil {
				panic(err)
			}
			result[key] = v
		}
	}
	return result
}
