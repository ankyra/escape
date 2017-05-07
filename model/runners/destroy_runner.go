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

package runners

import (
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/release"
	"github.com/ankyra/escape-client/model/types"
	"os"
)

type destroyRunner struct {
	Stage string
}

func NewDestroyRunner(stage string) Runner {
	return &destroyRunner{
		Stage: stage,
	}
}

func (b *destroyRunner) Run(ctx RunnerContext) error {
	if err := b.runDependencies(ctx); err != nil {
		return err
	}
	//        self._check_file_integrity(metadata.files)
	if err := NewPreDestroyRunner(b.Stage).Run(ctx); err != nil {
		return err
	}
	ctx.Logger().Log("destroy.destroy_step", nil)
	typ, err := types.ResolveType(ctx.GetReleaseMetadata().GetType())
	if err != nil {
		return err
	}
	if err := typ.Destroy(ctx); err != nil {
		return err
	}
	if err := NewPostDestroyRunner(b.Stage).Run(ctx); err != nil {
		return err
	}
	ctx.Logger().Log("destroy.step_finished", nil)
	return nil
}

func (b *destroyRunner) runDependencies(ctx RunnerContext) error {
	metadata := ctx.GetReleaseMetadata()
	for _, depend := range metadata.GetDependencies() {
		if err := b.runDependency(ctx, depend); err != nil {
			return err
		}
	}
	return nil
}

func (b *destroyRunner) runDependency(ctx RunnerContext, dependency string) error {
	ctx.Logger().PushSection("Dependency " + dependency)
	ctx.Logger().Log("destroy.destroy_dependency", map[string]string{
		"dependency": dependency,
	})
	ctx.Logger().PushRelease(dependency)
	dep, err := release.NewDependencyFromString(dependency)
	if err != nil {
		return err
	}
	location := ctx.GetPath().UnpackedDepDirectory(dep)
	metadata, err := newMetadataFromReleaseDir(location)
	if err != nil {
		return err
	}
	depCtx := ctx.NewContextForDependency(metadata)
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(location); err != nil {
		return err
	}
	if err := NewDestroyRunner(b.Stage).Run(depCtx); err != nil {
		return err
	}
	if err := os.Chdir(currentDir); err != nil {
		return err
	}
	ctx.Logger().Log("destroy.destroy_dependency_finished", nil)
	ctx.Logger().PopRelease()
	ctx.Logger().PopSection()
	return nil
}
