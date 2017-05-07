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

package state

import (
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
)

type environmentState struct {
	Name         string                      `json:"name"`
	Inputs       map[string]interface{}      `json:"inputs"`
	Deployments  map[string]*deploymentState `json:"deployments"`
	projectState ProjectState                `json:"-"`
}

func NewEnvironmentState(prj ProjectState, envName string) EnvironmentState {
	return &environmentState{
		projectState: prj,
		Name:         envName,
		Inputs:       map[string]interface{}{},
		Deployments:  map[string]*deploymentState{},
	}
}

func (e *environmentState) GetDeployments() []DeploymentState {
	result := []DeploymentState{}
	for _, d := range e.Deployments {
		result = append(result, d)
	}
	return result
}

func (e *environmentState) GetProjectState() ProjectState {
	return e.projectState
}
func (e *environmentState) GetInputs() map[string]interface{} {
	return e.Inputs
}

func (e *environmentState) GetName() string {
	return e.Name
}

func (e *environmentState) IsRemote() bool {
	return e.projectState.IsRemote()
}

func (e *environmentState) Save() error {
	return e.projectState.Save()
}

func (e *environmentState) ValidateAndFix(name string, p ProjectState) error {
	e.Name = name
	e.projectState = p
	if e.Deployments == nil {
		e.Deployments = map[string]*deploymentState{}
	}
	for deplName, depl := range e.Deployments {
		if err := depl.ValidateAndFix(deplName, e); err != nil {
			return err
		}
	}
	if e.projectState == nil {
		return fmt.Errorf("EnvironmentState's ProjectState reference has not been set")
	}
	if e.Name == "" {
		return fmt.Errorf("Environment name is missing from the EnvironmentState")
	}
	return nil
}

func (e *environmentState) LookupDeploymentState(deploymentName string) (DeploymentState, error) {
	val, ok := e.Deployments[deploymentName]
	if !ok {
		return nil, fmt.Errorf("Deployment '%s' does not exist", deploymentName)
	}
	return val, nil
}

func (e *environmentState) GetDeploymentState(deps []string) (DeploymentState, error) {
	if deps == nil || len(deps) == 0 {
		return nil, fmt.Errorf("Missing name to resolve deployment state. This is a bug in Escape.")
	}
	if len(deps) == 1 {
		return e.getDeploymentState(deps[0])
	} else {
		return e.getDeploymentStateForDependency(deps)
	}
}

func (e *environmentState) getDeploymentState(versionlessReleaseId string) (DeploymentState, error) {
	deploymentName := versionlessReleaseId
	depl, ok := e.Deployments[deploymentName]
	if !ok {
		depl = NewDeploymentState(e, deploymentName).(*deploymentState)
		e.Deployments[deploymentName] = depl
	}
	return depl, nil
}

func (e *environmentState) getDeploymentStateForDependency(deps []string) (DeploymentState, error) {
	deploymentName := deps[0]
	result := e.getOrCreateRootDeploymentState(deploymentName)
	for _, dep := range deps[1:] {
		depl, ok := (*result.Deployments)[dep]
		if !ok {
			result = result.NewDependencyDeploymentState(dep).(*deploymentState)
		} else {
			result = depl
		}
	}
	return result, nil
}

func (e *environmentState) getOrCreateRootDeploymentState(deploymentName string) *deploymentState {
	depl, ok := e.Deployments[deploymentName]
	if !ok {
		depl = NewDeploymentState(e, deploymentName).(*deploymentState)
		e.Deployments[deploymentName] = depl
	}
	return depl
}
