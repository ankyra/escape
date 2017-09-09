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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ankyra/escape-core"
)

type DeploymentState struct {
	Name        string                 `json:"name"`
	Release     string                 `json:"release,omitempty"`
	Stages      map[string]*StageState `json:"stages,omitempty"`
	Inputs      map[string]interface{} `json:"inputs,omitempty"`
	environment *EnvironmentState      `json:"-"`
	parent      *DeploymentState       `json:"-"`
	parentStage *StageState            `json:"-"`
}

func NewDeploymentState(env *EnvironmentState, name, release string) *DeploymentState {
	return &DeploymentState{
		Name:        name,
		Release:     release,
		Stages:      map[string]*StageState{},
		Inputs:      map[string]interface{}{},
		environment: env,
	}
}

func (d *DeploymentState) GetName() string {
	return d.Name
}

func (d *DeploymentState) GetRootDeploymentName() string {
	prev := d
	p := prev
	for p != nil {
		prev = p
		p = p.parent
	}
	return prev.Name
}

func (d *DeploymentState) GetDependencyPath() string {
	result := []string{}
	p := d
	for p != nil {
		result = append(result, p.Name)
		p = p.parent
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return strings.Join(result, ":")
}

func (d *DeploymentState) GetReleaseId(stage string) string {
	return d.Release + "-v" + d.GetVersion(stage)
}

func (d *DeploymentState) GetVersion(stage string) string {
	return d.GetStageOrCreateNew(stage).Version
}

func (d *DeploymentState) GetEnvironmentState() *EnvironmentState {
	return d.environment
}

func (d *DeploymentState) GetDeploymentOrMakeNew(stage, deploymentName string) *DeploymentState {
	st := d.GetStageOrCreateNew(stage)
	depl, ok := st.Deployments[deploymentName]
	if !ok {
		depl = NewDeploymentState(d.environment, deploymentName, deploymentName)
		depl.parent = d
	}
	depl.parentStage = st
	st.Deployments[deploymentName] = depl
	return depl
}

func (d *DeploymentState) GetUserInputs(stage string) map[string]interface{} {
	return d.GetStageOrCreateNew(stage).UserInputs
}

func (d *DeploymentState) GetCalculatedInputs(stage string) map[string]interface{} {
	return d.GetStageOrCreateNew(stage).Inputs
}

func (d *DeploymentState) GetCalculatedOutputs(stage string) map[string]interface{} {
	return d.GetStageOrCreateNew(stage).Outputs
}

func (d *DeploymentState) UpdateInputs(stage string, inputs map[string]interface{}) error {
	d.GetStageOrCreateNew(stage).SetInputs(inputs)
	return d.Save()
}

func (d *DeploymentState) UpdateUserInputs(stage string, inputs map[string]interface{}) error {
	d.GetStageOrCreateNew(stage).SetUserInputs(inputs)
	return d.Save()
}

func (d *DeploymentState) UpdateOutputs(stage string, outputs map[string]interface{}) error {
	d.GetStageOrCreateNew(stage).SetOutputs(outputs)
	return d.Save()
}

func (d *DeploymentState) CommitVersion(stage string, metadata *core.ReleaseMetadata) error {
	d.GetStageOrCreateNew(stage).SetVersion(metadata.Version)
	d.GetStageOrCreateNew(stage).Provides = metadata.GetProvides()
	return nil
}

func (d *DeploymentState) UpdateStatus(stage string, status *Status) error {
	d.GetStageOrCreateNew(stage).Status = status
	return d.Save()
}
func (d *DeploymentState) GetStatus(stage string) *Status {
	return d.GetStageOrCreateNew(stage).Status
}

func (d *DeploymentState) IsDeployed(stage string, metadata *core.ReleaseMetadata) bool {
	return d.GetStageOrCreateNew(stage).Version == metadata.Version
}

func (d *DeploymentState) Save() error {
	return d.environment.Save(d)
}

func (p *DeploymentState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (d *DeploymentState) SetProvider(stage, name, deplName string) {
	d.GetStageOrCreateNew(stage).Providers[name] = deplName
}

func (d *DeploymentState) GetProviders(stage string) map[string]string {
	result := map[string]string{}
	d.walkStatesAndStages(stage, func(p *DeploymentState, st *StageState) {
		for key, val := range st.Providers {
			result[key] = val
		}
	})
	return result
}

func (d *DeploymentState) GetPreStepInputs(stage string) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range d.environment.Inputs {
		result[key] = val
	}
	d.walkStatesAndStages(stage, func(p *DeploymentState, st *StageState) {
		if p.Inputs != nil {
			for key, val := range p.Inputs {
				result[key] = val
			}
		}
		if st.UserInputs != nil {
			for key, val := range st.UserInputs {
				result[key] = val
			}
		}
	})
	return result
}

func (d *DeploymentState) walkStatesAndStages(startStage string, cb func(*DeploymentState, *StageState)) {
	deps := d.getDependencyStates()
	stages := d.getDependencyStages(startStage)
	for i := len(deps) - 1; i >= 0; i-- {
		p := deps[i]
		stage := stages[i]
		cb(p, stage)
	}
}

func (d *DeploymentState) getDependencyStates() []*DeploymentState {
	deps := []*DeploymentState{}
	p := d
	for p != nil {
		deps = append(deps, p)
		p = p.parent
	}
	return deps
}

func (d *DeploymentState) getDependencyStages(startStage string) []*StageState {
	stages := []*StageState{d.GetStageOrCreateNew(startStage)}
	p := d
	for p != nil {
		stages = append(stages, p.parentStage)
		p = p.parent
	}
	return stages
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
	if d.Stages == nil {
		d.Stages = map[string]*StageState{}
	}
	for name, st := range d.Stages {
		st.validateAndFix(name, env, d)
	}
	return nil
}

func (d *DeploymentState) validateAndFixSubDeployment(stage *StageState, env *EnvironmentState, parent *DeploymentState) error {
	d.parent = parent
	d.parentStage = stage
	return d.validateAndFix(d.Name, env)
}

func (d *DeploymentState) GetStageOrCreateNew(stage string) *StageState {
	st, ok := d.Stages[stage]
	if !ok || st == nil {
		st = newStage()
		d.Stages[stage] = st
	}
	st.validateAndFix(stage, d.environment, d)
	return st
}
