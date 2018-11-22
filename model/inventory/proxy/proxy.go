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

package proxy

import (
	"sort"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/inventory/types"
)

type InventoryProxy struct {
	From              types.Inventory
	To                types.Inventory
	ProxiedNamespaces map[string]bool
}

func NewInventoryProxy(from, to types.Inventory, proxiedNamespaces []string) *InventoryProxy {
	proxy := map[string]bool{}
	for _, p := range proxiedNamespaces {
		proxy[p] = true
	}
	return &InventoryProxy{
		From:              from,
		To:                to,
		ProxiedNamespaces: proxy,
	}
}

func (r *InventoryProxy) GetInventory(project string) types.Inventory {
	_, shouldProxy := r.ProxiedNamespaces[project]
	if shouldProxy {
		return r.To
	}
	return r.From
}

func (r *InventoryProxy) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	return r.GetInventory(project).QueryReleaseMetadata(project, name, version)
}

func (r *InventoryProxy) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return r.GetInventory(project).QueryNextVersion(project, name, versionPrefix)
}

func (r *InventoryProxy) DownloadRelease(project, name, version, targetFile string) error {
	return r.GetInventory(project).DownloadRelease(project, name, version, targetFile)
}

func (r *InventoryProxy) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return r.GetInventory(project).UploadRelease(project, releasePath, metadata)
}

func (r *InventoryProxy) TagRelease(project, name, version, tag string) error {
	return r.GetInventory(project).TagRelease(project, name, version, tag)
}

func (r *InventoryProxy) ListProjects() ([]string, error) {
	projects, err := r.From.ListProjects()
	if err != nil {
		return nil, err
	}
	for p := range r.ProxiedNamespaces {
		projects = append(projects, p)
	}
	sort.Strings(projects)
	return projects, nil
}

func (r *InventoryProxy) ListApplications(project string) ([]string, error) {
	return r.GetInventory(project).ListApplications(project)
}

func (r *InventoryProxy) ListVersions(project, app string) ([]string, error) {
	return r.GetInventory(project).ListVersions(project, app)
}

func (r *InventoryProxy) Login(url, username, password string) (string, error)    { return "", nil }
func (r *InventoryProxy) LoginWithBasicAuth(url, username, password string) error { return nil }
func (r *InventoryProxy) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	return nil, nil
}
