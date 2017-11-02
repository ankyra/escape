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

package controllers

import (
	"fmt"

	. "github.com/ankyra/escape/model/interfaces"
)

type RegistryController struct{}

func (r RegistryController) Query(context Context, project, application, appVersion string) error {
	registry := context.GetRegistry()
	var result []string
	var err error
	if project == "" {
		result, err = registry.ListProjects()
	} else if application == "" {
		result, err = registry.ListApplications(project)
	} else if appVersion == "" {
		result, err = registry.ListVersions(project, application)
	} else {
		metadata, err := registry.QueryReleaseMetadata(project, application, appVersion)
		if err != nil {
			return err
		}
		fmt.Println(metadata.ToJson())
		return nil
	}
	if err != nil {
		return err
	}
	for _, line := range result {
		if application == "" {
			fmt.Println(line)
		} else {
			fmt.Printf("v%s\n", line)
		}
	}
	if len(result) == 0 {
		fmt.Println("Registry returned 0 results.")
	}
	return nil
}
