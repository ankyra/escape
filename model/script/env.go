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

package script

import (
	. "github.com/ankyra/escape-client/model/interfaces"
)

func NewScriptEnvironment() *ScriptEnvironment {
	result := ScriptEnvironment{}
	return &result
}

func NewScriptEnvironmentWithGlobals(globals map[string]Script) *ScriptEnvironment {
	result := ScriptEnvironment{}
	globalsDict := LiftDict(globals)
	result["$"] = globalsDict
	return &result
}

func NewScriptEnvironmentForCompileStep(variableCtx map[string]ReleaseMetadata) *ScriptEnvironment {
	result := map[string]Script{}
	for key, metadata := range variableCtx {
		result[key] = releaseMetadataToEscapeDict(metadata)
	}
	return NewScriptEnvironmentWithGlobals(result)
}

func NewScriptEnvironmentForStage(metadataCtx *map[string]ReleaseMetadata, deployCtx *map[string]DeploymentState, depl DeploymentState, stage string) *ScriptEnvironment {
	result := map[string]Script{}
	for key, deplState := range *deployCtx {
		metadata, _ := (*metadataCtx)[key]
		result[key] = deploymentStateToEscapeDict(deplState, metadata, stage)
	}
	for key, metadata := range *metadataCtx {
		reference, exists := result[metadata.GetReleaseId()]
		if exists {
			result[key] = reference
		}
	}
	return NewScriptEnvironmentWithGlobals(result)
}

func releaseMetadataToEscapeDict(metadata ReleaseMetadata) Script {
	return LiftDict(releaseMetadataToDict(metadata))
}

func deploymentStateToEscapeDict(deplState DeploymentState, metadata ReleaseMetadata, stage string) Script {
	result := releaseMetadataToDict(metadata)
	inputsDict := map[string]Script{}
	outputsDict := map[string]Script{}
	inputs := deplState.GetCalculatedInputs(stage)
	outputs := deplState.GetCalculatedOutputs(stage)
	if inputs != nil {
		for key, val := range *inputs {
			v, err := Lift(val)
			if err != nil {
				panic(err)
			}
			inputsDict[key] = v
		}
	}
	if outputs != nil {
		for key, val := range *outputs {
			v, err := Lift(val)
			if err != nil {
				panic(err)
			}
			outputsDict[key] = v
		}
	}
	result["inputs"] = LiftDict(inputsDict)
	result["outputs"] = LiftDict(outputsDict)
	env := deplState.GetEnvironmentState()
	prj := env.GetProjectState()
	result["project"] = LiftString(prj.GetName())
	result["environment"] = LiftString(env.GetName())
	result["deployment"] = LiftString(deplState.GetName())
	return LiftDict(result)
}

func releaseMetadataToDict(metadata ReleaseMetadata) map[string]Script {
	metadataDict := map[string]Script{}
	if metadata == nil {
		return map[string]Script{
			"metadata": LiftDict(metadataDict),
		}
	}
	for key, val := range metadata.GetMetadata() {
		metadataDict[key] = LiftString(val)
	}
	return map[string]Script{
		"metadata": LiftDict(metadataDict),

		"branch":      LiftString(metadata.GetBranch()),
		"description": LiftString(metadata.GetDescription()),
		"logo":        LiftString(metadata.GetLogo()),
		"build":       LiftString(metadata.GetName()),
		"revision":    LiftString(metadata.GetRevision()),
		"id":          LiftString(metadata.GetReleaseId()),
		"type":        LiftString(metadata.GetType()),
		"version":     LiftString(metadata.GetVersion()),
	}
}
