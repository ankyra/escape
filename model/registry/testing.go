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
	core "github.com/ankyra/escape-core"
)

type mockRegistry struct {
	ReleaseMetadata func(string, string, string) (*core.ReleaseMetadata, error)
	NextVersion     func(string, string, string) (string, error)
	Download        func(string, string, string, string) error
	Upload          func(string, string, *core.ReleaseMetadata) error
}

func NewMockRegistry() *mockRegistry {
	return &mockRegistry{}
}

func (m *mockRegistry) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	return m.ReleaseMetadata(project, name, version)
}
func (m *mockRegistry) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return m.NextVersion(project, name, versionPrefix)
}
func (m *mockRegistry) DownloadRelease(project, name, version, targetFile string) error {
	return m.Download(project, name, version, targetFile)
}
func (m *mockRegistry) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return m.Upload(project, releasePath, metadata)
}
