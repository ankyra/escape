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

package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
	core "github.com/ankyra/escape-core"
)

type client struct {
	Config    EscapeConfig
	endpoints ServerEndpoints
}

func NewClient(cfg EscapeConfig) Client {
	return &client{
		Config:    cfg,
		endpoints: NewServerEndpoints(cfg),
	}
}

func (c *client) getHTTPClient() *http.Client {
	return &http.Client{}
}

func (c *client) postJson(url string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return c.getHTTPClient().Do(req)
}
func (c *client) authGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "JWT "+c.Config.GetCurrentTarget().GetAuthToken())
	return c.getHTTPClient().Do(req)
}
func (c *client) authPostJson(url string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "JWT "+c.Config.GetCurrentTarget().GetAuthToken())
	return c.getHTTPClient().Do(req)
}

func (c *client) authPostFile(url, path string) (*http.Response, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", path)
	if err != nil {
		return nil, err
	}
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "JWT "+c.Config.GetCurrentTarget().GetAuthToken())
	return c.getHTTPClient().Do(req)
}

func (c *client) Login(url, username, password string, storeCredentials bool) error {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	authUrl := url + "auth"
	data := map[string]string{
		"username": username,
		"password": password,
	}
	resp, err := c.postJson(authUrl, data)
	if err != nil {
		return errors.New("Failed to login to the Escape server at " + url + ": " + err.Error())
	}
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized.")
	}
	if resp.StatusCode != 200 {
		return errors.New("Oh you done it now. Status: " + resp.Status)
	}
	defer resp.Body.Close()

	result := map[string]string{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	accessToken, ok := result["access_token"]
	if !ok {
		return errors.New("Expecting access token in login payload.")
	}
	configTarget := c.Config.GetCurrentTarget()
	configTarget.SetApiServer(url)
	configTarget.SetAuthToken(accessToken)
	if storeCredentials {
		configTarget.SetUsername(username)
		configTarget.SetPassword(password)
	}
	configTarget.Save()
	return nil
}

func (c *client) ReleaseQuery(project, name, version string) (*core.ReleaseMetadata, error) {
	url := c.endpoints.ReleaseQuery(project, name, version)
	resp, err := c.authGet(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, errors.New("Unauthorized")
	} else if resp.StatusCode != 200 {
		releaseQuery := project + "/" + name + "-" + version
		if project == "_" {
			releaseQuery = name + "-" + version
		}
		return nil, errors.New("Couldn't query release " + releaseQuery + ": " + resp.Status)
	}
	result := core.NewEmptyReleaseMetadata()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *client) DownloadRelease(project, name, version, targetFile string) error {
	url := c.endpoints.DownloadRelease(project, name, version)
	resp, err := c.authGet(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized")
	} else if resp.StatusCode != 200 {
		releaseId := project + "/" + name + "-v" + version
		if project == "_" {
			releaseId = name + "-v" + version
		}
		return errors.New("Couldn't download release " + releaseId + ": " + resp.Status)
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

func (c *client) NextVersionQuery(project, name, prefix string) (string, error) {
	url := c.endpoints.NextReleaseVersion(project, name, prefix)
	resp, err := c.authGet(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 401 {
		return "", errors.New("Unauthorized")
	} else if resp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("There was a problem with the query: %s", body)
	} else if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("Could not query release version.")
		}
		return "", fmt.Errorf("Could not query release version: %s", body)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (c *client) Register(project string, metadata *core.ReleaseMetadata) error {
	url := c.endpoints.RegisterPackage(project)
	resp, err := c.authPostJson(url, metadata)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized")
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Couldn't register package: %s", resp.Status)
		}
		return fmt.Errorf("Couldn't register package (%s): %s", resp.Status, result)
	}
	return nil
}

func (c *client) UploadRelease(project, name, version, releasePath string) error {
	url := c.endpoints.UploadRelease(project, name, version)
	resp, err := c.authPostFile(url, releasePath)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized")
	} else if resp.StatusCode != 200 {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Couldn't upload package (%s)", resp.Status)
		}
		return fmt.Errorf("Couldn't upload package (%s): %s", resp.Status, result)
	}
	return nil
}

func (c *client) loginIfAuthTokenIsNotSet() error {
	configTarget := c.Config.GetCurrentTarget()
	if configTarget.GetApiServer() != "" && configTarget.GetAuthToken() != "" {
		return nil
	}
	return c.tryLoginWithStoredCredentials()
}
func (c *client) tryLoginWithStoredCredentials() error {
	configTarget := c.Config.GetCurrentTarget()
	apiServer := c.endpoints.ApiServer()
	username := configTarget.GetUsername()
	password := configTarget.GetPassword()
	if apiServer != "" && username != "" && password != "" {
		return c.Login(apiServer, username, password, false)
	}
	return errors.New("Please login using `escape login`")
}
