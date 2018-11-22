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

package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape/model/inventory/types"
	"github.com/ankyra/escape/model/remote"
)

type inventory struct {
	apiServer string
	client    *remote.InventoryClient
	endpoints *remote.ServerEndpoints
}

func NewRemoteInventory(apiServer, escapeToken, basicAuthUsername, basicAuthPassword string, insecureSkipVerify bool) *inventory {
	inv := &inventory{
		client:    remote.NewRemoteClient(escapeToken, basicAuthUsername, basicAuthPassword, insecureSkipVerify),
		endpoints: remote.NewServerEndpoints(apiServer),
	}
	inv.apiServer = inv.endpoints.ApiServer()
	return inv
}

const error_InventoryConnection = ", because the Inventory at '%s' could not be reached: %s"
const error_InventoryServerSide = ", because the Inventory at '%s' responded with a server-side error code. Please try again or contact an administrator if the problem persists."
const error_InventoryUserSide = ", because the Inventory at '%s' says there's a problem with the request: %s"
const error_InventoryUnknownStatus = ", because the Inventory at '%s' responded with status code %d: %s"

const error_QueryReleaseMetadata = "Couldn't get release metadata for '%s'"
const error_QueryReleaseMetadataNotFound = ", because the release metadata could not be found in the Inventory at '%s'. You probably need to release the '%s' package first."
const error_QueryReleaseMetadataForbidden = ", because you don't have permission to view the '%s' release in the Inventory at '%s'. Please ask an administrator for access."
const error_Unauthorized = "You don't have a valid authentication token for the Inventory at %s. Use `escape login --url %s` to login."
const error_QueryNextVersion = "Couldn't resolve next version for '%s'"
const error_ListProjects = "Couldn't list projects"
const error_ListApplications = "Couldn't list applications for project '%s'"
const error_ListApplicationsNotFound = ", because the project '%s' could not be found in the Inventory at '%s'."
const error_ListVersions = "Couldn't list versions for application '%s' in project '%s'"
const error_ListVersionsNotFound = ", because the project '%s' or application '%s' could not be found in the Inventory at '%s'."
const error_ListProjectForbidden = ", because you don't have permissions to view this project in the Inventory at '%s'. Please ask an administrator for access."
const error_AuthMethods = "Couldn't get authentication methods from server"
const error_Login = "Couldn't login to the Inventory"
const error_LoginCredentials = ", because the username or password was incorrect."
const error_Download = "Couldn't download release '%s'"
const error_DownloadNotFound = ", because the package could not be found in the Inventory at '%s'"
const error_Upload = "Couldn't upload release '%s/%s'"
const error_Register = "Couldn't register release '%s/%s'"

func (r *inventory) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	query, err := parsers.ParseVersionQuery(version)
	if err != nil {
		return nil, err
	}
	version = query.ToString()
	releaseQuery := project + "/" + name + query.ToVersionSuffix()
	if project == "_" {
		releaseQuery = name + query.ToVersionSuffix()
	}

	url := r.endpoints.ReleaseQuery(project, name, version)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_InventoryConnection, releaseQuery, r.apiServer, err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	if resp.StatusCode == 400 {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_InventoryUserSide, releaseQuery, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return nil, fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_QueryReleaseMetadataForbidden, releaseQuery, releaseQuery, r.apiServer)
	} else if resp.StatusCode == 404 {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_QueryReleaseMetadataNotFound, releaseQuery, r.apiServer, releaseQuery)
	} else if resp.StatusCode == 500 {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_InventoryServerSide, releaseQuery, r.apiServer)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(error_QueryReleaseMetadata+error_InventoryUnknownStatus, releaseQuery, r.apiServer, resp.StatusCode, body)
	}
	metadata, err := core.NewReleaseMetadataFromJsonString(body)
	if err != nil {
		return nil, fmt.Errorf(`The Inventory returned release metadata for '%s/%s-%s' that could not be understood: %s`, project, name, version, err.Error())
	}
	return metadata, nil
}

func (r *inventory) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	url := r.endpoints.NextReleaseVersion(project, name, versionPrefix)
	releaseQuery := project + "/" + name + "-v" + versionPrefix
	if project == "_" {
		releaseQuery = name + "-v" + versionPrefix
	}
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return "", fmt.Errorf(error_QueryNextVersion+error_InventoryConnection, releaseQuery, r.apiServer, err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	if resp.StatusCode == 400 {
		return "", fmt.Errorf(error_QueryNextVersion+error_InventoryUserSide, releaseQuery, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return "", fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return "", fmt.Errorf(error_QueryNextVersion+error_QueryReleaseMetadataForbidden, releaseQuery, releaseQuery, r.apiServer)
	} else if resp.StatusCode == 500 {
		return "", fmt.Errorf(error_QueryNextVersion+error_InventoryServerSide, releaseQuery, r.apiServer)
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(error_QueryNextVersion+error_InventoryUnknownStatus, releaseQuery, r.apiServer, resp.StatusCode, body)
	}
	return body, nil
}

func (r *inventory) ListProjects() ([]string, error) {
	return r.urlToList(r.endpoints.ListProjects(), error_ListProjects, "",
		func(result map[string]interface{}) []string {
			projects := []string{}
			for key, _ := range result {
				projects = append(projects, key)
			}
			sort.Strings(projects)
			return projects
		})
}

func (r *inventory) ListApplications(project string) ([]string, error) {
	return r.urlToList(r.endpoints.ListApplications(project),
		fmt.Sprintf(error_ListApplications, project),
		fmt.Sprintf(error_ListApplicationsNotFound, project, r.apiServer),
		func(result map[string]interface{}) []string {
			apps := []string{}
			for key, _ := range result {
				apps = append(apps, key)
			}
			sort.Strings(apps)
			return apps
		})
}

func (r *inventory) ListVersions(project, app string) ([]string, error) {
	return r.urlToList(r.endpoints.ProjectNameQuery(project, app),
		fmt.Sprintf(error_ListVersions, app, project),
		fmt.Sprintf(error_ListVersionsNotFound, project, app, r.apiServer),
		func(result map[string]interface{}) []string {
			versions := make([]string, len(result["versions"].([]interface{})))
			for i, v := range result["versions"].([]interface{}) {
				versions[i] = v.(string)
			}
			return versions
		})
}

func (r *inventory) urlToList(url, baseErrorMessage, notFoundMessage string, transformToList func(map[string]interface{}) []string) ([]string, error) {
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return nil, fmt.Errorf(baseErrorMessage+error_InventoryConnection, r.apiServer, err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	if resp.StatusCode == 400 {
		return nil, fmt.Errorf(baseErrorMessage+error_InventoryUserSide, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return nil, fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return nil, fmt.Errorf(baseErrorMessage+error_ListProjectForbidden, r.apiServer)
	} else if resp.StatusCode == 404 && notFoundMessage != "" {
		return nil, fmt.Errorf(baseErrorMessage + notFoundMessage)
	} else if resp.StatusCode == 500 {
		return nil, fmt.Errorf(baseErrorMessage+error_InventoryServerSide, r.apiServer)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(baseErrorMessage+error_InventoryUnknownStatus, r.apiServer, resp.StatusCode, body)
	}
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return transformToList(result), nil
}

func (r *inventory) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	authUrl := r.endpoints.AuthMethods(url)
	resp, err := r.client.GET(authUrl)
	if err != nil {
		return nil, fmt.Errorf(error_AuthMethods+error_InventoryConnection, url, err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if resp.StatusCode == 400 {
		return nil, fmt.Errorf(error_AuthMethods+error_InventoryUserSide, url, body)
	} else if resp.StatusCode == 401 {
		realm := resp.Header.Get("WWW-Authenticate")
		if strings.HasPrefix(realm, `Basic realm`) {
			result := map[string]*types.AuthMethod{}
			result["Basic Authentication"] = &types.AuthMethod{
				URL:  url,
				Type: "basic-auth",
			}
			return result, nil
		}
		return nil, fmt.Errorf(error_AuthMethods+error_InventoryUnknownStatus, url, resp.StatusCode, body)
	} else if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode == 500 {
		return nil, fmt.Errorf(error_AuthMethods+error_InventoryServerSide, url)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(error_AuthMethods+error_InventoryUnknownStatus, url, resp.StatusCode, body)
	}
	result := map[string]*types.AuthMethod{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *inventory) LoginWithBasicAuth(url, username, password string) error {
	resp, err := r.client.GET_with_basic_authentication(url, username, password)
	if err != nil {
		return fmt.Errorf(error_Login+error_InventoryConnection, url, err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if resp.StatusCode == 400 {
		return fmt.Errorf(error_Login+error_InventoryUserSide, url, body)
	} else if resp.StatusCode == 401 {
		return fmt.Errorf(error_Login + error_LoginCredentials)
	} else if resp.StatusCode == 500 {
		return fmt.Errorf(error_Login+error_InventoryServerSide, url)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf(error_Login+error_InventoryUnknownStatus, url, resp.StatusCode, body)
	}
	return nil
}

func (r *inventory) Login(url, username, password string) (string, error) {
	payload := map[string]string{
		"username":     username,
		"secret_token": password,
	}
	resp, err := r.client.POST_json_with_authentication(url, payload)
	if err != nil {
		return "", fmt.Errorf(error_Login+error_InventoryConnection, url, err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if resp.StatusCode == 400 {
		return "", fmt.Errorf(error_Login+error_InventoryUserSide, url, body)
	} else if resp.StatusCode == 401 {
		return "", fmt.Errorf(error_Login + error_LoginCredentials)
	} else if resp.StatusCode == 500 {
		return "", fmt.Errorf(error_Login+error_InventoryServerSide, url)
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(error_Login+error_InventoryUnknownStatus, url, resp.StatusCode, body)
	}
	return resp.Header.Get("X-Escape-Token"), nil
}

func (r *inventory) DownloadRelease(project, name, version, targetFile string) error {
	releaseQuery := project + "/" + name + "-v" + version
	if project == "_" {
		releaseQuery = name + "-v" + version
	}
	url := r.endpoints.DownloadRelease(project, name, version)
	resp, err := r.client.GET_with_authentication(url)
	if err != nil {
		return fmt.Errorf(error_Download+error_InventoryConnection, releaseQuery, r.apiServer, err.Error())
	}

	if resp.StatusCode == 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(error_Download+error_InventoryUserSide, releaseQuery, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return fmt.Errorf(error_Download+error_ListProjectForbidden, releaseQuery, r.apiServer)
	} else if resp.StatusCode == 404 {
		return fmt.Errorf(error_Download+error_DownloadNotFound, releaseQuery, r.apiServer)
	} else if resp.StatusCode == 500 {
		return fmt.Errorf(error_Download+error_InventoryServerSide, releaseQuery, r.apiServer)
	} else if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(error_Download+error_InventoryUnknownStatus, releaseQuery, r.apiServer, resp.StatusCode, body)
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
	baseError := fmt.Sprintf(error_Upload, project, metadata.GetReleaseId())
	if err != nil {
		return err
	}
	if resp.StatusCode == 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUserSide, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return fmt.Errorf(baseError+error_ListProjectForbidden, r.apiServer)
	} else if resp.StatusCode == 404 {
		return fmt.Errorf(baseError+error_ListApplicationsNotFound, project, r.apiServer)
	} else if resp.StatusCode == 500 {
		return fmt.Errorf(baseError+error_InventoryServerSide, r.apiServer)
	} else if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUnknownStatus, r.apiServer, resp.StatusCode, body)
	}
	return nil
}

func (r *inventory) register(project string, metadata *core.ReleaseMetadata) error {
	url := r.endpoints.RegisterPackage(project)
	resp, err := r.client.POST_json_with_authentication(url, metadata)
	baseError := fmt.Sprintf(error_Register, project, metadata.GetReleaseId())
	if err != nil {
		return fmt.Errorf(baseError+error_InventoryConnection, r.apiServer, err.Error())
	}
	if resp.StatusCode == 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUserSide, r.apiServer, body)
	} else if resp.StatusCode == 401 {
		return fmt.Errorf(error_Unauthorized, r.apiServer, r.apiServer)
	} else if resp.StatusCode == 403 {
		return fmt.Errorf(baseError+error_ListProjectForbidden, r.apiServer)
	} else if resp.StatusCode == 404 {
		return fmt.Errorf(baseError+error_ListApplicationsNotFound, project, r.apiServer)
	} else if resp.StatusCode == 500 {
		return fmt.Errorf(baseError+error_InventoryServerSide, r.apiServer)
	} else if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUnknownStatus, r.apiServer, resp.StatusCode, body)
	}
	return nil
}

func (r *inventory) TagRelease(project, name, version, tag string) error {
	query, err := parsers.ParseVersionQuery(version)
	if err != nil {
		return err
	}
	url := r.endpoints.TagRelease(project, name)
	data := map[string]interface{}{
		"release_id": project + "/" + name + query.ToVersionSuffix(),
		"tag":        tag,
	}
	resp, err := r.client.POST_json_with_authentication(url, data)
	if err != nil {
		return err
	}
	baseError := fmt.Sprintf("Couldn't tag '%s' with '%s'", data["release_id"], data["tag"])
	if resp.StatusCode == 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUserSide, r.apiServer, body)
	} else if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		return fmt.Errorf(baseError+error_InventoryUnknownStatus, r.apiServer, resp.StatusCode, body)
	}
	return nil
}
