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
	"os"

	"github.com/ankyra/escape-core/state"
	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape/model/runners"
	"github.com/ankyra/escape/model/runners/destroy"
)

type DestroyController struct{}

func (DestroyController) Destroy(context Context, destroyBuild, destroyDeployment bool) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Destroy")
	context.Log("destroy.start", nil)
	if destroyBuild {
		runnerContext, err := runners.NewRunnerContext(context)
		if err != nil {
			return err
		}
		if err := destroy.NewDestroyRunner("build").Run(runnerContext); err != nil {
			return err
		}
	}
	if destroyDeployment {
		runnerContext, err := runners.NewRunnerContext(context)
		if err != nil {
			return MarkDeploymentFailed(context, err, state.DestroyFailure)
		}
		if err := destroy.NewDestroyRunner("deploy").Run(runnerContext); err != nil {
			return err
		}
	}
	context.Log("destroy.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (d DestroyController) FetchAndDestroy(context Context, releaseId string, destroyBuild, destroyDeployment bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return MarkDeploymentFailed(context, err, state.DestroyFailure)
	}
	fetcher := FetchController{}
	if err := fetcher.ResolveFetchAndLoad(context, releaseId); err != nil {
		os.Chdir(currentDir)
		return MarkDeploymentFailed(context, err, state.DestroyFailure)
	}
	if err := d.Destroy(context, destroyBuild, destroyDeployment); err != nil {
		os.Chdir(currentDir)
		return err
	}
	return os.Chdir(currentDir)
}
