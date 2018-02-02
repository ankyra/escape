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

package local

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	. "github.com/ankyra/escape-core/state"
)

type localStateProvider struct {
	state        *ProjectState
	saveLocation string
}

func NewLocalStateProvider(file string) *localStateProvider {
	return &localStateProvider{
		saveLocation: file,
	}
}

func (l *localStateProvider) Load(project, env string) (*EnvironmentState, error) {
	var err error
	if project == "" {
		project, err = getDefaultName()
		if err != nil {
			return nil, err
		}
	}
	prj, err := NewProjectStateFromFile(project, l.saveLocation, l)
	if err != nil {
		return nil, err
	}
	l.state = prj
	path, err := filepath.Abs(l.saveLocation)
	if err != nil {
		return nil, err
	}
	l.saveLocation = path
	return prj.GetEnvironmentStateOrMakeNew(env)
}

func (l *localStateProvider) Save(depl *DeploymentState) error {
	if l.saveLocation == "" {
		return fmt.Errorf("Save location has not been set. Inexplicably")
	}
	contents := []byte(l.state.ToJson())
	return ioutil.WriteFile(l.saveLocation, contents, 0644)
}

func getDefaultName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.Username, nil
}
