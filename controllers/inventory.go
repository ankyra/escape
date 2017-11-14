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

type InventoryController struct{}

func (r InventoryController) Query(context Context, project, application, appVersion string) *ControllerResult {
	result := NewControllerResult()

	inventory := context.GetInventory()
	var resultData []string
	var err error
	if project == "" {
		resultData, err = inventory.ListProjects()
	} else if application == "" {
		resultData, err = inventory.ListApplications(project)
	} else if appVersion == "" {
		resultData, err = inventory.ListVersions(project, application)
		for i, _ := range resultData {
			resultData[i] = "v" + resultData[i]
		}
	} else {
		metadata, err := inventory.QueryReleaseMetadata(project, application, appVersion)
		if err != nil {
			result.Error = err
			return result
		}
		fmt.Println(metadata.ToJson())
		return nil
	}
	if err != nil {
		result.Error = err
		return result
	}

	result.MarshalableOutput = resultData

	if len(resultData) == 0 {
		result.HumanOutput.AddLine("Inventory returned 0 results.")
		return result
	}

	result.HumanOutput.AddStringList(resultData)
	return result
}
