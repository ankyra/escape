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
	. "github.com/ankyra/escape-client/model/interfaces"
	"strings"
)

type StateController struct{}

func (p StateController) ShowDeployments(context Context) error {
	envState := context.GetEnvironmentState()
	for _, depl := range envState.GetDeployments() {
		fmt.Println(depl.GetName())
	}
	return nil
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

func (p StateController) ShowProviders(context Context) error {
	envState := context.GetEnvironmentState()
	exists := false
	for provider, implementations := range envState.GetProviders() {
		exists = true
		fmt.Printf("%s:\n", provider)
		fmt.Printf("\t%s\n", strings.Join(implementations, ", "))
	}
	if !exists {
		fmt.Println("No providers found in the environment state. Try deploying one.")
	}
	return nil
}

func (p StateController) CreateState(context Context, stage string) error {
	envState := context.GetEnvironmentState()
	metadata := context.GetReleaseMetadata()
	deplState := envState.GetOrCreateDeploymentState(context.GetRootDeploymentName())
	deplState.Release = metadata.GetVersionlessReleaseId()
	inputs := deplState.GetPreStepInputs(stage)
	changed := false
	for _, i := range metadata.GetInputs() {
		val, ok := inputs[i.GetId()]
		if !ok {
			val = i.AskUserInput()
			if val != nil {
				changed = true
				inputs[i.GetId()] = val
			}
		}
	}
	if changed {
		deplState.UpdateUserInputs(stage, inputs)
	}
	providers := envState.GetProviders()
	for _, c := range metadata.GetConsumes() {
		implementations := providers[c]
		if len(implementations) == 1 {
			deplState.SetProvider(stage, c, implementations[0])
		}
	}
	fmt.Println(deplState.ToJson())
	return deplState.Save()
}
