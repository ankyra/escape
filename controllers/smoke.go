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
	"github.com/ankyra/escape/model/runners/deploy"
)

type SmokeController struct{}

func (SmokeController) Smoke(context Context) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetQualifiedReleaseId())
	context.PushLogSection("Smoke tests")
	context.Log("smoke.start", nil)
	runnerContext, err := runners.NewRunnerContext(context)
	if err != nil {
		return MarkDeploymentFailed(context, err, state.TestFailure)
	}
	if err := deploy.NewSmokeRunner().Run(runnerContext); err != nil {
		return err
	}
	context.Log("smoke.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (s SmokeController) FetchAndSmoke(context Context, releaseId string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return MarkDeploymentFailed(context, err, state.TestFailure)
	}
	fetcher := FetchController{}
	if err := fetcher.ResolveFetchAndLoad(context, releaseId); err != nil {
		os.Chdir(currentDir)
		return MarkDeploymentFailed(context, err, state.TestFailure)
	}
	if err := s.Smoke(context); err != nil {
		os.Chdir(currentDir)
		return err
	}
	return os.Chdir(currentDir)
}
