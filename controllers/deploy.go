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
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/model/runners"
	"github.com/ankyra/escape-client/model/runners/deploy"
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"os"
)

type DeployController struct{}

func (d DeployController) Deploy(context Context) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Deploy")
	context.Log("deploy.start", nil)
	runnerContext, err := runners.NewRunnerContext(context, "deploy")
	if err != nil {
		return err
	}
	if err := deploy.NewDeployRunner().Run(runnerContext); err != nil {
		return err
	}
	if err := (SmokeController{}).Smoke(context); err != nil {
		return err
	}
	context.Log("deploy.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (d DeployController) FetchAndDeploy(context Context, releaseId string) error {
	// TODO cd into temp directory
	parsed, err := parsers.ParseQualifiedReleaseId(releaseId)
	if err != nil {
		return err
	}
	if parsed.NeedsResolving() {
		metadata, err := context.QueryReleaseMetadata(parsed.Project, parsed.Name, parsed.Version)
		if err != nil {
			return err
		}
		parsed.Version = metadata.Version
		releaseId = metadata.GetQualifiedReleaseId()
	}
	fetcher := FetchController{}
	if err := fetcher.Fetch(context, []string{releaseId}); err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	root := paths.NewPath().UnpackedDepDirectory(core.NewDependencyFromQualifiedReleaseId(parsed))
	err = os.Chdir(root)
	if err := context.LoadReleaseJson(); err != nil {
		return err
	}
	if err := d.Deploy(context); err != nil {
		os.Chdir(currentDir)
		return err
	}
	return os.Chdir(currentDir)
}
