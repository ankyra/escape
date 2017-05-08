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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/runners"
	"github.com/ankyra/escape-client/model/runners/errand"
)

type ErrandsController struct{}

func (ErrandsController) List(context Context) error {
	metadata := context.GetReleaseMetadata()
	if metadata.GetErrands() == nil {
		return nil
	}
	for _, errand := range metadata.GetErrands() {
		description := "No description given."
		if errand.GetDescription() != "" {
			description = errand.GetDescription()
		}
		fmt.Println("- " + errand.GetName() + ": " + description)
	}
	return nil
}

func (ErrandsController) Run(context Context, errandStr string) error {
	//        applog("errand.start", errand=errand, release=escape_plan.get_versionless_build_id())
	metadata := context.GetReleaseMetadata()
	if metadata.GetErrands() == nil {
		return fmt.Errorf("This release doesn't have any errands.")
	}
	errandObj, ok := metadata.GetErrands()[errandStr]
	if !ok {
		return fmt.Errorf("The errand '%s' does not exist", errandStr)
	}
	runner := errand.NewErrandRunner(errandObj)
	runnerContext, err := runners.NewRunnerContext(context)
	if err != nil {
		return err
	}
	return runner.Run(runnerContext)
}
