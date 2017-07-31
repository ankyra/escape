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
)

type EnvironmentState struct {
	Name        string                      `json:"name"`
	Inputs      map[string]interface{}      `json:"inputs,omitempty"`
	Deployments map[string]*DeploymentState `json:"deployments,omitempty"`
	Project     *ProjectState               `json:"-"`
}

func NewEnvironmentState(envName string, project *ProjectState) *EnvironmentState {
	return &EnvironmentState{
		Name:        envName,
		Inputs:      map[string]interface{}{},
		Deployments: map[string]*DeploymentState{},
		Project:     project,
	}
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
	e.Name = name
	e.Project = project
	if e.Deployments == nil {
		e.Deployments = map[string]*DeploymentState{}
	}
	for deplName, depl := range e.Deployments {
		if err := depl.validateAndFix(deplName, e); err != nil {
			return err
		}
	}
	if e.Name == "" {
		return fmt.Errorf("Environment name is missing from the EnvironmentState")
	}
	return nil
}

func (e *EnvironmentState) LookupDeploymentState(deploymentName string) (*DeploymentState, error) {
	val, ok := e.Deployments[deploymentName]
	if !ok {
		return nil, fmt.Errorf("Deployment '%s' does not exist", deploymentName)
	}
	return val, nil
}

func (e *EnvironmentState) GetOrCreateDeploymentState(deploymentName string) *DeploymentState {
	depl, ok := e.Deployments[deploymentName]
	if !ok {
		depl = NewDeploymentState(e, deploymentName, deploymentName)
		e.Deployments[deploymentName] = depl
	}
	return depl
}

func (e *EnvironmentState) GetProviders() map[string][]string {
	result := map[string][]string{}
	for deplName, depl := range e.Deployments {
		st := depl.GetStageOrCreateNew("deploy")
		for _, provides := range st.Provides {
			result[provides] = append(result[provides], deplName)
		}
	}
	return result
}

func (e *EnvironmentState) GetProvidersOfType(typ string) []string {
	result := []string{}
	for deplName, depl := range e.Deployments {
		st := depl.GetStageOrCreateNew("deploy")
		for _, provides := range st.Provides {
			if provides == typ {
				result = append(result, deplName)
			}
		}
	}
	return result
}
