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
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
	"path/filepath"
)

type ProjectState struct {
	Name         string                       `json:"name"`
	Environments map[string]*EnvironmentState `json:"environments"`
	saveLocation string                       `json:"-"`
	provider     StateProvider                `json:"-"`
}

func newProjectState(prjName string) (*ProjectState, error) {
	return &ProjectState{
		Name:         prjName,
		Environments: map[string]*EnvironmentState{},
	}, nil
}

func NewProjectStateFromJsonString(data string, provider StateProvider) (*ProjectState, error) {
	prjState, err := newProjectState("")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), prjState); err != nil {
		return nil, err
	}
	if err := prjState.validateAndFix(); err != nil {
		return nil, err
	}
	prjState.provider = provider
	return prjState, nil
}

func NewProjectStateFromFile(prjName, cfgFile string, provider StateProvider) (*ProjectState, error) {
	if cfgFile == "" {
		return nil, fmt.Errorf("Configuration file path is required.")
	}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	if !util.PathExists(cfgFile) {
		p, err := newProjectState(prjName)
		if err != nil {
			return nil, err
		}
		p.saveLocation = cfgFile
		p.provider = provider
		return p, p.validateAndFix()
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	result, err := NewProjectStateFromJsonString(string(data), provider)
	if err != nil {
		return nil, err
	}
	result.saveLocation = cfgFile
	if result.Name == "" {
		result.Name = prjName
	}
	return result, nil
}

func (p *ProjectState) validateAndFix() error {
	if p.Name == "" {
		return fmt.Errorf("State is missing project name")
	}
	if p.Environments == nil {
		p.Environments = map[string]*EnvironmentState{}
	}
	for name, env := range p.Environments {
		if err := env.ValidateAndFix(name, p.Name, p.provider); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProjectState) GetName() string {
	return p.Name
}

func (p *ProjectState) SetName(name string) {
	p.Name = name
}

func (p *ProjectState) GetEnvironmentStateOrMakeNew(env string) *EnvironmentState {
	e, ok := p.Environments[env]
	if !ok || e == nil {
		p.Environments[env] = NewEnvironmentState(p.Name, env, p.provider)
		return p.Environments[env]
	}
	e.provider = p.provider
	return e
}

func (p *ProjectState) Save() error {
	if p.saveLocation == "" {
		return fmt.Errorf("Save location has not been set. Inexplicably")
	}
	contents := []byte(p.ToJson())
	return ioutil.WriteFile(p.saveLocation, contents, 0644)
}

func (p *ProjectState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}
