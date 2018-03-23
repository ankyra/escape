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

package state

import (
	"fmt"
	"strings"

	"github.com/ankyra/escape-core/state/validate"
)

func DeploymentDoesNotExistError(deploymentName string) error {
	return fmt.Errorf("Deployment '%s' does not exist", deploymentName)
}

func DeploymentPathResolveError(stage, deploymentPath, deploymentName string) error {
	return fmt.Errorf("Failed to resolve deployment path '%s': the deployment '%s' could not be found in the %s stage",
		deploymentPath, deploymentName, stage)
}

type EnvironmentState struct {
	Name        string                      `json:"name"`
	Inputs      map[string]interface{}      `json:"inputs,omitempty"`
	Deployments map[string]*DeploymentState `json:"deployments,omitempty"`
	Project     *ProjectState               `json:"-"`
}

func NewEnvironmentState(envName string, project *ProjectState) (*EnvironmentState, error) {
	e := &EnvironmentState{
		Name:        envName,
		Inputs:      map[string]interface{}{},
		Deployments: map[string]*DeploymentState{},
		Project:     project,
	}
	return e, e.ValidateAndFix(envName, project)
}

func (e *EnvironmentState) GetDeployments() []*DeploymentState {
	result := []*DeploymentState{}
	for _, d := range e.Deployments {
		result = append(result, d)
	}
	return result
}

func (e *EnvironmentState) GetProjectName() string {
	return e.Project.Name
}

func (e *EnvironmentState) Save(d *DeploymentState) error {
	return e.Project.Save(d)
}

func (e *EnvironmentState) ValidateAndFix(name string, project *ProjectState) error {
	if !validate.IsValidEnvironmentName(name) {
		return validate.InvalidEnvironmentNameError(name)
	}
	e.Name = name
	e.Project = project
	if e.Inputs == nil {
		e.Inputs = map[string]interface{}{}
	}
	if e.Deployments == nil {
		e.Deployments = map[string]*DeploymentState{}
	}
	for deplName, depl := range e.Deployments {
		if err := depl.validateAndFix(deplName, e); err != nil {
			return err
		}
	}
	return nil
}

func (e *EnvironmentState) LookupDeploymentState(deploymentName string) (*DeploymentState, error) {
	val, ok := e.Deployments[deploymentName]
	if !ok {
		return nil, DeploymentDoesNotExistError(deploymentName)
	}
	return val, nil
}

func (e *EnvironmentState) DeleteDeployment(deploymentName string) error {
	_, ok := e.Deployments[deploymentName]
	if !ok {
		return DeploymentDoesNotExistError(deploymentName)
	}
	delete(e.Deployments, deploymentName)
	return e.Project.CommitDeleteDeployment(e.Name, deploymentName)
}

func (e *EnvironmentState) ResolveDeploymentPath(stage, deploymentPath string) (*DeploymentState, error) {
	parts := strings.Split(deploymentPath, ":")
	if len(parts) == 0 {
		return nil, DeploymentDoesNotExistError(deploymentPath)
	}
	deploymentName := parts[0]
	val, ok := e.Deployments[deploymentName]
	if !ok {
		return nil, DeploymentDoesNotExistError(deploymentPath)
	}
	for _, p := range parts[1:] {
		newVal, err := val.GetDeployment(stage, p)
		if err != nil {
			return nil, DeploymentPathResolveError(stage, deploymentPath, p)
		}
		val = newVal
		stage = DeployStage
	}
	return val, nil
}

func (e *EnvironmentState) GetOrCreateDeploymentState(deploymentName string) (*DeploymentState, error) {
	depl, ok := e.Deployments[deploymentName]
	if !ok {
		depl, err := NewDeploymentState(e, deploymentName, deploymentName)
		if err != nil {
			return nil, err
		}
		e.Deployments[deploymentName] = depl
		return depl, nil
	}
	return depl, nil
}

func (e *EnvironmentState) GetProviders() map[string][]string {
	result := map[string][]string{}
	for deplName, depl := range e.Deployments {
		st := depl.GetStageOrCreateNew(DeployStage)
		for _, provides := range st.Provides {
			result[provides] = append(result[provides], deplName)
		}
	}
	return result
}

func (e *EnvironmentState) GetProvidersOfType(typ string) []string {
	result := []string{}
	for deplName, depl := range e.Deployments {
		st := depl.GetStageOrCreateNew(DeployStage)
		for _, provides := range st.Provides {
			if provides == typ {
				result = append(result, deplName)
			}
		}
	}
	return result
}
