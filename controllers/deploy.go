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

package controllers

import (
	"os"

	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/model/runners"
	"github.com/ankyra/escape/model/runners/deploy"
)

type DeployController struct{}

func SetExtraProviders(context *model.Context, stage string, extraProviders map[string]string) error {
	envState := context.GetEnvironmentState()
	deplState, err := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	if err != nil {
		return err
	}
	metadata := context.GetReleaseMetadata()
	return deplState.ConfigureProviders(metadata, stage, extraProviders)
}

func SaveExtraInputsAndProvidersInDeploymentState(context *model.Context, stage string, extraVars map[string]interface{}, extraProviders map[string]string) error {
	envState := context.GetEnvironmentState()
	deplState, err := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	if err != nil {
		return err
	}
	inputs := deplState.GetUserInputs(stage)
	for key, val := range extraVars {
		inputs[key] = val
	}
	if err := SetExtraProviders(context, stage, extraProviders); err != nil {
		return err
	}
	return deplState.UpdateUserInputs(stage, inputs)
}

func (d DeployController) Deploy(context *model.Context, extraVars map[string]interface{}, extraProviders map[string]string) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetQualifiedReleaseId())
	context.PushLogSection("Deploy")
	context.Log("deploy.start", nil)
	if err := SaveExtraInputsAndProvidersInDeploymentState(context, "deploy", extraVars, extraProviders); err != nil {
		return MarkDeploymentFailed(context, err, state.Failure)
	}
	runnerContext, err := runners.NewRunnerContext(context)
	if err != nil {
		return MarkDeploymentFailed(context, err, state.Failure)
	}
	if err := deploy.NewDeployRunner().Run(runnerContext); err != nil {
		return err
	}
	context.Log("deploy.finished", map[string]string{
		"deployment":  context.GetRootDeploymentName(),
		"environment": context.GetEnvironmentState().Name,
	})
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (d DeployController) FetchAndDeploy(context *model.Context, releaseId string, extraVars map[string]interface{}, extraProviders map[string]string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return MarkDeploymentFailed(context, err, state.Failure)
	}
	fetcher := FetchController{}
	if err := fetcher.ResolveFetchAndLoad(context, releaseId); err != nil {
		os.Chdir(currentDir)
		return MarkDeploymentFailed(context, err, state.Failure)
	}
	if err := d.Deploy(context, extraVars, extraProviders); err != nil {
		os.Chdir(currentDir)
		return err
	}
	return os.Chdir(currentDir)
}

func MarkDeploymentFailed(context *model.Context, err error, errorCode state.StatusCode) error {
	if context.RootDeploymentName == "" {
		return err
	}
	envState := context.GetEnvironmentState()
	deplState, err2 := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	if err2 != nil {
		return err2
	}
	return deplState.SetFailureStatus(state.DeployStage, err, errorCode)
}
