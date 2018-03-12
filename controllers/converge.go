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
	"math"
	"time"

	"github.com/ankyra/escape-core/state"
	. "github.com/ankyra/escape/model/interfaces"
)

const BackoffStart = 1

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
	stage := depl.GetStageOrCreateNew(state.DeployStage)
	if stage.Version == "" {
		return fmt.Errorf("No 'version' set for deployment of '%s' in deployment '%s'",
			depl.Release, depl.Name)
	}
	releaseId := depl.Release + "-v" + stage.Version
	status := stage.Status
	context.SetRootDeploymentName(depl.Name)
	if status.Code == state.TestPending {
		return SmokeController{}.FetchAndSmoke(context, releaseId)
	}
	if status.Code == state.DestroyPending {
		return DestroyController{}.FetchAndDestroy(context, releaseId, false, true)
	}
	if status.IsError() {

		now := time.Now()
		if status.TryAgainAt.IsZero() {
			// The action has not been retried so set an initial retry time and save
			// the new status.
			status.TryAgainAt = now.Add(time.Duration(BackoffStart) * time.Second)
			return depl.UpdateStatus(state.DeployStage, status)
		}
		if stage.Status.TryAgainAt.Before(now) {
			// The action has to be retried.
			return retryAction(context, depl)
		} else {
			// The action will be retried in a later round. Do nothing.
			return nil
		}
	}
	if !refresh && status.Code == state.OK {
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
	return DeployController{}.FetchAndDeploy(context, releaseId, nil, nil)
}

func retryAction(context Context, depl *state.DeploymentState) error {
	stage := depl.GetStageOrCreateNew(state.DeployStage)
	status := stage.Status
	releaseId := depl.Release + "-v" + stage.Version
	var err error
	if status.Code == state.Failure {
		err = DeployController{}.FetchAndDeploy(context, releaseId, nil, nil)
	} else if status.Code == state.TestFailure {
		err = SmokeController{}.FetchAndSmoke(context, releaseId)
	} else if status.Code == state.DestroyFailure {
		err = DestroyController{}.FetchAndDestroy(context, releaseId, false, true)
	} else {
		return fmt.Errorf("Unknown error status '%s'. This is a bug in Esape.", status.Code)
	}
	if err == nil {
		return nil
	}
	now := time.Now()
	status.Tried += 1
	backOff := time.Duration(BackoffStart*math.Exp(float64(status.Tried))) * time.Second
	fmt.Println(backOff)
	status.TryAgainAt = now.Add(backOff)
	if err2 := depl.UpdateStatus(state.DeployStage, status); err2 != nil {
		return fmt.Errorf("Couldn't update status '%s'. Trying to set failure status, because: %s", err2.Error(), err.Error())
	}
	return err
}
