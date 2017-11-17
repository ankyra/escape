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
	"strings"

	. "github.com/ankyra/escape/model/interfaces"
)

type StateController struct{}

func (p StateController) ListDeployments(context Context) *ControllerResult {
	result := NewControllerResult()
	envState := context.GetEnvironmentState()

	deployments := []interface{}{}
	for _, depl := range envState.GetDeployments() {
		deployments = append(deployments, depl.GetName())
	}

	result.HumanOutput.AddList(deployments)
	result.MarshalableOutput = deployments

	return result
}

func (p StateController) ShowDeployment(context Context, dep string) error {
	envState := context.GetEnvironmentState()
	for _, depl := range envState.GetDeployments() {
		if depl.GetName() == dep {
			fmt.Println(depl.ToJson())
			return nil
		}
	}
	return fmt.Errorf("Deployment '%s' not found", dep)
}

func (p StateController) ShowProviders(context Context) *ControllerResult {
	result := NewControllerResult()
	envState := context.GetEnvironmentState()
	exists := false
	providers := map[string][]string{}
	for provider, implementations := range envState.GetProviders() {
		exists = true
		result.HumanOutput.AddLine("%s:", provider)
		result.HumanOutput.AddLine("\t%s", strings.Join(implementations, ", "))
		providers[provider] = implementations
	}
	if !exists {
		result.HumanOutput.AddLine("No providers found in the environment state. Try deploying one.")
		result.MarshalableOutput = providers
		return result
	}
	result.MarshalableOutput = providers
	return result
}

func (p StateController) CreateState(context Context, stage string, extraVars, extraProviders map[string]string) error {
	envState := context.GetEnvironmentState()
	metadata := context.GetReleaseMetadata()
	deplState := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	deplState.Release = metadata.GetVersionlessReleaseId()
	inputs := deplState.GetUserInputs(stage)
	changed := false
	for key, val := range extraVars {
		inputs[key] = val
		changed = true
	}
	for _, i := range metadata.GetInputs(stage) {
		val, ok := inputs[i.Id]
		if !ok {
			val = i.AskUserInput()
			if val != nil {
				changed = true
				inputs[i.Id] = val
			}
		}
	}
	if changed {
		deplState.UpdateUserInputs(stage, inputs)
	}
	if err := SetExtraProviders(context, stage, extraProviders); err != nil {
		return err
	}
	fmt.Println(deplState.ToJson())
	return deplState.Save()
}
