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
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type projectState struct {
	Name         string                       `json:"name"`
	Inputs       map[string]interface{}       `json:"inputs"`
	Environments map[string]*environmentState `json:"environments"`
	saveLocation string
}

func getDefaultName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.Name, nil
}

func newProjectState() (*projectState, error) {
	defaultName, err := getDefaultName()
	if err != nil {
		return nil, err
	}
	return &projectState{
		Name:         defaultName,
		Inputs:       map[string]interface{}{},
		Environments: map[string]*environmentState{},
	}, nil
}

func NewProjectStateFromJsonString(data string) (*projectState, error) {
	prjState, err := newProjectState()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), prjState); err != nil {
		return nil, err
	}
	if err := prjState.validateAndFix(); err != nil {
		return nil, err
	}
	return prjState, nil
}

func NewProjectStateFromFile(cfgFile string) (*projectState, error) {
	if cfgFile == "" {
		return nil, fmt.Errorf("Configuration file path is required.")
	}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	if !util.PathExists(cfgFile) {
		p, err := newProjectState()
		if err != nil {
			return nil, err
		}
		p.saveLocation = cfgFile
		return p, p.validateAndFix()
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	result, err := NewProjectStateFromJsonString(string(data))
	if err != nil {
		return nil, err
	}
	result.saveLocation = cfgFile
	return result, nil
}

func (p *projectState) getInputs() map[string]interface{} {
	return p.Inputs
}

func (p *projectState) validateAndFix() error {
	if p.Name == "" {
		defaultName, err := getDefaultName()
		if err != nil {
			return err
		}
		p.Name = defaultName
	}
	if p.Inputs == nil {
		p.Inputs = map[string]interface{}{}
	}
	if p.Environments == nil {
		p.Environments = map[string]*environmentState{}
	}
	for name, env := range p.Environments {
		if err := env.validateAndFix(name, p); err != nil {
			return err
		}
	}
	return nil
}

func (p *projectState) GetEnvironmentStateOrMakeNew(env string) EnvironmentState {
	e, ok := p.Environments[env]
	if !ok {
		e := NewEnvironmentState(p, env)
		p.Environments[env] = e.(*environmentState)
		return e
	}
	return e
}

func (p *projectState) Save() error {
	if p.saveLocation == "" {
		return fmt.Errorf("Save location has not been set. Inexplicably")
	}
	contents := []byte(p.ToJson())
	return ioutil.WriteFile(p.saveLocation, contents, 0644)
}

func (p *projectState) ToJson() string {
	str, err := json.MarshalIndent(p, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (p *projectState) GetName() string {
	return p.Name
}

func (p *projectState) SetName(name string) {
	p.Name = name
}
