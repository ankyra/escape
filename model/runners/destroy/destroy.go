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

package destroy

import (
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/runners"
	"github.com/ankyra/escape-client/model/types"
)

func NewPreDestroyRunner(stage string) Runner {
	return runners.NewPreScriptStepRunner(stage, "pre_destroy")
}

func NewDestroyRunner(stage string) Runner {
	deferred := func() Runner { return NewDestroyRunner(stage) }
	return runners.NewCompoundRunner(
		runners.NewDependencyRunner(stage, deferred),
		NewPreDestroyRunner(stage),
		runners.NewRunner(destroyStep),
		NewPostDestroyRunner(stage),
	)
}

func NewPostDestroyRunner(stage string) Runner {
	return runners.NewRunner(func(ctx RunnerContext) error {
		step := runners.NewScriptStep(ctx, stage, "post_destroy", true)
		step.Commit = deleteCommit
		return step.Run(ctx)
	})
}

func destroyStep(ctx RunnerContext) error {
	ctx.Logger().Log("destroy.destroy_step", nil)
	typ, err := types.ResolveType(ctx.GetReleaseMetadata().GetType())
	if err != nil {
		return err
	}
	if err := typ.Destroy(ctx); err != nil {
		return err
	}
	ctx.Logger().Log("destroy.step_finished", nil)
	return nil
}

func deleteCommit(ctx RunnerContext, depl DeploymentState, stage string) error {
	if err := depl.SetVersion(stage, ""); err != nil {
		return err
	}
	if err := depl.UpdateInputs(stage, nil); err != nil {
		return err
	}
	if err := depl.UpdateOutputs(stage, nil); err != nil {
		return err
	}
	return nil
}
