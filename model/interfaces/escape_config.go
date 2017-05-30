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

package interfaces

import (
	core "github.com/ankyra/escape-core"
)

type Registry interface {
	QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error)
	QueryNextVersion(project, name, versionPrefix string) (string, error)
	DownloadRelease(project, name, version, targetFile string) error
	UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error
}

type EscapeTargetConfig interface {
	ToJson() string
	Save() error

	GetApiServer() string
	GetUsername() string
	GetPassword() string
	GetAuthToken() string
	GetStorageBackend() string
	GetGcsBucketUrl() string
	GetProject() string

	SetApiServer(string)
	SetUsername(string)
	SetPassword(string)
	SetAuthToken(string)
	SetStorageBackend(string)
	SetGcsBucketUrl(string)
}

type EscapeConfig interface {
	LoadConfig(cfgFile string, cfgProfile string) error
	Save() error
	GetCurrentTarget() EscapeTargetConfig
	GetRegistry() Registry
}
