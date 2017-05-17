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
	. "github.com/ankyra/escape-client/model/interfaces"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
)

type stage struct {
	UserInputs *map[string]interface{} `json:"inputs"`
	Inputs     *map[string]interface{} `json:"calculated_inputs"`
	Outputs    *map[string]interface{} `json:"calculated_outputs"`
	Version    string                  `json:"version"`
	Step       string                  `json:"step"`
}

type deploymentState struct {
	Name        string                       `json:"name"`
	Stages      map[string]*stage            `json:"stages"`
	Inputs      *map[string]interface{}      `json:"inputs"`
	Deployments *map[string]*deploymentState `json:"deployments"`
	Providers   *map[string]string           `json:"providers"`
	environment EnvironmentState             `json:"-"`
	parent      *deploymentState             `json:"-"`
}

func newStage() *stage {
	return &stage{
		UserInputs: &map[string]interface{}{},
		Inputs:     &map[string]interface{}{},
		Outputs:    &map[string]interface{}{},
	}
}

func NewDeploymentState(env EnvironmentState, name string) DeploymentState {
	return &deploymentState{
		Name:        name,
		Stages:      map[string]*stage{},
		Inputs:      &map[string]interface{}{},
		Providers:   &map[string]string{},
		Deployments: &map[string]*deploymentState{},
		environment: env,
	}
}

func (d *deploymentState) GetName() string {
	return d.Name
}
func (d *deploymentState) GetVersion(stage string) string {
	return d.getStage(stage).Version
}
func (d *deploymentState) SetVersion(stage, version string) error {
	st := d.getStage(stage)
	st.Version = version
	return nil
}
func (d *deploymentState) IsDeployed(stage, version string) bool {
	return d.getStage(stage).Version == version
}

func (d *deploymentState) GetEnvironmentState() EnvironmentState {
	return d.environment
}

func (d *deploymentState) NewDependencyDeploymentState(dep string) DeploymentState {
	depl := NewDeploymentState(d.environment, dep).(*deploymentState)
	depl.parent = d
	(*d.Deployments)[dep] = depl
	return depl
}

func (d *deploymentState) ValidateAndFix(name string, env EnvironmentState) error {
	d.Name = name
	d.environment = env
	if d.Name == "" {
		return fmt.Errorf("Deployment name is missing from DeploymentState")
	}
	if d.Inputs == nil {
		d.Inputs = &map[string]interface{}{}
	}
	if d.Providers == nil {
		d.Providers = &map[string]string{}
	}
	if d.Deployments == nil {
		d.Deployments = &map[string]*deploymentState{}
	}
	for name, depl := range *d.Deployments {
		depl.Name = name
		if err := depl.ValidateAndFixSubDeployment(env, d); err != nil {
			return err
		}
	}
	if d.Stages == nil {
		d.Stages = map[string]*stage{}
	}
	for _, st := range d.Stages {
		if st.UserInputs == nil {
			st.UserInputs = &map[string]interface{}{}
		}
		if st.Inputs == nil {
			st.Inputs = &map[string]interface{}{}
		}
		if st.Outputs == nil {
			st.Outputs = &map[string]interface{}{}
		}
	}
	return nil
}

func (d *deploymentState) ValidateAndFixSubDeployment(env EnvironmentState, parent *deploymentState) error {
	d.parent = parent
	return d.ValidateAndFix(d.Name, env)
}

func (d *deploymentState) getStage(stage string) *stage {
	st, ok := d.Stages[stage]
	if !ok {
		st = newStage()
		d.Stages[stage] = st
	}
	return st
}
func (d *deploymentState) GetUserInputs(stage string) *map[string]interface{} {
	return d.getStage(stage).UserInputs
}
func (d *deploymentState) GetCalculatedInputs(stage string) *map[string]interface{} {
	return d.getStage(stage).Inputs
}
func (d *deploymentState) GetCalculatedOutputs(stage string) *map[string]interface{} {
	return d.getStage(stage).Outputs
}

func (d *deploymentState) UpdateInputs(stage string, inputs *map[string]interface{}) error {
	st := d.getStage(stage)
	st.Inputs = inputs
	return d.Save()
}
func (d *deploymentState) UpdateUserInputs(stage string, inputs *map[string]interface{}) error {
	st := d.getStage(stage)
	st.UserInputs = inputs
	return d.Save()
}
func (d *deploymentState) UpdateOutputs(stage string, outputs *map[string]interface{}) error {
	st := d.getStage(stage)
	st.Outputs = outputs
	return d.Save()
}

func (d *deploymentState) GetPreStepInputs(stage string) *map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range d.environment.GetProjectState().GetInputs() {
		result[key] = val
	}
	for key, val := range d.environment.GetInputs() {
		result[key] = val
	}
	deps := []*deploymentState{d}
	p := d.parent
	for p != nil {
		deps = append(deps, p)
		p = p.parent
	}
	for i := len(deps) - 1; i >= 0; i-- {
		p = deps[i]
		if p.Inputs != nil {
			for key, val := range *p.Inputs {
				result[key] = val
			}
		}
		st := p.getStage(stage)
		if st.UserInputs != nil {
			for key, val := range *st.UserInputs {
				result[key] = val
			}
		}
		p = p.parent
	}
	return &result
}

func (d *deploymentState) Save() error {
	if d.environment.IsRemote() {
		//            self._update_remote_config(release_inputs, release_outputs)
		return nil
	} else {
		return d.environment.Save()
	}
}

func (d *deploymentState) GetProviders() map[string]string {
	result := map[string]string{}
	if d.Providers == nil {
		return result
	}
	for key, val := range *d.Providers {
		result[key] = val
	}
	current := d
	for current.parent != nil {
		current = current.parent
		for key, val := range *current.Providers {
			_, alreadySet := result[key]
			if !alreadySet {
				result[key] = val
			}
		}
	}
	return result
}

func (d *deploymentState) ResolveConsumer(consumer string) (DeploymentState, error) {

	providers := d.GetProviders()
	val, ok := providers[consumer]
	if !ok {
		return nil, fmt.Errorf("No provider of type '%s' was configured in the deployment state.", consumer)
	}
	return d.environment.LookupDeploymentState(val)
}

func (p *deploymentState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (d *deploymentState) ToScriptEnvironment(metadataMap map[string]*core.ReleaseMetadata, stage string) (*script.ScriptEnvironment, error) {
	result := map[string]script.Script{}
	result["this"] = d.ToScript(metadataMap["this"], stage)
	for key, deplName := range d.GetProviders() {
		deplState, err := d.environment.LookupDeploymentState(deplName)
		if err != nil {
			return nil, err
		}
		result[key] = deplState.ToScript(metadataMap[key], "deploy")
	}
	for key, deplState := range *d.Deployments {
		metadata := metadataMap[key]
		version := metadata.GetVersion()
		if deplState.IsDeployed("deploy", version) {
			result[key+"-v"+version] = deplState.ToScript(metadataMap[key], "deploy")
		}
	}
	for key, metadata := range metadataMap {
		reference, exists := result[metadata.GetReleaseId()]
		if exists {
			result[key] = reference
		}
	}
	return script.NewScriptEnvironmentWithGlobals(result), nil

}

func (d *deploymentState) ToScript(metadata *core.ReleaseMetadata, stage string) script.Script {
	result := map[string]script.Script{}
	if metadata != nil {
		result = metadata.ToScriptMap()
	}
	result["inputs"] = script.LiftDict(d.liftScriptValues(d.GetCalculatedInputs(stage)))
	result["outputs"] = script.LiftDict(d.liftScriptValues(d.GetCalculatedOutputs(stage)))
	env := d.GetEnvironmentState()
	prj := env.GetProjectState()
	result["project"] = script.LiftString(prj.GetName())
	result["environment"] = script.LiftString(env.GetName())
	result["deployment"] = script.LiftString(d.GetName())
	return script.LiftDict(result)
}

func (d *deploymentState) liftScriptValues(values *map[string]interface{}) map[string]script.Script {
	result := map[string]script.Script{}
	if values != nil {
		for key, val := range *values {
			v, err := script.Lift(val)
			if err != nil {
				panic(err)
			}
			result[key] = v
		}
	}
	return result
}
