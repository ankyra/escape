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

package types

import (
	"fmt"
)

type EnvironmentState struct {
	Name        string                      `json:"name"`
	Inputs      map[string]interface{}      `json:"inputs"`
	Deployments map[string]*DeploymentState `json:"deployments"`
	ProjectName string                      `json:"-"`
	provider    StateProvider               `json:"-"`
}

func NewEnvironmentState(prjName, envName string, provider StateProvider) *EnvironmentState {
	return &EnvironmentState{
		ProjectName: prjName,
		Name:        envName,
		Inputs:      map[string]interface{}{},
		Deployments: map[string]*DeploymentState{},
		provider:    provider,
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
	return e.ProjectName
}
func (e *EnvironmentState) getInputs() map[string]interface{} {
	return e.Inputs
}
func (e *EnvironmentState) GetName() string {
	return e.Name
}
func (e *EnvironmentState) Save(d *DeploymentState) error {
	if e.provider == nil {
		return fmt.Errorf("No state provider configured. This is a bug in Escape.")
	}
	return e.provider.Save(d)
}

func (e *EnvironmentState) ValidateAndFix(name, prjName string, provider StateProvider) error {
	e.Name = name
	e.provider = provider
	e.ProjectName = prjName
	if e.Deployments == nil {
		e.Deployments = map[string]*DeploymentState{}
	}
	for deplName, depl := range e.Deployments {
		if err := depl.validateAndFix(deplName, e); err != nil {
			return err
		}
	}
	if e.ProjectName == "" {
		return fmt.Errorf("EnvironmentState's ProjectState reference has not been set")
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
