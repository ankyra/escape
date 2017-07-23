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

package deploy

import (
	. "github.com/ankyra/escape-client/model/runners"
	"github.com/ankyra/escape-core/state"
)

var Stage = "deploy"

func NewPreDeployRunner() Runner {
	return NewPreScriptStepRunner(Stage, "pre_deploy", state.RunningPreStep, state.Failure)
}
func NewMainDeployRunner() Runner {
	return NewMainStepRunner(Stage, "deploy", state.RunningMainStep, state.Failure)
}

func NewPostDeployRunner() Runner {
	return NewPostScriptStepRunner(Stage, "post_deploy", state.RunningPostStep, state.Failure)
}

func NewSmokeRunner() Runner {
	return NewScriptRunner(Stage, "smoke", state.OK, state.TestFailure)
}

func NewDeployRunner() Runner {
	return NewCompoundRunner(
		NewDependencyRunner(Stage, Stage, NewDeployRunner, state.Failure),
		NewPreDeployRunner(),
		NewMainDeployRunner(),
		NewPostDeployRunner(),
		NewStatusCodeRunner(Stage, state.OK),
	)
}
