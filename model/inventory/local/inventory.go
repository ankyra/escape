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
	"fmt"
	"io/ioutil"
	"os"
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
	metaPath := filepath.Join(r.BaseDir, project, name, name+"-"+version+".json")
	return core.NewReleaseMetadataFromFile(metaPath)
}

func (r *LocalInventory) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	indexPath := filepath.Join(r.BaseDir, project, name, "index.json")
	index, err := LoadVersionIndexFromFileOrCreateNew(name, indexPath)
	if err != nil {
		return "", fmt.Errorf("Failed to load application index for local inventory: %s", err.Error())
	}
	versions := []string{}
	for version := range index.Versions {
		versions = append(versions, version)
	}
	semver := getMaxFromVersions(versions, versionPrefix)
	semver.OnlyKeepLeadingVersionPart()
	if err := semver.IncrementSmallest(); err != nil {
		return "", err
	}
	return versionPrefix + semver.ToString(), nil
}

func getMaxFromVersions(versions []string, prefix string) *core.SemanticVersion {
	current := core.NewSemanticVersion("-1")
	for _, v := range versions {
		if strings.HasPrefix(v, prefix) {
			release_version := v[len(prefix):]
			newver := core.NewSemanticVersion(release_version)
			if current.LessOrEqual(newver) {
				current = newver
			}
		}
	}
	return current
}

func (r *LocalInventory) DownloadRelease(project, name, version, targetFile string) error {
	path := filepath.Join(r.BaseDir, project, name, name+"-"+version+".tgz")
	if !util.PathExists(path) {
		return fmt.Errorf("The release %s/%s-%s could not be found in the local inventory (expected at %s)", project, name, version, path)
	}
	return util.CopyFile(path, targetFile)
}

func (r *LocalInventory) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	path := filepath.Join(r.BaseDir, project, metadata.Name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("Could not create application directory '%s': %s", path, err.Error())
	}
	dstPath := filepath.Join(path, metadata.GetReleaseId()+".tgz")
	if err := util.CopyFile(releasePath, dstPath); err != nil {
		return fmt.Errorf("Couldn't copy %s to %s: %s", releasePath, dstPath, err.Error())
	}
	metaPath := filepath.Join(path, metadata.GetReleaseId()+".json")
	if err := metadata.WriteJsonFile(metaPath); err != nil {
		return fmt.Errorf("Could not write release metadata file %s: %s", metaPath, err.Error())
	}
	indexPath := filepath.Join(r.BaseDir, project, metadata.Name, "index.json")
	index, err := LoadVersionIndexFromFileOrCreateNew(metadata.Name, indexPath)
	if err != nil {
		return fmt.Errorf("Failed to load application index for local inventory: %s", err.Error())
	}
	if err := index.AddRelease(metadata); err != nil {
		return err
	}
	return index.Save()
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

func (r *LocalInventory) ListVersions(project, app string) ([]string, error) {
	path := filepath.Join(r.BaseDir, project, app, "index.json")
	if !util.PathExists(path) {
		return nil, fmt.Errorf("The application '%s/%s' could not be found in the local inventory at %s.", project, app, r.BaseDir)
	}
	index, err := LoadVersionIndexFromFile(path)
	if err != nil {
		return nil, err
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
