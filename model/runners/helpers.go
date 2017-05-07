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
	"encoding/json"
	"fmt"
	"github.com/ankyra/escape-client/model"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
)

func initScript(ctx RunnerContext, stage, field string) (string, error) {
	metadata := ctx.GetReleaseMetadata()
	scriptPath := metadata.GetScript(field)
	if scriptPath == "" {
		return "", nil
	}
	script := ctx.GetPath().Script(scriptPath)
	ctx.Logger().Log(stage+".step", map[string]string{
		"step":   field,
		"script": script,
	})
	if !util.PathExists(script) {
		return "", fmt.Errorf("Referenced %s script '%s' does not exist", field, script)
	}
	if err := util.MakeExecutable(script); err != nil {
		return "", err
	}
	return script, nil
}

func initDeploymentState(ctx RunnerContext, stage string, shouldBeDeployed bool) (DeploymentState, error) {
	deploymentState, err := ctx.GetEnvironmentState().GetDeploymentState(ctx.GetDepends())
	if err != nil {
		return nil, err
	}
	version := ctx.GetReleaseMetadata().GetVersion()
	if shouldBeDeployed && !deploymentState.IsDeployed(stage, version) {
		return nil, fmt.Errorf("Deployment '%s' could not be found", ctx.GetDepends()[0])
	}
	ctx.SetDeploymentState(deploymentState)
	ctx.SetBuildInputs(deploymentState.GetCalculatedInputs(stage))
	ctx.SetBuildOutputs(deploymentState.GetCalculatedOutputs(stage))
	return deploymentState, nil
}

func runPreScript(ctx RunnerContext, stage, field string, shouldBeDeployed bool) error {
	version := ctx.GetReleaseMetadata().GetVersion()
	scriptPath, err := initScript(ctx, stage, field)
	if err != nil {
		return err
	}
	deploymentState, err := initDeploymentState(ctx, stage, shouldBeDeployed)
	if err != nil {
		return err
	}
	inputs, err := NewEnvironmentBuilder().GetInputsForPreStep(ctx, stage)
	if err != nil {
		return err
	}
	ctx.SetBuildInputs(inputs)
	if scriptPath == "" {
		return updateAfterPreStep(stage, version, inputs, deploymentState)
	}
	env := NewEnvironmentBuilder().MergeInputsWithOsEnvironment(ctx)
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(ctx.GetPath().GetBaseDir())
	if err := proc.Run([]string{scriptPath}, env, ctx.Logger()); err != nil {
		return err
	}
	return updateAfterPreStep(stage, version, inputs, deploymentState)
}

func updateAfterPreStep(stage, version string, inputs *map[string]interface{}, deploymentState DeploymentState) error {
	if err := deploymentState.SetVersion(stage, version); err != nil {
		return err
	}
	return deploymentState.UpdateInputs(stage, inputs)
}

func updateAfterPostStep(ctx RunnerContext, stage string, outputs *map[string]interface{}, deploymentState DeploymentState) error {
	processedOutputs, err := NewEnvironmentBuilder().GetOutputs(ctx, stage)
	if err != nil {
		return err
	}
	return deploymentState.UpdateOutputs("deploy", processedOutputs)
}

func runPostScript(ctx RunnerContext, stage, field string) error {
	scriptPath, err := initScript(ctx, stage, field)
	if err != nil {
		return err
	}
	deploymentState, err := initDeploymentState(ctx, stage, true)
	if err != nil {
		return err
	}
	if scriptPath == "" {
		return updateAfterPostStep(ctx, stage, ctx.GetBuildOutputs(), deploymentState)
	}

	env := NewEnvironmentBuilder().MergeInputsAndOutputsWithOsEnvironment(ctx)
	if err := writeOutputsToFile(deploymentState.GetCalculatedOutputs(stage)); err != nil {
		return err
	}
	outputsJsonLocation := ctx.GetPath().OutputsFile()
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(ctx.GetPath().GetBaseDir())
	err = proc.Run([]string{scriptPath, outputsJsonLocation}, env, ctx.Logger())
	if err != nil {
		return err
	}
	outputOverrides, err := readOutputsFromFile(outputsJsonLocation)
	if err != nil {
		return err
	}
	outputs := ctx.GetBuildOutputs()
	if outputs == nil {
		outputs = &map[string]interface{}{}
	}
	for key, val := range outputOverrides {
		switch val.(type) {
		case string:
			break
		default:
			return fmt.Errorf("Expecting string value for output variable '%s'", key)
		}
		(*outputs)[key] = val
		//            applog("build.output_override_variable_value", variable=key, value=value)
	}
	ctx.SetBuildOutputs(outputs)
	return updateAfterPostStep(ctx, stage, outputs, deploymentState)
}

func runScript(ctx RunnerContext, script, field string) error {
	if !util.PathExists(script) {
		return fmt.Errorf("Referenced %s script '%s' does not exist", field, script)
	}
	if err := util.MakeExecutable(script); err != nil {
		return err
	}
	env := NewEnvironmentBuilder().MergeInputsAndOutputsWithOsEnvironment(ctx)
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(ctx.GetPath().GetBaseDir())
	return proc.Run([]string{script}, env, ctx.Logger())
}

func writeOutputsToFile(outputs *map[string]interface{}) error {
	path := model.NewPath()
	path.EnsureEscapeDirectoryExists()
	contents, err := json.Marshal(outputs)
	if err != nil {
		return err
	}
	outputsFile := path.OutputsFile()
	if err := ioutil.WriteFile(outputsFile, contents, 0644); err != nil {
		return err
	}
	return nil
}

func readOutputsFromFile(outputsJsonLocation string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if !util.PathExists(outputsJsonLocation) {
		return result, nil
	}
	payload, err := ioutil.ReadFile(outputsJsonLocation)
	if err != nil {
		return nil, err
	}
	if len(payload) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, err
	}
	return result, nil
}
