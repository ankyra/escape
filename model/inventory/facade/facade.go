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

package facade

import (
	"fmt"
	"sort"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/inventory/types"
)

type inventories struct {
	Inventories []types.Inventory
}

func NewInventoriesFacade() types.Inventory {
	return &inventories{
		Inventories: []types.Inventory{},
	}
}

// Gets the release metadata for project/app-version from the first inventory
// that has it.  Ignores other errors so could potentially do something
// unexpected.
func (r *inventories) QueryReleaseMetadata(project, name, version string) (*core.ReleaseMetadata, error) {
	for _, inv := range r.Inventories {
		result, err := inv.QueryReleaseMetadata(project, name, version)
		if err != nil {
			continue
		}
		return result, nil
	}
	return nil, fmt.Errorf("Couldn't get release metadata for '%s/%s-%s'", project, name, version)
}

func (r *inventories) DownloadRelease(project, name, version, targetFile string) error {
	return nil
}

// Combines the projects found in each inventory into a list.
func (r *inventories) ListProjects() ([]string, error) {
	projectSet := map[string]bool{}
	for _, inv := range r.Inventories {
		result, err := inv.ListProjects()
		if err != nil {
			return nil, err
		}
		for _, prj := range result {
			projectSet[prj] = true
		}
	}
	result := []string{}
	for key, _ := range projectSet {
		result = append(result, key)
	}
	sort.Strings(result)
	return result, nil
}

// Combines the applications found in each inventory into a list.
func (r *inventories) ListApplications(project string) ([]string, error) {
	applicationSet := map[string]bool{}
	for _, inv := range r.Inventories {
		result, err := inv.ListApplications(project)
		if err != nil {
			return nil, err
		}
		for _, app := range result {
			applicationSet[app] = true
		}
	}
	result := []string{}
	for key, _ := range applicationSet {
		result = append(result, key)
	}
	sort.Strings(result)
	return result, nil
}

// Lists the versions for project/app from the first inventory that has it.
// Ignore errors so could cause interesting behaviour...
func (r *inventories) ListVersions(project, app string) ([]string, error) {
	for _, inv := range r.Inventories {
		result, err := inv.ListVersions(project, app)
		if err != nil {
			continue
		}
		return result, nil
	}
	return nil, fmt.Errorf("Couldn't list versions for application '%s' in project '%s', because the project '%s' or application '%s' could not be found in any of the configured Inventories.", app, project, project, app)
}

// These methods are only implemented for the actual Inventory types.
func (r *inventories) Login(url, username, password string) (string, error)    { return "", nil }
func (r *inventories) LoginWithBasicAuth(url, username, password string) error { return nil }
func (r *inventories) GetAuthMethods(url string) (map[string]*types.AuthMethod, error) {
	return nil, nil
}
func (r *inventories) QueryNextVersion(project, name, versionPrefix string) (string, error) {
	return "", nil
}
func (r *inventories) UploadRelease(project, releasePath string, metadata *core.ReleaseMetadata) error {
	return nil
}
