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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ankyra/escape-core/util"
)

type Backend interface {
	Save(d *DeploymentState) error
}

type ProjectState struct {
	Name         string                       `json:"name"`
	Environments map[string]*EnvironmentState `json:"environments,omitempty"`
	Backend      Backend                      `json:"-"`
}

func NewProjectState(prjName string) (*ProjectState, error) {
	return &ProjectState{
		Name:         prjName,
		Environments: map[string]*EnvironmentState{},
	}, nil
}

func NewProjectStateFromJsonString(data string, backend Backend) (*ProjectState, error) {
	prjState, err := NewProjectState("")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), prjState); err != nil {
		return nil, err
	}
	if err := prjState.ValidateAndFix(); err != nil {
		return nil, err
	}
	prjState.Backend = backend
	return prjState, nil
}

func NewProjectStateFromFile(prjName, cfgFile string, backend Backend) (*ProjectState, error) {
	if cfgFile == "" {
		return nil, fmt.Errorf("Configuration file path is required.")
	}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	if !util.PathExists(cfgFile) {
		p, err := NewProjectState(prjName)
		if err != nil {
			return nil, err
		}
		p.Backend = backend
		return p, p.ValidateAndFix()
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	result, err := NewProjectStateFromJsonString(string(data), backend)
	if err != nil {
		return nil, err
	}
	if result.Name == "" {
		result.Name = prjName
	}
	return result, nil
}

func (p *ProjectState) Save(d *DeploymentState) error {
	return p.Backend.Save(d)
}

func (p *ProjectState) ValidateAndFix() error {
	if p.Name == "" {
		return fmt.Errorf("State is missing project name")
	}
	if p.Environments == nil {
		p.Environments = map[string]*EnvironmentState{}
	}
	for name, env := range p.Environments {
		if err := env.ValidateAndFix(name, p); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProjectState) GetEnvironmentStateOrMakeNew(env string) *EnvironmentState {
	e, ok := p.Environments[env]
	if !ok || e == nil {
		p.Environments[env] = NewEnvironmentState(env, p)
		return p.Environments[env]
	}
	return e
}

func (p *ProjectState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}
