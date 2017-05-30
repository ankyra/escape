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
)

type PushController struct{}

func (p PushController) Push(context Context, buildFatPackage bool) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Push")
	if err := p.saveLocally(context); err != nil {
		return err
	}
	if err := p.upload(context); err != nil {
		return err
	}
	context.Log("push.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (p PushController) saveLocally(context Context) error {
	path := paths.NewPath()
	metadata := context.GetReleaseMetadata()
	if err := path.EnsureDependencyCacheDirectoryExists(metadata.Project); err != nil {
		return err
	}
	localRegister := path.LocalReleaseMetadata(metadata)
	if err := metadata.WriteJsonFile(localRegister); err != nil {
		return err
	}
	return nil
}

func (p PushController) upload(context Context) error {
	context.Log("upload.start", nil)
	releasePath := paths.NewPath().ReleaseLocation(context.GetReleaseMetadata())
	metadata := context.GetReleaseMetadata()
	project := context.GetEscapeConfig().GetCurrentTarget().GetProject()
	if err := context.GetRegistry().UploadRelease(project, releasePath, metadata); err != nil {
		return err
	}
	context.Log("upload.finished", nil)
	return nil
}
