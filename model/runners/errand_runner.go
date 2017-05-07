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
)

type ErrandRunner struct {
	Errand Errand
}

func NewErrandRunner(errand Errand) Runner {
	return &ErrandRunner{
		Errand: errand,
	}
}

func (e *ErrandRunner) Run(ctx RunnerContext) error {
	if e.Errand.GetScript() == "" {
		return nil
	}
	state := ctx.GetEnvironmentState()
	deploymentState, err := state.GetDeploymentState(ctx.GetDepends())
	if err != nil {
		return err
	}
	if !deploymentState.IsDeployed("deploy", ctx.GetReleaseMetadata().GetVersion()) {
		return fmt.Errorf("Deployment '%s' could not be found", ctx.GetDepends()[0])
	}
	ctx.SetDeploymentState(deploymentState)

	inputs, err := NewEnvironmentBuilder().GetInputsForErrand(ctx, e.Errand)
	if err != nil {
		return err
	}
	outputs := deploymentState.GetCalculatedOutputs("deploy")
	ctx.SetBuildInputs(inputs)
	ctx.SetBuildOutputs(outputs)
	scriptPath := ctx.GetPath().Script(e.Errand.GetScript())
	return runScript(ctx, scriptPath, e.Errand.GetName())
}
