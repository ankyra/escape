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
	"github.com/ankyra/escape-client/model/paths"
)

type PushController struct{}

func (p PushController) Push(context Context, buildFatPackage bool) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Push")
	if err := p.register(context); err != nil {
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

func (p PushController) register(context Context) error {
	context.Log("register.start", nil)
	backend := context.GetEscapeConfig().GetCurrentTarget().GetStorageBackend()
	path := paths.NewPath()
	if err := path.EnsureDependencyCacheDirectoryExists(); err != nil {
		return err
	}
	localRegister := path.LocalReleaseMetadata(context.GetReleaseMetadata())
	if err := context.GetReleaseMetadata().WriteJsonFile(localRegister); err != nil {
		return err
	}
	if backend == "" || backend == "local" {
	} else if backend == "escape" {
		if err := context.GetClient().Register(context.GetReleaseMetadata()); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unsupported Escape storage backend: '%s'", backend)
	}
	context.Log("register.finished", nil)
	return nil
}

func (p PushController) upload(context Context) error {
	context.Log("upload.start", nil)
	backend := context.GetEscapeConfig().GetCurrentTarget().GetStorageBackend()
	releasePath := paths.NewPath().ReleaseLocation(context.GetReleaseMetadata())
	if backend == "" || backend == "local" {
		context.Log("upload.finished", nil)
		return nil
	} else if backend == "escape" {
		return p.uploadToEscapeServer(context, releasePath)
	}
	return fmt.Errorf("Unknown storage backend: '%s'", backend)
}

func (p PushController) uploadToEscapeServer(context Context, releasePath string) error {
	releaseId := context.GetReleaseMetadata().GetReleaseId()
	err := context.GetClient().UploadRelease(releaseId, releasePath)
	if err != nil {
		return err
	}
	context.Log("upload.finished", nil)
	return nil
}
