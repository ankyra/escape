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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/runners"
	"github.com/ankyra/escape-client/model/runners/build"
)

type BuildController struct{}

func (BuildController) Build(context Context, buildFatPackage bool) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Build")
	context.Log("build.start", nil)
	runnerContext, err := runners.NewRunnerContext(context)
	if err != nil {
		return err
	}
	if err := build.NewBuildRunner().Run(runnerContext); err != nil {
		return err
	}
	if err := context.LoadMetadata(); err != nil {
		return err
	}
	context.Log("build.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}
