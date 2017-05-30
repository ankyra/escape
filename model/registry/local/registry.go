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
	core "github.com/ankyra/escape-core"
)

type LocalRegistry struct {
}

func NewLocalRegistry() *LocalRegistry {
	return &LocalRegistry{}
}

func (r *LocalRegistry) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *LocalRegistry) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return "", fmt.Errorf("Auto versioning is not implemented for local registries. The registry can be configured in the global Escape configuration (see `escape config`)")
}

func (r *LocalRegistry) DownloadRelease(project, name, version, targetFile string) error {
	return fmt.Errorf("Not implemented")
}

func (r *LocalRegistry) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return nil
}
