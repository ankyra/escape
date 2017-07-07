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
	"github.com/ankyra/escape-client/model/registry/types"
	core "github.com/ankyra/escape-core"
	"io"
	"io/ioutil"
	"os"
)

type registry struct {
	client    *RegistryClient
	endpoints *ServerEndpoints
}

func NewRemoteRegistry(apiServer, escapeToken string) *registry {
	return &registry{
		client:    NewRegistryClient(escapeToken),
		endpoints: NewServerEndpoints(apiServer),
	}
}

func (r *registry) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	url := r.endpoints.ReleaseQuery(project, name, version)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		releaseQuery := project + "/" + name + "-" + version
		if project == "_" {
			releaseQuery = name + "-" + version
		}
		return nil, fmt.Errorf("Couldn't query release '%s': %s", releaseQuery, resp.Status)
	}
	result := core.NewEmptyReleaseMetadata()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *registry) QueryPreviousReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	url := r.endpoints.PreviousReleaseQuery(project, name, version)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("Unauthorized")
	} else if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		releaseQuery := project + "/" + name + "-" + version
		if project == "_" {
			releaseQuery = name + "-" + version
		}
		return nil, fmt.Errorf("Couldn't query previous release '%s': %s", releaseQuery, resp.Status)
	}
	result := core.NewEmptyReleaseMetadata()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *registry) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	url := r.endpoints.NextReleaseVersion(project, name, versionPrefix)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 401 {
		return "", fmt.Errorf("Unauthorized")
	} else if resp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("There was a problem with the query: %s", body)
	} else if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Could not query next release version.")
		}
		return "", fmt.Errorf("Could not query next release version: %s", body)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (r *registry) DownloadRelease(project, name, version, targetFile string) error {
	url := r.endpoints.DownloadRelease(project, name, version)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		releaseId := project + "/" + name + "-v" + version
		if project == "_" {
			releaseId = name + "-v" + version
		}
		return fmt.Errorf("Couldn't download release '%s': %s", releaseId, resp.Status)
	}
	fmt.Println("Writing: " + targetFile)
	fp, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, resp.Body); err != nil {
		return err
	}
	return nil
}

func (r *registry) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	if err := r.register(project, metadata); err != nil {
		return err
	}
	url := r.endpoints.UploadRelease(project, metadata.Name, metadata.Version)
	resp, err := r.client.POST_file_with_authentication(url, releasePath)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Couldn't upload package (%s)", resp.Status)
		}
		return fmt.Errorf("Couldn't upload package (%s): %s", resp.Status, result)
	}
	return nil
}

func (r *registry) register(project string, metadata *core.ReleaseMetadata) error {
	url := r.endpoints.RegisterPackage(project)
	resp, err := r.client.POST_json_with_authentication(url, metadata)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Couldn't register package: %s", resp.Status)
		}
		return fmt.Errorf("Couldn't register package (%s): %s", resp.Status, result)
	}
	return nil
}

func (r *registry) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	url = r.endpoints.AuthMethods(url)
	resp, err := r.client.GET(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get auth methods from server '%s': %s", url, resp.Status)
		}
		return nil, fmt.Errorf("Couldn't get auth methods from server '%s': %s", url, resp.Status, result)
	}
	result := map[string]*types.AuthMethod{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *registry) LoginWithSecretToken(url, username, password string) (string, error) {
	payload := map[string]string{
		"username":     username,
		"secret_token": password,
	}
	resp, err := r.client.POST_json_with_authentication(url, payload)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 401 {
		return "", fmt.Errorf("Invalid credentials")
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Failed to login: %s", resp.Status)
		}
		return "", fmt.Errorf("Failed to login (%s): %s", resp.Status, result)
	}
	return resp.Header.Get("X-Escape-Token"), nil
}

func (r *registry) urlToList(url, doingMessage string) ([]string, error) {
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("Unauthorized")
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Couldn't %s: %s", doingMessage, resp.Status)
	}
	result := []string{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *registry) ListProjects() ([]string, error) {
	return r.urlToList(r.endpoints.ListProjects(), "list projects")
}
func (r *registry) ListApplications(project string) ([]string, error) {
	return r.urlToList(r.endpoints.ProjectQuery(project), "list applications for project '"+project+"'")
}
func (r *registry) ListVersions(project, app string) ([]string, error) {
	return r.urlToList(r.endpoints.ProjectNameQuery(project, app), "list version for project '"+project+"/"+app+"'")
}
