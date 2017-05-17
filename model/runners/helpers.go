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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
)

type ScriptStep struct {
	ShouldBeDeployed        bool
	ModifiesOutputVariables bool
	Stage                   string
	Step                    string
	Inputs                  func(ctx RunnerContext, stage string) (*map[string]interface{}, error)
	LoadOutputs             bool
	ScriptPath              string
	Commit                  func(ctx RunnerContext, state DeploymentState, stage string) error
}

func NewScriptStep(ctx RunnerContext, stage, step string, shouldBeDeployed bool) *ScriptStep {
	return &ScriptStep{
		ShouldBeDeployed:        shouldBeDeployed,
		Stage:                   stage,
		Step:                    step,
		Inputs:                  nil,
		LoadOutputs:             shouldBeDeployed,
		ScriptPath:              ctx.GetReleaseMetadata().GetScript(step),
		Commit:                  nil,
		ModifiesOutputVariables: false,
	}
}

func NewPreScriptStepRunner(stage, field string) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, false)
		step.Inputs = NewEnvironmentBuilder().GetInputsForPreStep
		step.Commit = preCommit
		return step.Run(ctx)
	})
}
func NewMainStepRunner(stage, field string) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, true)
		step.Commit = mainCommit
		step.ModifiesOutputVariables = true
		return step.Run(ctx)
	})
}
func NewPostScriptStepRunner(stage, field string) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, true)
		step.Commit = postCommit
		step.ModifiesOutputVariables = true
		return step.Run(ctx)
	})
}

func NewScriptRunner(stage, field string) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		return NewScriptStep(ctx, stage, field, true).Run(ctx)
	})
}

func compileTemplates(ctx RunnerContext, stage string) error {
	env, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return err
	}
	templates := ctx.GetReleaseMetadata().GetTemplates()
	for _, tpl := range templates {
		if err := tpl.Render(stage, env); err != nil {
			return err
		}
	}
	return nil
}

func preCommit(ctx RunnerContext, deploymentState DeploymentState, stage string) error {
	inputs := ctx.GetBuildInputs()
	version := ctx.GetReleaseMetadata().GetVersion()
	if err := deploymentState.SetVersion(stage, version); err != nil {
		return err
	}
	if err := deploymentState.UpdateInputs(stage, inputs); err != nil {
		return err
	}
	return compileTemplates(ctx, stage)
}

func mainCommit(ctx RunnerContext, deploymentState DeploymentState, stage string) error {
	return ctx.GetDeploymentState().UpdateOutputs(stage, ctx.GetBuildOutputs())
}

func postCommit(ctx RunnerContext, deploymentState DeploymentState, stage string) error {
	processedOutputs, err := NewEnvironmentBuilder().GetOutputs(ctx, stage)
	if err != nil {
		return err
	}
	return deploymentState.UpdateOutputs(stage, processedOutputs)
}

func (b *ScriptStep) Run(ctx RunnerContext) error {
	if b.ScriptPath != "" {
		scriptPath, err := b.initScript(ctx)
		if err != nil {
			return err
		}
		b.ScriptPath = scriptPath
	}
	deploymentState, err := b.initDeploymentState(ctx)
	if err != nil {
		return err
	}
	if b.ScriptPath != "" {
		if err := b.runScript(ctx); err != nil {
			return err
		}
	}
	if b.Commit != nil {
		return b.Commit(ctx, deploymentState, b.Stage)
	}
	return nil
}

func (b *ScriptStep) initScript(ctx RunnerContext) (string, error) {
	script := ctx.GetPath().Script(b.ScriptPath)
	ctx.Logger().Log(b.Stage+".step", map[string]string{
		"step":   b.Step,
		"script": script,
	})
	if !util.PathExists(script) {
		return "", fmt.Errorf("Referenced %s script '%s' does not exist", b.Step, script)
	}
	if err := util.MakeExecutable(script); err != nil {
		return "", err
	}
	return script, nil
}

func (b *ScriptStep) initDeploymentState(ctx RunnerContext) (DeploymentState, error) {
	deploymentState, err := ctx.GetEnvironmentState().GetDeploymentState(ctx.GetDepends())
	if err != nil {
		return nil, err
	}
	version := ctx.GetReleaseMetadata().GetVersion()
	if b.ShouldBeDeployed && !deploymentState.IsDeployed(b.Stage, version) {
		stageName := "Build"
		if b.Stage == "deploy" {
			stageName = "Deployment"
		}
		return nil, fmt.Errorf("%s state '%s' (version %s) could not be found", stageName, ctx.GetDepends()[0], version)
	}
	ctx.SetDeploymentState(deploymentState)
	if b.Inputs != nil {
		inputs, err := b.Inputs(ctx, b.Stage)
		if err != nil {
			return nil, err
		}
		ctx.SetBuildInputs(inputs)
	} else {
		ctx.SetBuildInputs(deploymentState.GetCalculatedInputs(b.Stage))
	}
	if b.LoadOutputs {
		ctx.SetBuildOutputs(deploymentState.GetCalculatedOutputs(b.Stage))
	}
	return deploymentState, nil
}

func (b *ScriptStep) getEnv(ctx RunnerContext) []string {
	if !b.LoadOutputs {
		return NewEnvironmentBuilder().MergeInputsWithOsEnvironment(ctx)
	}
	return NewEnvironmentBuilder().MergeInputsAndOutputsWithOsEnvironment(ctx)
}

func (b *ScriptStep) getCmd(ctx RunnerContext) ([]string, error) {
	if b.ModifiesOutputVariables {
		if err := writeOutputsToFile(ctx.GetBuildOutputs()); err != nil {
			return nil, err
		}
		outputsJsonLocation := ctx.GetPath().OutputsFile()
		return []string{b.ScriptPath, outputsJsonLocation}, nil
	}
	return []string{b.ScriptPath}, nil
}

func (b *ScriptStep) runScript(ctx RunnerContext) error {
	env := b.getEnv(ctx)
	cmd, err := b.getCmd(ctx)
	if err != nil {
		return err
	}
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(ctx.GetPath().GetBaseDir())
	if err := proc.Run(cmd, env, ctx.Logger()); err != nil {
		return err
	}
	return b.readOutputVariables(ctx)
}

func (b *ScriptStep) readOutputVariables(ctx RunnerContext) error {
	if !b.ModifiesOutputVariables {
		return nil
	}
	outputsJsonLocation := ctx.GetPath().OutputsFile()
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
	return nil
}

func writeOutputsToFile(outputs *map[string]interface{}) error {
	path := paths.NewPath()
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
