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
	"github.com/ankyra/escape-client/model/registry/local"
	"github.com/ankyra/escape-client/model/registry/remote"
	"github.com/ankyra/escape-client/model/registry/types"
	core "github.com/ankyra/escape-core"
)

type Registry interface {
	QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error)
	QueryPreviousReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error)
	QueryNextVersion(project, name, versionPrefix string) (string, error)
	DownloadRelease(project, name, version, targetFile string) error
	UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error
	GetAuthMethods(url string) (map[string]*types.AuthMethod, error)
	LoginWithSecretToken(url, username, password string) (string, error)
}

func NewLocalRegistry() Registry {
	return local.NewLocalRegistry()
}

func NewRemoteRegistry(apiServer, authToken string) Registry {
	return remote.NewRemoteRegistry(apiServer, authToken)
}
