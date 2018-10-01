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

package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/inventory/types"
	"github.com/ankyra/escape/util"
)

type LocalInventory struct {
	BaseDir string
}

func NewLocalInventory(baseDir string) *LocalInventory {
	return &LocalInventory{
		BaseDir: baseDir,
	}
}

func (r *LocalInventory) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	if version == "latest" || strings.HasSuffix(version, ".@") {
		return nil, fmt.Errorf("Dynamic version release querying not implemented in local inventory. The inventory can be configured in the Global Escape configuration (see `escape config`)")
	}
	return nil, fmt.Errorf("Not implemented")
}

func (r *LocalInventory) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return "", fmt.Errorf("Auto versioning is not implemented in local inventory. The inventory can be configured in the global Escape configuration (see `escape config`)")
}

func (r *LocalInventory) DownloadRelease(project, name, version, targetFile string) error {
	return fmt.Errorf("Release download not implemented in local inventory. The inventory can be configured in the Global Escape configuration (see `escape config`)")
}

func (r *LocalInventory) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return nil
}

func (r *LocalInventory) ListProjects() ([]string, error) {
	path := r.BaseDir
	result := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		if file.IsDir() {
			name := file.Name()
			if !strings.HasPrefix(name, ".") {
				result = append(result, name)
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

func (r *LocalInventory) ListApplications(project string) ([]string, error) {
	path := filepath.Join(r.BaseDir, project)
	if !util.PathExists(path) {
		return nil, fmt.Errorf("The project '%s' could not be found in the local inventory at %s.", project, r.BaseDir)
	}
	result := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		if file.IsDir() {
			name := file.Name()
			indexPath := filepath.Join(r.BaseDir, project, name, "index.json")
			if util.PathExists(indexPath) {
				result = append(result, name)
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

type VersionIndex struct {
	Name          string
	EscapeVersion string
	CoreVersion   string
	Versions      map[string]*core.ReleaseMetadata
}

func (r *LocalInventory) ListVersions(project, app string) ([]string, error) {
	path := filepath.Join(r.BaseDir, project, app, "index.json")
	if !util.PathExists(path) {
		return nil, fmt.Errorf("The application '%s/%s' could not be found in the local inventory at %s.", project, app, r.BaseDir)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	index := VersionIndex{}
	if err := json.Unmarshal(content, &index); err != nil {
		return nil, fmt.Errorf("Could not read local version index at '%s': %s", path, err.Error())
	}

	result := []string{}
	for version := range index.Versions {
		result = append(result, version)
	}
	return result, nil
}

// Not required.
func (r *LocalInventory) Login(url, username, password string) (string, error)    { return "", nil }
func (r *LocalInventory) LoginWithBasicAuth(url, username, password string) error { return nil }
func (r *LocalInventory) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	return nil, nil
}
