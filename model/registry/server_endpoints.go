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
	"net/url"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
)

type ServerEndpoints interface {
	ApiServer() string
	ReleaseQuery(project, name, version string) string
	RegisterPackage(project string) string
	NextReleaseVersion(project, name, prefix string) string
	UploadRelease(project, name, version string) string
	DownloadRelease(project, name, version string) string
}

type serverEndpoints struct {
	escapeConfig EscapeConfig
}

func NewServerEndpoints(cfg EscapeConfig) *serverEndpoints {
	return &serverEndpoints{
		escapeConfig: cfg,
	}
}
func (s *serverEndpoints) ApiServer() string {
	apiServer := s.escapeConfig.GetCurrentTarget().GetApiServer()
	if apiServer == "" {
		return ""
	}
	if strings.HasSuffix(apiServer, "/") {
		return apiServer
	}
	return apiServer + "/"
}
func (s *serverEndpoints) ReleaseQuery(project, name, version string) string {
	return s.ApiServer() + "a/" + project + "/" + name + "/" + version + "/"
}

func (s *serverEndpoints) NextReleaseVersion(project, name, prefix string) string {
	v := url.Values{}
	v.Set("prefix", prefix)
	return s.ProjectNameQuery(project, name) + "next-version?" + v.Encode()
}

func (s *serverEndpoints) ProjectQuery(project string) string {
	return s.ApiServer() + "a/" + project + "/"
}
func (s *serverEndpoints) ProjectNameQuery(project, name string) string {
	return s.ProjectQuery(project) + name + "/"
}
func (s *serverEndpoints) ProjectReleaseQuery(project, name, version string) string {
	return s.ProjectNameQuery(project, name) + "v" + version + "/"
}
func (s *serverEndpoints) RegisterPackage(project string) string {
	return s.ProjectQuery(project) + "register"
}
func (s *serverEndpoints) UploadRelease(project, name, version string) string {
	return s.ProjectReleaseQuery(project, name, version) + "upload"
}
func (s *serverEndpoints) DownloadRelease(project, name, version string) string {
	return s.ProjectReleaseQuery(project, name, version) + "download"
}
