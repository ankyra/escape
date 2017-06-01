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
	"net/url"
	"strings"
)

type ServerEndpoints struct {
	apiServer string
}

func NewServerEndpoints(apiServer string) *ServerEndpoints {
	return &ServerEndpoints{
		apiServer: apiServer,
	}
}
func (s *ServerEndpoints) ApiServer() string {
	if s.apiServer == "" {
		return ""
	}
	if strings.HasSuffix(s.apiServer, "/") {
		return s.apiServer
	}
	return s.apiServer + "/"
}
func (s *ServerEndpoints) ReleaseQuery(project, name, version string) string {
	return s.ApiServer() + "a/" + project + "/" + name + "/" + version + "/"
}

func (s *ServerEndpoints) NextReleaseVersion(project, name, prefix string) string {
	v := url.Values{}
	v.Set("prefix", prefix)
	return s.ProjectNameQuery(project, name) + "next-version?" + v.Encode()
}
func (s *ServerEndpoints) ProjectQuery(project string) string {
	return s.ApiServer() + "a/" + project + "/"
}
func (s *ServerEndpoints) ProjectNameQuery(project, name string) string {
	return s.ProjectQuery(project) + name + "/"
}
func (s *ServerEndpoints) ProjectReleaseQuery(project, name, version string) string {
	return s.ProjectNameQuery(project, name) + "v" + version + "/"
}
func (s *ServerEndpoints) RegisterPackage(project string) string {
	return s.ProjectQuery(project) + "register"
}
func (s *ServerEndpoints) UploadRelease(project, name, version string) string {
	return s.ProjectReleaseQuery(project, name, version) + "upload"
}
func (s *ServerEndpoints) DownloadRelease(project, name, version string) string {
	return s.ProjectReleaseQuery(project, name, version) + "download"
}
func (s *ServerEndpoints) AuthMethods(baseUrl string) string {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	return baseUrl + "auth-methods"
}
