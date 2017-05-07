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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
)

type postdestroy_runner struct {
	Stage string
}

func NewPostDestroyRunner(stage string) Runner {
	return &postdestroy_runner{
		Stage: stage,
	}
}

func (p *postdestroy_runner) Run(ctx RunnerContext) error {
	scriptPath, err := initScript(ctx, p.Stage, "pre_destroy")
	if err != nil {
		return err
	}
	deploymentState, err := initDeploymentState(ctx, p.Stage, true)
	if err != nil {
		return err
	}
	ctx.SetBuildInputs(deploymentState.GetCalculatedInputs(p.Stage))
	ctx.SetBuildOutputs(deploymentState.GetCalculatedOutputs(p.Stage))
	if scriptPath == "" {
		return p.deleteState(ctx, deploymentState)
	}
	env := NewEnvironmentBuilder().MergeInputsAndOutputsWithOsEnvironment(ctx)
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(ctx.GetPath().GetBaseDir())
	if err := proc.Run([]string{scriptPath}, env, ctx.Logger()); err != nil {
		return err
	}
	return p.deleteState(ctx, deploymentState)
}

func (p *postdestroy_runner) deleteState(ctx RunnerContext, depl DeploymentState) error {
	if err := depl.SetVersion(p.Stage, ""); err != nil {
		return err
	}
	if err := depl.UpdateInputs(p.Stage, nil); err != nil {
		return err
	}
	if err := depl.UpdateOutputs(p.Stage, nil); err != nil {
		return err
	}
	return nil
}
