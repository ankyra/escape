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

type Client interface {
	Login(url, username, password string, storeCredentials bool) error
	ReleaseQuery(releaseQuery string) (*core.ReleaseMetadata, error)
	NextVersionQuery(releaseId, prefix string) (string, error)
	Register(*core.ReleaseMetadata) error
	UploadRelease(releaseId, releasePath string) error
	DownloadRelease(releaseId, targetDir string) error
}
