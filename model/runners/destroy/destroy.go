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
)

func NewPreDestroyRunner(stage string) Runner {
	return runners.NewScriptRunner(stage, "pre_destroy")
}

func NewDestroyRunner(stage string) Runner {
	return runners.NewCompoundRunner(
		NewPreDestroyRunner(stage),
		NewMainDestroyRunner(stage),
		NewPostDestroyRunner(stage),
	)
}

func NewMainDestroyRunner(stage string) Runner {
	return runners.NewRunner(func(ctx RunnerContext) error {
		step := runners.NewScriptStep(ctx, stage, "destroy", true)
		step.ModifiesOutputVariables = true
		return step.Run(ctx)
	})
}

func NewPostDestroyRunner(stage string) Runner {
	return runners.NewRunner(func(ctx RunnerContext) error {
		step := runners.NewScriptStep(ctx, stage, "post_destroy", true)
		step.Commit = deleteCommit
		return step.Run(ctx)
	})
}

func deleteCommit(ctx RunnerContext, depl DeploymentState, stage string) error {
	deferred := func() Runner { return NewDestroyRunner("deploy") }
	err := runners.NewDependencyRunner("destroy", deferred).Run(ctx)
	if err != nil {
		return err
	}
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
