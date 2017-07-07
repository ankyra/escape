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

package local

import (
	"fmt"
	"github.com/ankyra/escape-client/model/registry/types"
	core "github.com/ankyra/escape-core"
	"strings"
)

type LocalRegistry struct{}

func NewLocalRegistry() *LocalRegistry {
	return &LocalRegistry{}
}

func (r *LocalRegistry) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	if version == "latest" || strings.HasSuffix(version, ".@") {
		return nil, fmt.Errorf("Dynamic version release querying not implemented in local registry. The registry can be configured in the Global Escape configuration (see `escape config`)")
	}
	return nil, fmt.Errorf("Not implemented")
}
func (r *LocalRegistry) QueryPreviousReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *LocalRegistry) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return "", fmt.Errorf("Auto versioning is not implemented for local registries. The registry can be configured in the global Escape configuration (see `escape config`)")
}

func (r *LocalRegistry) DownloadRelease(project, name, version, targetFile string) error {
	return fmt.Errorf("Release download not implemented in local registry. The registry can be configured in the Global Escape configuration (see `escape config`)")
}

func (r *LocalRegistry) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return nil
}

func (r *LocalRegistry) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	return nil, nil
}

func (r *LocalRegistry) LoginWithSecretToken(url, username, password string) (string, error) {
	return "", nil
}
