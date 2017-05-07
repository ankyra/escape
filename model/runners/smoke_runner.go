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

type smokeRunner struct {
}

func NewSmokeRunner() Runner {
	return &smokeRunner{}
}

func (t *smokeRunner) Run(ctx RunnerContext) error {
	metadata := ctx.GetReleaseMetadata()
	state := ctx.GetEnvironmentState()
	deploymentState, err := state.GetDeploymentState(ctx.GetDepends())
	if err != nil {
		return err
	}
	version := ctx.GetReleaseMetadata().GetVersion()
	if !deploymentState.IsDeployed("deploy", version) {
		return fmt.Errorf("Deployment '%s' of version '%s' could not be found", ctx.GetDepends()[0], version)
	}
	if metadata.GetScript("smoke") == "" {
		return nil
	}
	ctx.SetBuildInputs(deploymentState.GetCalculatedInputs("deploy"))
	ctx.SetBuildOutputs(deploymentState.GetCalculatedOutputs("deploy"))
	ctx.SetDeploymentState(deploymentState)

	scriptPath := ctx.GetPath().Script(metadata.GetScript("smoke"))
	return runScript(ctx, scriptPath, "test")
}
