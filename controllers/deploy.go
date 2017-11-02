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

package controllers

import (
	"fmt"
	"os"

	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape/model/runners"
	"github.com/ankyra/escape/model/runners/deploy"
)

type DeployController struct{}

func SetExtraProviders(context Context, stage string, extraProviders map[string]string) error {
	envState := context.GetEnvironmentState()
	deplState := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	metadata := context.GetReleaseMetadata()
	configuredProviders := deplState.GetProviders(stage)
	availableProviders := envState.GetProviders()
	for _, c := range metadata.GetConsumes(stage) {
		provider, override := extraProviders[c]
		if override {
			deplState.SetProvider(stage, c, provider)
			continue
		}
		_, configured := configuredProviders[c]
		if configured {
			continue
		}
		implementations := availableProviders[c]
		if len(implementations) == 1 {
			deplState.SetProvider(stage, c, implementations[0])
		} else {
			return fmt.Errorf("Missing provider of type '%s'. This can be configured using the -p / --extra-provider flag.", c)
		}
	}
	return nil
}

func SaveExtraInputsAndProvidersInDeploymentState(context Context, stage string, extraVars, extraProviders map[string]string) error {
	envState := context.GetEnvironmentState()
	deplState := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	inputs := deplState.GetUserInputs(stage)
	for key, val := range extraVars {
		inputs[key] = val
	}
	if err := SetExtraProviders(context, stage, extraProviders); err != nil {
		return err
	}
	return deplState.UpdateUserInputs(stage, inputs)
}

func (d DeployController) Deploy(context Context, extraVars, extraProviders map[string]string) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Deploy")
	context.Log("deploy.start", nil)
	if err := SaveExtraInputsAndProvidersInDeploymentState(context, "deploy", extraVars, extraProviders); err != nil {
		return err
	}
	runnerContext, err := runners.NewRunnerContext(context, "deploy")
	if err != nil {
		return err
	}
	if err := deploy.NewDeployRunner().Run(runnerContext); err != nil {
		return err
	}
	if err := (SmokeController{}).Smoke(context); err != nil {
		return err
	}
	context.Log("deploy.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (d DeployController) FetchAndDeploy(context Context, releaseId string, extraVars, extraProviders map[string]string) error {
	// TODO cd into temp directory
	parsed, err := parsers.ParseQualifiedReleaseId(releaseId)
	if err != nil {
		return err
	}
	if parsed.NeedsResolving() {
		metadata, err := context.QueryReleaseMetadata(core.NewDependencyFromQualifiedReleaseId(parsed))
		if err != nil {
			return err
		}
		parsed.Version = metadata.Version
		metadata.Project = parsed.Project // registry needs to be updated to latest core
		releaseId = metadata.GetQualifiedReleaseId()
	}
	fetcher := FetchController{}
	if err := fetcher.Fetch(context, []string{releaseId}); err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	root := paths.NewPath().UnpackedDepDirectory(core.NewDependencyFromQualifiedReleaseId(parsed))
	err = os.Chdir(root)
	if err := context.LoadReleaseJson(); err != nil {
		return err
	}
	if err := d.Deploy(context, extraVars, extraProviders); err != nil {
		os.Chdir(currentDir)
		return err
	}
	return os.Chdir(currentDir)
}
