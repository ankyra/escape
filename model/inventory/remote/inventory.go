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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/inventory/types"
	"github.com/ankyra/escape/model/remote"
)

type inventory struct {
	client    *remote.InventoryClient
	endpoints *remote.ServerEndpoints
}

func NewRemoteInventory(apiServer, escapeToken string, insecureSkipVerify bool) *inventory {
	return &inventory{
		client:    remote.NewRemoteClient(escapeToken, insecureSkipVerify),
		endpoints: remote.NewServerEndpoints(apiServer),
	}
}

func (r *inventory) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	if !strings.HasPrefix(version, "v") && version != "latest" {
		version = "v" + version
	}
	releaseQuery := project + "/" + name + "-" + version
	if project == "_" {
		releaseQuery = name + "-" + version
	}

	url := r.endpoints.ReleaseQuery(project, name, version)
	apiServer := r.endpoints.ApiServer()
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get release metadata for '%s', because the Inventory at '%s' could not be reached: %s", releaseQuery, apiServer, err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("You don't have a valid authentication token for the Inventory at %s. Use `escape login --url %s` to login.", apiServer, apiServer)
	} else if resp.StatusCode == 403 {
		return nil, fmt.Errorf(`You don't have permissions to view the '%s' release in the Inventory at %s. Please ask an administrator for access.`, releaseQuery, apiServer)
	} else if resp.StatusCode == 404 {
		return nil, fmt.Errorf(`Dependency '%s' could not be found. It may not exist in the Inventory you're using (%s) and you need to release it first, or you may not have been given access to it.`, releaseQuery, apiServer)
	} else if resp.StatusCode == 500 {
		return nil, fmt.Errorf(`Couldn't get release metadata for '%s', because the Inventory at %s responded with a server-side error code. Please try again or contact an administrator if the problem persists.`, releaseQuery, apiServer)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(`Couldn't get release metadata for '%s', because the Inventory at '%s' responded with status code %d: %s`, releaseQuery, apiServer, resp.StatusCode, body)
	}
	metadata, err := core.NewReleaseMetadataFromJsonString(body)
	if err != nil {
		return nil, fmt.Errorf(`The Inventory returned release metadata for '%s/%s-%s' that could not be understood: %s`, project, name, version, err.Error())
	}
	return metadata, nil
}

func (r *inventory) QueryPreviousReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
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

func (r *inventory) QueryNextVersion(project, name, versionPrefix string) (string, error) {
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

func (r *inventory) DownloadRelease(project, name, version, targetFile string) error {
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

func (r *inventory) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
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

func (r *inventory) register(project string, metadata *core.ReleaseMetadata) error {
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

func (r *inventory) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	authUrl := r.endpoints.AuthMethods(url)
	resp, err := r.client.GET(authUrl)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get auth methods from server '%s'", url)
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get auth methods from server '%s': %s", authUrl, resp.Status)
		}
		return nil, fmt.Errorf("Couldn't get auth methods from server '%s': %s\n%s", authUrl, resp.Status, string(result))
	}
	result := map[string]*types.AuthMethod{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *inventory) LoginWithSecretToken(url, username, password string) (string, error) {
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

func (r *inventory) urlToList(url, doingMessage string, transformToList func(map[string]interface{}) []string) ([]string, error) {
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf(doingMessage)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(doingMessage)
	}
	result := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return transformToList(result), nil
}

var helpText = "It may not exist in the inventory you're using (%s) and you need to release it first, or you may not have been given access to it."

func (r *inventory) ListProjects() ([]string, error) {
	return r.urlToList(r.endpoints.ListProjects(), "list projects", func(result map[string]interface{}) []string {
		projects := []string{}
		for key, _ := range result {
			projects = append(projects, key)
		}
		return projects
	})
}
func (r *inventory) ListApplications(project string) ([]string, error) {
	return r.urlToList(r.endpoints.ListApplications(project), fmt.Sprintf("Project '%s' could not be found. "+helpText, project, r.endpoints.ApiServer()), func(result map[string]interface{}) []string {
		projects := []string{}
		for key, _ := range result {
			projects = append(projects, key)
		}
		return projects
	})
}
func (r *inventory) ListVersions(project, app string) ([]string, error) {
	return r.urlToList(r.endpoints.ProjectNameQuery(project, app), fmt.Sprintf("Application '%s' could not be found. "+helpText, app, r.endpoints.ApiServer()), func(result map[string]interface{}) []string {
		versions := make([]string, len(result["versions"].([]interface{})))
		for _, v := range result["versions"].([]interface{}) {
			versions = append(versions, v.(string))
		}
		return versions
	})
}
