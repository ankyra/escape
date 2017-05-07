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
	"github.com/ankyra/escape-client/model/types"
)

type buildRunner struct{}

func NewBuildRunner() Runner {
	return &buildRunner{}
}

func (b *buildRunner) Run(ctx RunnerContext) error {
	if err := NewPreBuildRunner().Run(ctx); err != nil {
		return err
	}
	ctx.Logger().Log("build.build_step", nil)
	typ, err := types.ResolveType(ctx.GetReleaseMetadata().GetType())
	if err != nil {
		return err
	}
	outputs, err := typ.Run(ctx)
	if err != nil {
		return err
	}
	if err := ctx.GetDeploymentState().UpdateOutputs("build", outputs); err != nil {
		return err
	}

	ctx.SetBuildOutputs(outputs)
	if err := NewPostBuildRunner().Run(ctx); err != nil {
		return err
	}
	ctx.Logger().Log("build.build_step_finished", nil)
	return nil
}
