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
	ReleaseQuery(releaseQuery string) string
	NextReleaseVersion(releaseId, prefix string) string
	RegisterPackage() string
	UploadRelease(releaseId string) string
	DownloadRelease(releaseId string) string
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
func (s *serverEndpoints) ReleaseQuery(releaseQuery string) string {
	return s.ApiServer() + "r/" + releaseQuery + "/"
}

func (s *serverEndpoints) NextReleaseVersion(releaseId, prefix string) string {
	v := url.Values{}
	v.Set("prefix", prefix)
	return s.ReleaseQuery(releaseId) + "next-version?" + v.Encode()
}

func (s *serverEndpoints) RegisterPackage() string {
	return s.ApiServer() + "r/"
}

func (s *serverEndpoints) UploadRelease(releaseId string) string {
	return s.ReleaseQuery(releaseId) + "upload"
}
func (s *serverEndpoints) DownloadRelease(releaseId string) string {
	return s.ReleaseQuery(releaseId) + "download"
}
