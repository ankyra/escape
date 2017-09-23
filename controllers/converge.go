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
	"github.com/ankyra/escape-core/state"
)

type ConvergeController struct{}

func (ConvergeController) Converge(context Context, refresh bool) error {
	context.PushLogSection("Converge")
	dag, err := context.GetEnvironmentState().GetDeploymentStateDAG("deploy")
	dag.Walk(func(d *state.DeploymentState) {
		if err != nil {
			return
		}
		err = ConvergeDeployment(context, d, refresh)
	})
	context.PopLogSection()
	return err
}

func ConvergeDeployment(context Context, depl *state.DeploymentState, refresh bool) error {
	if depl.Release == "" {
		return fmt.Errorf("No release set for deployment '%s'", depl.Name)
	}
	stage := depl.GetStageOrCreateNew("deploy")
	if stage.Version == "" {
		return fmt.Errorf("No 'version' set for deployment of '%s' in deployment '%s'",
			depl.Release, depl.Name)
	}
	if !refresh && stage.Status.Code == state.OK {
		context.Log("converge.skip_ok", map[string]string{
			"deployment": depl.Name,
			"release":    depl.Release + "-v" + stage.Version,
		})
		return nil
	}
	context.Log("converge", map[string]string{
		"deployment": depl.Name,
		"release":    depl.Release + "-v" + stage.Version,
	})
	context.SetRootDeploymentName(depl.Name)
	return DeployController{}.FetchAndDeploy(context, depl.Release+"-v"+stage.Version, nil, nil)
}
