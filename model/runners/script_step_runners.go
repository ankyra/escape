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
	"encoding/json"
	"fmt"
	"io/ioutil"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model/dependency_resolvers"
	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape/util"
)

type ScriptStep struct {
	ShouldBeDeployed        bool
	ShouldDownload          bool
	ModifiesOutputVariables bool
	Stage                   string
	Step                    string
	Inputs                  func(ctx *RunnerContext, stage string) (map[string]interface{}, error)
	LoadOutputs             bool
	Script                  *core.ExecStage
	Commit                  func(ctx *RunnerContext, d *state.DeploymentState, stage string) error
}

func NewScriptStep(ctx *RunnerContext, stage, step string, shouldBeDeployed bool) *ScriptStep {
	return &ScriptStep{
		ShouldBeDeployed:        shouldBeDeployed,
		ShouldDownload:          false,
		Stage:                   stage,
		Step:                    step,
		Inputs:                  nil,
		LoadOutputs:             shouldBeDeployed,
		Script:                  ctx.GetReleaseMetadata().GetExecStage(step),
		Commit:                  nil,
		ModifiesOutputVariables: false,
	}
}

func NewPreScriptStepRunner(stage, field string, startCode, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, false)
		step.ShouldDownload = true
		step.Inputs = NewEnvironmentBuilder().GetInputsForPreStep
		step.Commit = preCommit
		return RunOrReportFailure(ctx, stage, step, startCode, errorCode)
	})
}
func NewMainStepRunner(stage, field string, startCode, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, true)
		step.Commit = mainCommit
		step.LoadOutputs = false
		step.ModifiesOutputVariables = true
		return RunOrReportFailure(ctx, stage, step, startCode, errorCode)
	})
}
func NewPostScriptStepRunner(stage, field string, startCode, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		step := NewScriptStep(ctx, stage, field, true)
		step.Commit = postCommit
		step.ModifiesOutputVariables = true
		return RunOrReportFailure(ctx, stage, step, startCode, errorCode)
	})
}

func NewScriptRunner(stage, field string, successCode, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		err := NewScriptStep(ctx, stage, field, true).Run(ctx)
		if err != nil {
			return ReportFailure(ctx, stage, err, errorCode)
		}
		st := state.NewStatus(successCode)
		return ctx.GetDeploymentState().UpdateStatus(stage, st)
	})
}

func compileTemplates(ctx *RunnerContext, stage string) error {
	env, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return err
	}
	templates := ctx.GetReleaseMetadata().GetTemplates(stage)
	for _, tpl := range templates {
		if err := tpl.Render(stage, env); err != nil {
			return err
		}
	}
	return nil
}

func preCommit(ctx *RunnerContext, deploymentState *state.DeploymentState, stage string) error {
	inputs := ctx.GetBuildInputs()
	metadata := ctx.GetReleaseMetadata()
	if err := deploymentState.CommitVersion(stage, metadata); err != nil {
		return err
	}
	if err := deploymentState.UpdateInputs(stage, inputs); err != nil {
		return err
	}
	return compileTemplates(ctx, stage)
}

func mainCommit(ctx *RunnerContext, deploymentState *state.DeploymentState, stage string) error {
	return ctx.GetDeploymentState().UpdateOutputs(stage, ctx.GetBuildOutputs())
}

func postCommit(ctx *RunnerContext, deploymentState *state.DeploymentState, stage string) error {
	processedOutputs, err := NewEnvironmentBuilder().GetOutputs(ctx, stage)
	if err != nil {
		return err
	}
	return deploymentState.UpdateOutputs(stage, processedOutputs)
}

func (b *ScriptStep) Run(ctx *RunnerContext) error {
	ctx.GetPath().EnsureEscapeDirectoryExists()
	if b.Script != nil && !b.Script.IsEmpty() {
		if err := b.initScript(ctx); err != nil {
			return err
		}
	}
	deploymentState, err := b.initDeploymentState(ctx)
	if err != nil {
		return err
	}
	if err := b.handleDownloads(ctx); err != nil {
		return err
	}
	if b.Script != nil && !b.Script.IsEmpty() {
		if err := b.runScript(ctx); err != nil {
			return err
		}
	}
	if b.Commit != nil {
		return b.Commit(ctx, deploymentState, b.Stage)
	}
	return nil
}

func (b *ScriptStep) initScript(ctx *RunnerContext) error {
	ctx.Logger().Log(b.Stage+".step", map[string]string{
		"step":   b.Step,
		"script": b.Script.String(),
	})
	if b.Script.RelativeScript == "" {
		return nil
	}
	script := b.Script.RelativeScript
	if !util.PathExists(script) {
		return fmt.Errorf("Referenced %s script '%s' does not exist", b.Step, script)
	}
	return util.MakeExecutable(script)
}

func (b *ScriptStep) initDeploymentState(ctx *RunnerContext) (*state.DeploymentState, error) {
	deploymentState := ctx.GetDeploymentState()

	metadata := ctx.GetReleaseMetadata()
	if b.ShouldBeDeployed && !deploymentState.IsDeployed(b.Stage, metadata) {
		var cmd string
		stageName := "Build"

		if b.Stage == "deploy" {
			stageName = "Deployment"
			if ctx.GetRootDeploymentName() == metadata.GetVersionlessReleaseId() {
				cmd = fmt.Sprintf("escape run deploy %s", metadata.GetReleaseId())
			} else {
				cmd = fmt.Sprintf("escape run deploy -d %s %s", ctx.GetRootDeploymentName(), metadata.GetReleaseId())
			}
		} else {
			if ctx.GetRootDeploymentName() == metadata.GetVersionlessReleaseId() {
				cmd = "escape run build"
			} else {
				cmd = fmt.Sprintf("escape run build -d %s", ctx.GetRootDeploymentName())
			}
		}
		return nil, fmt.Errorf("%s state '%s' for release '%s' could not be found\n\nYou may need to run `%s` to resolve this issue",
			stageName, ctx.GetRootDeploymentName(), metadata.GetReleaseId(), cmd)
	}
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

func (b *ScriptStep) getEnv(ctx *RunnerContext) []string {
	if !b.LoadOutputs {
		return NewEnvironmentBuilder().MergeInputsWithOsEnvironment(ctx)
	}
	return NewEnvironmentBuilder().MergeInputsAndOutputsWithOsEnvironment(ctx)
}

func (b *ScriptStep) handleDownloads(ctx *RunnerContext) error {
	if !b.ShouldDownload {
		return nil
	}
	downloads := ctx.GetReleaseMetadata().GetDownloads(b.Stage)
	return dependency_resolvers.DoDownloads(downloads, ctx.Logger())
}

func (b *ScriptStep) getCmd(ctx *RunnerContext) ([]string, error) {
	if b.ModifiesOutputVariables {
		if err := writeOutputsToFile(ctx.GetBuildOutputs()); err != nil {
			return nil, err
		}
	}
	return b.Script.GetAsCommand(), nil
}

func (b *ScriptStep) runScript(ctx *RunnerContext) error {
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

func (b *ScriptStep) readOutputVariables(ctx *RunnerContext) error {
	if !b.ModifiesOutputVariables {
		return nil
	}
	outputsJsonLocation := ctx.GetPath().OutputsFile()
	outputOverrides, err := readOutputsFromFile(outputsJsonLocation)
	if err != nil {
		return err
	}
	outputs := ctx.GetBuildOutputs()
	if outputs == nil || !b.LoadOutputs {
		outputs = map[string]interface{}{}
	}
	for key, val := range outputOverrides {
		switch val.(type) {
		case string:
			break
		default:
			return fmt.Errorf("Expecting string value for output variable '%s'", key)
		}
		outputs[key] = val
	}
	ctx.SetBuildOutputs(outputs)
	return nil
}

func writeOutputsToFile(outputs map[string]interface{}) error {
	path := paths.NewPath()
	path.EnsureEscapeDirectoryExists()
	if outputs == nil {
		outputs = map[string]interface{}{}
	}
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
