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

type StageState struct {
	UserInputs  map[string]interface{}      `json:"inputs"`
	Inputs      map[string]interface{}      `json:"calculated_inputs"`
	Outputs     map[string]interface{}      `json:"calculated_outputs"`
	Deployments map[string]*DeploymentState `json:"deployments"`
	Providers   map[string]string           `json:"providers"`
	Provides    []string                    `json:"provides"`
	Version     string                      `json:"version"`
	Step        string                      `json:"step"`
	Name        string                      `json:"-"`
}

func newStage() *StageState {
	return &StageState{
		UserInputs:  map[string]interface{}{},
		Inputs:      map[string]interface{}{},
		Outputs:     map[string]interface{}{},
		Providers:   map[string]string{},
		Provides:    []string{},
		Deployments: map[string]*DeploymentState{},
	}
}

func (st *StageState) validateAndFix(name string, envState *EnvironmentState, deplState *DeploymentState) error {
	st.Name = name
	if st.UserInputs == nil {
		st.UserInputs = map[string]interface{}{}
	}
	if st.Provides == nil {
		st.Provides = []string{}
	}
	if st.Inputs == nil {
		st.Inputs = map[string]interface{}{}
	}
	if st.Outputs == nil {
		st.Outputs = map[string]interface{}{}
	}
	if st.Providers == nil {
		st.Providers = map[string]string{}
	}
	if st.Deployments == nil {
		st.Deployments = map[string]*DeploymentState{}
	}
	for name, depl := range st.Deployments {
		depl.Name = name
		if err := depl.validateAndFixSubDeployment(st, envState, deplState); err != nil {
			return err
		}
	}
	return nil
}

func (st *StageState) setVersion(v string) *StageState {
	st.Version = v
	return st
}

func (st *StageState) setInputs(v map[string]interface{}) *StageState {
	st.Inputs = st.initIfNil(v)
	return st
}

func (st *StageState) setUserInputs(v map[string]interface{}) *StageState {
	st.UserInputs = st.initIfNil(v)
	return st
}

func (st *StageState) setOutputs(v map[string]interface{}) *StageState {
	st.Outputs = st.initIfNil(v)
	return st
}

func (st *StageState) initIfNil(v map[string]interface{}) map[string]interface{} {
	if v == nil {
		v = map[string]interface{}{}
	}
	return v
}
