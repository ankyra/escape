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

package destroy

import (
	"github.com/ankyra/escape-core/state"
	. "github.com/ankyra/escape/model/runners"
)

func NewDestroyRunner(stage string) Runner {
	return NewCompoundRunner(
		NewPreDestroyRunner(stage),
		NewMainDestroyRunner(stage),
		NewPostDestroyRunner(stage),
	)
}

func NewPreDestroyRunner(stage string) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		deferred := func() Runner { return NewDestroyRunner("deploy") }
		err := NewDependencyRunner("destroy", stage, deferred, state.DestroyFailure).Run(ctx)
		if err != nil {
			return err
		}
		return NewScriptRunner(stage, "pre_destroy", state.RunningPreStep, state.DestroyFailure).Run(ctx)
	})
}

func NewMainDestroyRunner(stage string) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		step := NewScriptStep(ctx, stage, "destroy", true)
		step.ModifiesOutputVariables = true
		return RunOrReportFailure(ctx, stage, step, state.RunningMainDestroyStep, state.DestroyFailure)
	})
}

func NewPostDestroyRunner(stage string) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		step := NewScriptStep(ctx, stage, "post_destroy", true)
		step.Commit = deleteCommit
		return RunOrReportFailure(ctx, stage, step, state.RunningPostDestroyStep, state.DestroyFailure)
	})
}

func deleteCommit(ctx *RunnerContext, depl *state.DeploymentState, stage string) error {
	if err := depl.CommitVersion(stage, ctx.GetReleaseMetadata()); err != nil {
		return err
	}
	if err := depl.UpdateInputs(stage, nil); err != nil {
		return err
	}
	if err := depl.UpdateOutputs(stage, nil); err != nil {
		return err
	}
	return ctx.GetDeploymentState().UpdateStatus(stage, state.NewStatus(state.Empty))
}
