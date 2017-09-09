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

package remote

import (
	"encoding/json"
	"fmt"

	"github.com/ankyra/escape-client/model/remote"
	. "github.com/ankyra/escape-core/state"
)

type remoteStateProvider struct {
	client    *remote.RegistryClient
	endpoints *remote.ServerEndpoints
}

func NewRemoteStateProvider(apiServer, escapeToken string, insecureSkipVerify bool) *remoteStateProvider {
	return &remoteStateProvider{
		client:    remote.NewRemoteClient(escapeToken, insecureSkipVerify),
		endpoints: remote.NewServerEndpoints(apiServer),
	}
}

func (r *remoteStateProvider) Load(project, env string) (*EnvironmentState, error) {
	url := r.endpoints.ProjectEnvironmentState(project, env)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Couldn't load environment state: %s", resp.Status)
	}
	result := EnvironmentState{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	prjState, _ := NewProjectState(project)
	prjState.Backend = r
	prjState.Environments[env] = &result
	result.Project = prjState
	return &result, prjState.ValidateAndFix()
}

func (l *remoteStateProvider) Save(depl *DeploymentState) error {
	return fmt.Errorf("Saving environment state not implemented")
}
