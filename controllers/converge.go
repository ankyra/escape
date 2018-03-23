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
		return convergeSmoke(context, depl, releaseId, "converge.test_pending")
	}
	if status.Code == state.DestroyPending {
		return convergeDestroy(context, depl, releaseId, "converge.destroy_pending")
	}
	if status.Code == state.DestroyAndDeletePending {
		err := convergeDestroy(context, depl, releaseId, "converge.destroy_pending")
		if err == nil {
			return context.GetEnvironmentState().DeleteDeployment(depl.Name)
		}
		return err
	}
	if status.IsError() {
		return handleExponentialBackoff(context, depl, status)
	}
	if !refresh && status.Code == state.OK {
		context.Log("converge.skip_ok", map[string]string{
			"deployment": depl.Name,
			"release":    releaseId,
		})
		return nil
	}
	return convergeDeploy(context, depl, releaseId, "converge")
}

func handleExponentialBackoff(context Context, depl *state.DeploymentState, status *state.Status) error {
	now := time.Now()

	// The action has not been retried so set an initial retry time and save
	// the new status.
	if status.TryAgainAt == nil || status.TryAgainAt.IsZero() {
		backOff := time.Duration(BackoffStart) * time.Second
		now = now.Add(backOff)
		status.TryAgainAt = &now
		context.Log("converge.mark_retry", map[string]string{
			"deployment": depl.Name,
			"backoff":    backOff.String(),
		})
		return depl.UpdateStatus(state.DeployStage, status)
	}

	// Retry action
	if status.TryAgainAt.Before(now) {
		return retryAction(context, depl)
	}

	// The action will be retried in a later round. Do nothing.
	context.Log("converge.skip_retry_later", map[string]string{
		"deployment": depl.Name,
		"retriedIn":  status.TryAgainAt.Sub(now).String(),
	})
	return nil
}

func retryAction(context Context, depl *state.DeploymentState) error {
	stage := depl.GetStageOrCreateNew(state.DeployStage)
	status := stage.Status
	releaseId := depl.Release + "-v" + stage.Version
	var err error
	if status.Code == state.Failure {
		err = convergeDeploy(context, depl, releaseId, "converge.deploy_retry")
	} else if status.Code == state.TestFailure {
		err = convergeSmoke(context, depl, releaseId, "converge.test_retry")
	} else if status.Code == state.DestroyFailure {
		err = convergeDestroy(context, depl, releaseId, "converge.destroy_retry")
	} else {
		return fmt.Errorf("Unknown error status '%s'. This is a bug in Esape.", status.Code)
	}
	if err == nil {
		return nil
	}
	now := time.Now()
	status.Tried += 1
	backOff := time.Duration(BackoffStart*math.Exp(float64(status.Tried))) * time.Second
	now = now.Add(backOff)
	status.TryAgainAt = &now
	context.Log("converge.mark_retry", map[string]string{
		"deployment": depl.Name,
		"backoff":    backOff.String(),
	})
	if err2 := depl.UpdateStatus(state.DeployStage, status); err2 != nil {
		return fmt.Errorf("Couldn't update status '%s'. Trying to set failure status, because: %s", err2.Error(), err.Error())
	}
	return err
}

func convergeSmoke(context Context, depl *state.DeploymentState, releaseId, logKey string) error {
	context.Log(logKey, map[string]string{
		"project":     depl.GetEnvironmentState().GetProjectName(),
		"environment": depl.GetEnvironmentState().Name,
		"deployment":  depl.Name,
		"release":     releaseId,
		"action":      "smoke",
	})
	return SmokeController{}.FetchAndSmoke(context, releaseId)
}

func convergeDestroy(context Context, depl *state.DeploymentState, releaseId, logKey string) error {
	context.Log(logKey, map[string]string{
		"project":     depl.GetEnvironmentState().GetProjectName(),
		"environment": depl.GetEnvironmentState().Name,
		"deployment":  depl.Name,
		"release":     releaseId,
		"action":      "destroy",
	})
	return DestroyController{}.FetchAndDestroy(context, releaseId, false, true)
}

func convergeDeploy(context Context, depl *state.DeploymentState, releaseId, logKey string) error {
	context.Log(logKey, map[string]string{
		"project":     depl.GetEnvironmentState().GetProjectName(),
		"environment": depl.GetEnvironmentState().Name,
		"deployment":  depl.Name,
		"release":     releaseId,
		"action":      "deploy",
	})
	return DeployController{}.FetchAndDeploy(context, releaseId, nil, nil)
}
