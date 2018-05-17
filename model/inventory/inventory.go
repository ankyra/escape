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

package inventory

import (
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/inventory/local"
	"github.com/ankyra/escape/model/inventory/remote"
	"github.com/ankyra/escape/model/inventory/types"
)

type Inventory interface {
	QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error)
	QueryNextVersion(project, name, versionPrefix string) (string, error)
	DownloadRelease(project, name, version, targetFile string) error
	UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error
	GetAuthMethods(url string) (map[string]*types.AuthMethod, error)
	Login(url, username, password string) (string, error)
	LoginWithBasicAuth(url, username, password string) error

	ListProjects() ([]string, error)
	ListApplications(project string) ([]string, error)
	ListVersions(project, app string) ([]string, error)
}

func NewLocalInventory() Inventory {
	return local.NewLocalInventory()
}

func NewRemoteInventory(apiServer, authToken string, insecureSkipVerify bool) Inventory {
	return remote.NewRemoteInventory(apiServer, authToken, insecureSkipVerify)
}
