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
)

type StateController struct{}

func (p StateController) Show(context Context) error {
	fmt.Println(context.GetProjectState().ToJson())
	return nil
}

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

func (p StateController) CreateDeployment(context Context) error {

	str := []string{context.GetReleaseMetadata().GetVersionlessReleaseId()}
	deplState, err := context.GetEnvironmentState().GetDeploymentState(str)
	if err != nil {
		return err
	}
	inputs := *deplState.GetPreStepInputs("deploy")
	changed := false
	for _, i := range context.GetReleaseMetadata().GetInputs() {
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
		deplState.UpdateUserInputs("deploy", &inputs)
	}
	// TODO check and set providers
	fmt.Println(deplState.ToJson())
	return context.GetEnvironmentState().Save()
}
