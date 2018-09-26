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

package controllers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape/model/runners"
	"github.com/ankyra/escape/model/runners/errand"
)

type ErrandsController struct{}

func (ErrandsController) List(context *model.Context) *ControllerResult {
	result := NewControllerResult()

	metadata := context.GetReleaseMetadata()
	errands := metadata.GetErrands()
	result.MarshalableOutput = errands
	if metadata.GetErrands() == nil {
		result.HumanOutput.AddLine("No errands on the deployment")
		return result
	}
	result.HumanOutput.AddLine("Errands:")
	for _, errand := range errands {
		description := "No description given."
		if errand.Description != "" {
			description = errand.Description
		}
		result.HumanOutput.AddLine("- %s\n\n", errand.Name)
		result.HumanOutput.AddLine("  %s\n\n", description)
		for _, input := range errand.Inputs {
			result.HumanOutput.AddLine("  * %s: %s\n", input.Id, input.Description)
		}
	}
	return result
}

func (ErrandsController) Run(context *model.Context, errandStr string, extraVars map[string]interface{}) error {
	//        applog("errand.start", errand=errand, release=escape_plan.get_versionless_build_id())
	metadata := context.GetReleaseMetadata()
	if metadata.GetErrands() == nil {
		return fmt.Errorf("This release doesn't have any errands.")
	}
	errandObj, ok := metadata.GetErrands()[errandStr]
	if !ok {
		return fmt.Errorf("The errand '%s' could not be found in deployment '%s'. You can use 'escape errands list' to see the available errands.", errandStr, context.GetRootDeploymentName())
	}
	runner := errand.NewErrandRunner(errandObj, extraVars)
	runnerContext, err := runners.NewRunnerContext(context)
	if err != nil {
		return err
	}
	return runner.Run(runnerContext)
}

func (e ErrandsController) RunRemoteErrand(context *model.Context, errandStr string, extraVars map[string]interface{}) error {
	name, err := ioutil.TempDir("", "escape-errand")
	if err != nil {
		return fmt.Errorf("Could not create temporary directory for errand: %s", err.Error())
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("Could not change to temporary directory %s: %s", name, err.Error())
	}
	releaseId := context.GetReleaseMetadata().GetQualifiedReleaseId()
	if err := (FetchController{}.Fetch(context, []string{releaseId})); err != nil {
		return err
	}
	if err := os.Chdir(paths.NewPath().UnpackedDepDirectoryByReleaseMetadata(context.GetReleaseMetadata())); err != nil {
		return err
	}
	if err := e.Run(context, errandStr, extraVars); err != nil {
		return err
	}
	if err := os.Chdir(currentDir); err != nil {
		return err
	}
	return os.RemoveAll(name)
}
