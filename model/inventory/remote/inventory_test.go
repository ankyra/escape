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
	"io/ioutil"
	"os"
	"testing"

	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape/testing"
	. "gopkg.in/check.v1"
)

type suite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&suite{})

const queryReleaseMetadataURL = "/api/v1/registry/query-project/units/name/versions/v1.0.0/"
const queryReleaseMetadataLatestURL = "/api/v1/registry/query-project/units/name/versions/latest/"
const queryNextVersionURL = "/api/v1/registry/query-project/units/name/next-version"
const listProjectsURL = "/api/v1/registry/"
const listApplicationsURL = "/api/v1/registry/test/units/"
const listVersionsURL = "/api/v1/registry/test/units/app/"
const authMethodsURL = "/api/v1/auth/login-methods"
const downloadURL = "/api/v1/registry/prj/units/name/versions/v1.0/download"
const uploadURL = "/api/v1/registry/prj/upload"
const registerURL = "/api/v1/registry/prj/register"

/*

	QUERY RELEASE METADATA

*/

const validMetadata = `{"name": "name", "project": "query-project", "version": "1.0"}`

func (s *suite) Test_QueryReleaseMetadata_happy_path(c *C) {
	server := NewMockServer().WithBody(validMetadata).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	server.ExpectCalled(c, true, queryReleaseMetadataURL)
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) Test_QueryReleaseMetadata_happy_path_for_latest(c *C) {
	server := NewMockServer().WithBody(validMetadata).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "latest")
	server.ExpectCalled(c, true, queryReleaseMetadataLatestURL)
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) Test_QueryReleaseMetadata_happy_path_versions_prefixed_with_v(c *C) {
	server := NewMockServer().WithBody(validMetadata).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	metadata, err := unit.QueryReleaseMetadata("query-project", "name", "v1.0.0")
	server.ExpectCalled(c, true, queryReleaseMetadataURL)
	c.Assert(err, IsNil)
	c.Assert(metadata, Not(IsNil))
	c.Assert(metadata.Name, Equals, "name")
	c.Assert(metadata.Project, Equals, "query-project")
	c.Assert(metadata.Version, Equals, "1.0")
}

func (s *suite) queryReleaseMetadata(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.QueryReleaseMetadata("query-project", "name", "1.0.0")
	return err
}

func (s *suite) Test_QueryReleaseMetadata_Errors(c *C) {
	baseError := fmt.Sprintf(error_QueryReleaseMetadata, "query-project/name-v1.0.0")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_QueryReleaseMetadataForbidden, "query-project/name-v1.0.0", url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_QueryReleaseMetadataNotFound, url+"/", "query-project/name-v1.0.0")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, queryReleaseMetadataURL, s.queryReleaseMetadata)
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.queryReleaseMetadata, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, queryReleaseMetadataURL, url[7:])
		return fmt.Sprintf(error_QueryReleaseMetadata+error_InventoryConnection, "query-project/name-v1.0.0", url+"/", err)
	})
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_empty_body_is_returned(c *C) {
	server := NewMockServer().Start(c)
	defer server.Stop()
	err := s.queryReleaseMetadata(server.URL)
	server.ExpectCalled(c, true, queryReleaseMetadataURL)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Couldn't unmarshal JSON release metadata: unexpected end of JSON input")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_invalid_json_is_returned(c *C) {
	server := NewMockServer().WithBody(`{`).Start(c)
	defer server.Stop()
	err := s.queryReleaseMetadata(server.URL)
	server.ExpectCalled(c, true, queryReleaseMetadataURL)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Couldn't unmarshal JSON release metadata: unexpected end of JSON input")
}

func (s *suite) Test_QueryReleaseMetadata_fails_if_invalid_metadata_is_returned(c *C) {
	server := NewMockServer().WithBody(`{"name": "name", "project": "query-project"}`).Start(c)
	defer server.Stop()
	err := s.queryReleaseMetadata(server.URL)
	server.ExpectCalled(c, true, queryReleaseMetadataURL)
	c.Assert(err.Error(), Equals, "The Inventory returned release metadata for 'query-project/name-v1.0.0' that could not be understood: Missing version field in release metadata")
}

/*

	QUERY NEXT VERSION

*/

func (s *suite) Test_QueryNextVersion_happy_path(c *C) {
	server := NewMockServer().WithBody(`1.0`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	version, err := unit.QueryNextVersion("query-project", "name", "1.@")
	server.ExpectCalled(c, true, queryNextVersionURL)
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
}

func (s *suite) queryNextVersion(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.QueryNextVersion("query-project", "name", "1.0.1")
	return err
}

func (s *suite) Test_QueryNextVersion_Errors(c *C) {
	baseError := fmt.Sprintf(error_QueryNextVersion, "query-project/name-v1.0.1")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_QueryReleaseMetadataForbidden, "query-project/name-v1.0.1", url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, queryNextVersionURL, s.queryNextVersion)
}

func (s *suite) Test_QueryNextVersion_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.queryNextVersion, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, queryNextVersionURL+"?prefix=1.0.1", url[7:])
		return fmt.Sprintf(error_QueryNextVersion+error_InventoryConnection, "query-project/name-v1.0.1", url+"/", err)
	})
}

/*

	LIST PROJECTS

*/

func (s *suite) Test_ListProjects_happy_path(c *C) {
	server := NewMockServer().WithBody(`{"prj": {}, "prj2": {}}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	projects, err := unit.ListProjects()
	server.ExpectCalled(c, true, listProjectsURL)
	c.Assert(err, IsNil)
	c.Assert(projects, HasLen, 2)
	c.Assert(projects[0], Equals, "prj")
	c.Assert(projects[1], Equals, "prj2")
}

func (s *suite) listProjects(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.ListProjects()
	return err
}

func (s *suite) Test_ListProjects_Errors(c *C) {
	baseError := error_ListProjects
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_ListProjectForbidden, url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 404, "Server Error")
		},
	}, listProjectsURL, s.listProjects)
}

func (s *suite) Test_ListProjects_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.listProjects, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, listProjectsURL, url[7:])
		return fmt.Sprintf(error_ListProjects+error_InventoryConnection, url+"/", err)
	})
}

/*

	LIST APPLICATIONS

*/

func (s *suite) Test_ListApplications_happy_path(c *C) {
	server := NewMockServer().WithBody(`{"prj": {}, "prj2": {}}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	apps, err := unit.ListApplications("test")
	server.ExpectCalled(c, true, listApplicationsURL)
	c.Assert(err, IsNil)
	c.Assert(apps, HasLen, 2)
	c.Assert(apps[0], Equals, "prj")
	c.Assert(apps[1], Equals, "prj2")
}

func (s *suite) listApplications(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.ListApplications("test")
	return err
}

func (s *suite) Test_ListApplications_Errors(c *C) {
	baseError := fmt.Sprintf(error_ListApplications, "test")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_ListProjectForbidden, url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_ListApplicationsNotFound, "test", url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, listApplicationsURL, s.listApplications)
}

func (s *suite) Test_ListApplications_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.listApplications, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, listApplicationsURL, url[7:])
		return fmt.Sprintf(error_ListApplications+error_InventoryConnection, "test", url+"/", err)
	})
}

/*

	LIST VERSIONS

*/

func (s *suite) Test_ListVersions_happy_path(c *C) {
	server := NewMockServer().WithBody(`{"versions": ["1.0", "1.1"]}`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	versions, err := unit.ListVersions("test", "app")
	server.ExpectCalled(c, true, listVersionsURL)
	c.Assert(err, IsNil)
	c.Assert(versions, HasLen, 2)
	c.Assert(versions[0], Equals, "1.0")
	c.Assert(versions[1], Equals, "1.1")
}

func (s *suite) listVersions(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.ListVersions("test", "app")
	return err
}

func (s *suite) Test_ListVersions_Errors(c *C) {
	baseError := fmt.Sprintf(error_ListVersions, "app", "test")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_ListProjectForbidden, url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_ListVersionsNotFound, "test", "app", url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, listVersionsURL, s.listVersions)
}

func (s *suite) Test_ListVersions_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.listVersions, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, listVersionsURL, url[7:])
		return fmt.Sprintf(error_ListVersions+error_InventoryConnection, "app", "test", url+"/", err)
	})
}

/*

	GET AUTH METHODS

*/

const validAuthMethods = `
{
	"email": {
		"type": "email",
		"url": "http://test"
	}
}

`

func (s *suite) Test_GetAuthMethods_happy_path(c *C) {
	server := NewMockServer().WithBody(validAuthMethods).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	authMethods, err := unit.GetAuthMethods(server.URL)
	server.ExpectCalled(c, true, authMethodsURL)
	c.Assert(err, IsNil)
	c.Assert(authMethods, HasLen, 1)
	c.Assert(authMethods["email"].Type, Equals, "email")
	c.Assert(authMethods["email"].URL, Equals, "http://test")
}

func (s *suite) getAuthMethods(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.GetAuthMethods(url)
	return err
}

func (s *suite) Test_GetAuthMethods_NotFound_does_not_return_error_but_nil(c *C) {
	server := NewMockServer().WithResponseCode(404).Start(c)
	defer server.Stop()
	unit := NewRemoteInventory(server.URL, "token", false)
	authMethods, err := unit.GetAuthMethods(server.URL)
	server.ExpectCalled(c, true, authMethodsURL)
	c.Assert(err, IsNil)
	c.Assert(authMethods, IsNil)
}

func (s *suite) Test_GetAuthMethods_Errors(c *C) {
	baseError := error_AuthMethods
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url, "Server Error")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url)
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url, 416, "Server Error")
		},
	}, authMethodsURL, s.getAuthMethods)
}

func (s *suite) Test_GetAuthMethods_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.getAuthMethods, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, authMethodsURL, url[7:])
		return fmt.Sprintf(error_AuthMethods+error_InventoryConnection, url, err)
	})
}

/*

	LOGIN

*/

func (s *suite) Test_Login_happy_path(c *C) {
	server := NewMockServer().WithHeader("X-Escape-Token", "my-auth-token").Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	token, err := unit.Login(server.URL, "user", "password")
	server.ExpectCalled(c, true, "/")
	c.Assert(err, IsNil)
	c.Assert(token, Equals, "my-auth-token")
}

func (s *suite) login(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	_, err := unit.Login(url, "user", "password")
	return err
}

func (s *suite) Test_Login_Errors(c *C) {
	baseError := error_Login
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url, "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(baseError + error_LoginCredentials)
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url, 404, "Server Error")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url)
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url, 416, "Server Error")
		},
	}, "/", s.login)
}

func (s *suite) Test_Login_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.login, func(url string) string {
		err := fmt.Sprintf("Post %s: dial tcp %s: getsockopt: connection refused", url, url[7:])
		return fmt.Sprintf(error_Login+error_InventoryConnection, url, err)
	})
}

/*

	DOWNLOAD

*/

func (s *suite) Test_Download_happy_path(c *C) {
	os.RemoveAll("testdata.txt")
	server := NewMockServer().WithBody(`abcdef`).Start(c)
	defer server.Stop()

	unit := NewRemoteInventory(server.URL, "token", false)
	err := unit.DownloadRelease("prj", "name", "1.0", "testdata.txt")
	server.ExpectCalled(c, true, downloadURL)
	c.Assert(err, IsNil)
	content, err := ioutil.ReadFile("testdata.txt")
	c.Assert(err, IsNil)
	c.Assert(string(content), Equals, "abcdef")
	os.RemoveAll("testdata.txt")
}

func (s *suite) download(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	return unit.DownloadRelease("prj", "name", "1.0", "testdata.txt")
}

func (s *suite) Test_Download_Errors(c *C) {
	baseError := fmt.Sprintf(error_Download, "prj/name-v1.0")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_ListProjectForbidden, url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_DownloadNotFound, url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, downloadURL, s.download)
}

func (s *suite) Test_Download_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.download, func(url string) string {
		err := fmt.Sprintf("Get %s%s: dial tcp %s: getsockopt: connection refused", url, downloadURL, url[7:])
		return fmt.Sprintf(error_Download+error_InventoryConnection, "prj/name-v1.0", url+"/", err)
	})
}

/*

	UPLOAD

*/

/*

	REGISTER

*/

func (s *suite) register(url string) error {
	unit := NewRemoteInventory(url, "token", false)
	metadata := core.NewReleaseMetadata("name", "1.0")
	return unit.register("prj", metadata)
}

func (s *suite) Test_Register_Errors(c *C) {
	baseError := fmt.Sprintf(error_Register, "prj", "name-v1.0")
	s.test_RemoteErrorHandling(c, map[int]func(string) string{
		400: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUserSide, url+"/", "Server Error")
		},
		401: func(url string) string {
			return fmt.Sprintf(error_Unauthorized, url+"/", url+"/")
		},
		403: func(url string) string {
			return fmt.Sprintf(baseError+error_ListProjectForbidden, url+"/")
		},
		404: func(url string) string {
			return fmt.Sprintf(baseError+error_DownloadNotFound, url+"/")
		},
		500: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryServerSide, url+"/")
		},
		416: func(url string) string {
			return fmt.Sprintf(baseError+error_InventoryUnknownStatus, url+"/", 416, "Server Error")
		},
	}, registerURL, s.register)
}

func (s *suite) Test_Register_fails_if_server_doesnt_respond(c *C) {
	s.test_ConnectionError(c, s.register, func(url string) string {
		err := fmt.Sprintf("Post %s%s: dial tcp %s: getsockopt: connection refused", url, registerURL, url[7:])
		return fmt.Sprintf(error_Register+error_InventoryConnection, "prj", "name-v1.0", url+"/", err)
	})
}

/*

	HELPER FUNCTIONS

*/

func (s *suite) test_Error(c *C, url string, do func(string) error, expectCode int, errorFunc func(string) string) {
	server := NewMockServer().WithResponseCode(expectCode).WithBody("Server Error").Start(c)
	defer server.Stop()
	err := do(server.URL)
	server.ExpectError(c, err, url, errorFunc(server.URL))
}

func (s *suite) test_ConnectionError(c *C, do func(string) error, errorFunc func(string) string) {
	server := NewMockServer().Start(c)
	server.Stop()
	err := do(server.URL)
	c.Assert(err.Error(), Equals, errorFunc(server.URL))
}

func (s *suite) test_RemoteErrorHandling(c *C, errorMap map[int]func(string) string, url string, do func(string) error) {
	for statusCode, errorFunc := range errorMap {
		s.test_Error(c, url, do, statusCode, errorFunc)
	}
}
