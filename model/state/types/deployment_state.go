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
	"encoding/json"
	"fmt"
)

type DeploymentState struct {
	Name        string                      `json:"name"`
	Release     string                      `json:"release"`
	Stages      map[string]*stage           `json:"stages"`
	Inputs      map[string]interface{}      `json:"inputs"`
	Deployments map[string]*DeploymentState `json:"deployments"`
	Providers   map[string]string           `json:"providers"`
	environment *EnvironmentState           `json:"-"`
	parent      *DeploymentState            `json:"-"`
}

func NewDeploymentState(env *EnvironmentState, name, release string) *DeploymentState {
	return &DeploymentState{
		Name:        name,
		Release:     release,
		Stages:      map[string]*stage{},
		Inputs:      map[string]interface{}{},
		Providers:   map[string]string{},
		Deployments: map[string]*DeploymentState{},
		environment: env,
	}
}
func (d *DeploymentState) NewDependencyDeploymentState(dep string) *DeploymentState {
	depl := NewDeploymentState(d.environment, dep, dep)
	depl.parent = d
	d.Deployments[dep] = depl
	return depl
}

func (d *DeploymentState) GetName() string {
	return d.Name
}

func (d *DeploymentState) GetRelease() string {
	return d.Release
}

func (d *DeploymentState) GetReleaseId(stage string) string {
	return d.GetRelease() + "-v" + d.GetVersion(stage)
}

func (d *DeploymentState) GetVersion(stage string) string {
	return d.getStage(stage).Version
}

func (d *DeploymentState) GetEnvironmentState() *EnvironmentState {
	return d.environment
}

func (d *DeploymentState) GetDeployments() []*DeploymentState {
	result := []*DeploymentState{}
	for _, val := range d.Deployments {
		result = append(result, val)
	}
	return result
}
func (d *DeploymentState) GetDeployment(deploymentName string) *DeploymentState {
	for _, val := range d.Deployments {
		if val.GetName() == deploymentName {
			return val
		}
	}
	return nil
}

func (d *DeploymentState) GetUserInputs(stage string) map[string]interface{} {
	return d.getStage(stage).UserInputs
}

func (d *DeploymentState) GetCalculatedInputs(stage string) map[string]interface{} {
	return d.getStage(stage).Inputs
}

func (d *DeploymentState) GetCalculatedOutputs(stage string) map[string]interface{} {
	return d.getStage(stage).Outputs
}

func (d *DeploymentState) UpdateInputs(stage string, inputs map[string]interface{}) error {
	d.getStage(stage).setInputs(inputs)
	return d.Save()
}

func (d *DeploymentState) UpdateUserInputs(stage string, inputs map[string]interface{}) error {
	d.getStage(stage).setUserInputs(inputs)
	return d.Save()
}

func (d *DeploymentState) UpdateOutputs(stage string, outputs map[string]interface{}) error {
	d.getStage(stage).setOutputs(outputs)
	return d.Save()
}

func (d *DeploymentState) SetVersion(stage, version string) error {
	d.getStage(stage).setVersion(version)
	return nil
}

func (d *DeploymentState) IsDeployed(stage, version string) bool {
	return d.getStage(stage).Version == version
}

func (d *DeploymentState) Save() error {
	return d.GetEnvironmentState().Save(d)
}

func (p *DeploymentState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (d *DeploymentState) GetProviders() map[string]string {
	result := map[string]string{}
	for key, val := range d.Providers {
		result[key] = val
	}
	current := d
	for current.parent != nil {
		current = current.parent
		for key, val := range current.Providers {
			if _, alreadySet := result[key]; !alreadySet {
				result[key] = val
			}
		}
	}
	return result
}

func (d *DeploymentState) GetPreStepInputs(stage string) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range d.environment.getInputs() {
		result[key] = val
	}
	deps := []*DeploymentState{d}
	p := d.parent
	for p != nil {
		deps = append(deps, p)
		p = p.parent
	}
	for i := len(deps) - 1; i >= 0; i-- {
		p = deps[i]
		if p.Inputs != nil {
			for key, val := range p.Inputs {
				result[key] = val
			}
		}
		st := p.getStage(stage)
		if st.UserInputs != nil {
			for key, val := range st.UserInputs {
				result[key] = val
			}
		}
		p = p.parent
	}
	return result
}

func (d *DeploymentState) validateAndFix(name string, env *EnvironmentState) error {
	d.Name = name
	d.environment = env
	if d.Name == "" {
		return fmt.Errorf("Deployment name is missing from DeploymentState")
	}
	if d.Release == "" {
		d.Release = name
	}
	if d.Inputs == nil {
		d.Inputs = map[string]interface{}{}
	}
	if d.Providers == nil {
		d.Providers = map[string]string{}
	}
	if d.Deployments == nil {
		d.Deployments = map[string]*DeploymentState{}
	}
	for name, depl := range d.Deployments {
		depl.Name = name
		if err := depl.validateAndFixSubDeployment(env, d); err != nil {
			return err
		}
	}
	if d.Stages == nil {
		d.Stages = map[string]*stage{}
	}
	for _, st := range d.Stages {
		st.validateAndFix()
	}
	return nil
}

func (d *DeploymentState) validateAndFixSubDeployment(env *EnvironmentState, parent *DeploymentState) error {
	d.parent = parent
	return d.validateAndFix(d.Name, env)
}

func (d *DeploymentState) getStage(stage string) *stage {
	st, ok := d.Stages[stage]
	if !ok || st == nil {
		st = newStage()
		d.Stages[stage] = st
	}
	st.validateAndFix()
	return st
}
