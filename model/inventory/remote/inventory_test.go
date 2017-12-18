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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type suite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&suite{})

type MockServer struct {
	HandlerCalled bool
	Body          string
	ResponseCode  int
	Server        *httptest.Server
	URL           string
	CapturedPath  string
}

func NewMockServer() *MockServer {
	return &MockServer{
		HandlerCalled: false,
		ResponseCode:  200,
	}
}

func (m *MockServer) Start(c *C) *MockServer {
	m.HandlerCalled = false
	m.Server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(m.ResponseCode)
			w.Write([]byte(m.Body))
			m.CapturedPath = r.URL.Path
			m.HandlerCalled = true
		}))

	status := 0
	checkStarted := time.Now()
	for status != m.ResponseCode && time.Now().Before(checkStarted.Add(time.Second*10)) {
		time.Sleep(time.Second / 1000)
		resp, err := http.Get(m.Server.URL)
		if err == nil {
			status = resp.StatusCode
		}
	}
	c.Assert(status, Not(Equals), 0)
	m.HandlerCalled = false
	m.URL = m.Server.URL
	return m
}

func (m *MockServer) WithBody(body string) *MockServer {
	m.Body = body
	return m
}

func (m *MockServer) WithResponseCode(code int) *MockServer {
	m.ResponseCode = code
	return m
}

func (m *MockServer) Stop() {
	m.Server.Close()
}

func (m *MockServer) ExpectCalled(c *C, expectCalled bool, path string) {
	c.Assert(m.HandlerCalled, Equals, expectCalled)
	c.Assert(m.CapturedPath, Equals, path)
}

/*

	QUERY RELEASE METADATA

*/

func (s *suite) Test_QueryReleaseMetadata_happy_path(c *C) {
	server := NewMockServer().WithBody(`{"name": "name", "project": "query-project", "version": "1.0"}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) Test_QueryReleaseMetadata_happy_path_for_latest(c *C) {
	server := NewMockServer().WithBody(`{"name": "name", "project": "query-project", "version": "1.0"}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "latest")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/latest/")
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) Test_QueryReleaseMetadata_happy_path_versions_prefixed_with_v(c *C) {
	server := NewMockServer().WithBody(`{"name": "name", "project": "query-project", "version": "1.0"}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "v1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_empty_body_is_returned(c *C) {
	server := NewMockServer().Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Couldn't unmarshal JSON release metadata: unexpected end of JSON input")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_invalid_json_is_returned(c *C) {
	server := NewMockServer().WithBody(`{`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Couldn't unmarshal JSON release metadata: unexpected end of JSON input")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_invalid_metadata_is_returned(c *C) {
	server := NewMockServer().WithBody(`{"name": "name", "project": "query-project"}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Missing version field in release metadata")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_release_metadata_cant_be_found(c *C) {
	server := NewMockServer().WithResponseCode(404).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, fmt.Sprintf("Dependency 'query-project/name-v1.0.0' could not be found. It may not exist in the Inventory you're using (%s/) and you need to release it first, or you may not have been given access to it.", server.URL))
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_unauthorized(c *C) {
	server := NewMockServer().WithResponseCode(401).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, fmt.Sprintf("You don't have a valid authentication token for the Inventory at %s/. Use `escape login --url %s/` to login.", server.URL, server.URL))
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_forbidden(c *C) {
	server := NewMockServer().WithResponseCode(403).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, fmt.Sprintf("You don't have permissions to view the 'query-project/name-v1.0.0' release in the Inventory at %s/. Please ask an administrator for access.", server.URL))
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_server_error(c *C) {
	server := NewMockServer().WithResponseCode(500).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, fmt.Sprintf("Couldn't get release metadata for 'query-project/name-v1.0.0', because the Inventory at %s/ responded with a server-side error code. Please try again or contact an administrator if the problem persists.", server.URL))
}

func (s *suite) Test_QueryReleaseMetadata_fails_on_other_statuses(c *C) {
	server := NewMockServer().WithResponseCode(416).WithBody("Yo").Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, "/api/v1/registry/query-project/units/name/versions/v1.0.0/")
	c.Assert(err.Error(), Equals, fmt.Sprintf("Couldn't get release metadata for 'query-project/name-v1.0.0', because the Inventory at '%s/' responded with status code 416: Yo", server.URL))
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_server_doesnt_respond(c *C) {
	server := NewMockServer().WithResponseCode(416).WithBody("Yo").Start(c)
	server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	c.Assert(err.Error(), Equals, fmt.Sprintf("Couldn't get release metadata for 'query-project/name-v1.0.0', because the Inventory at '%s/' could not be reached: Get %s/api/v1/registry/query-project/units/name/versions/v1.0.0/: dial tcp %s: getsockopt: connection refused", server.URL, server.URL, server.URL[7:]))
}
