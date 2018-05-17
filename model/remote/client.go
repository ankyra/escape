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
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type InventoryClient struct {
	EscapeToken        string
	BasicAuthUsername  string
	BasicAuthPassword  string
	InsecureSkipVerify bool
}

func NewRemoteClient(escapeToken, basicAuthUsername, basicAuthPassword string, insecureSkipVerify bool) *InventoryClient {
	return &InventoryClient{
		EscapeToken:        escapeToken,
		BasicAuthUsername:  basicAuthUsername,
		BasicAuthPassword:  basicAuthPassword,
		InsecureSkipVerify: insecureSkipVerify,
	}
}

func (c *InventoryClient) GetHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
	}
	return &http.Client{
		Transport: transport,
	}
}

func (c *InventoryClient) NewRequest(method, url string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if c.BasicAuthPassword != "" {
		req.SetBasicAuth(c.BasicAuthUsername, c.BasicAuthPassword)
	}
	return req, nil
}

func (c *InventoryClient) POST_json(url string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) POST_json_with_authentication(url string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) POST_file_with_authentication(url, path string) (*http.Response, error) {
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

	req, err := c.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) PUT_json_with_authentication(url string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("PUT", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) GET_with_authentication(url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) GET(url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.GetHTTPClient().Do(req)
}

func (c *InventoryClient) DELETE_with_authentication(url string) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}

// This function is only used to check credentials (for LoginWithBasicAuth in the inventory).
func (c *InventoryClient) GET_with_basic_authentication(url, username, password string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Escape-Token", c.EscapeToken)
	return c.GetHTTPClient().Do(req)
}
