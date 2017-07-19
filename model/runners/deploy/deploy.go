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
	"github.com/ankyra/escape-client/model/state/types"
)

var Stage = "deploy"

func NewPreDeployRunner() Runner {
	return NewPreScriptStepRunner(Stage, "pre_deploy", types.RunningPreStep, types.Failure)
}
func NewMainDeployRunner() Runner {
	return NewMainStepRunner(Stage, "deploy", types.RunningMainStep, types.Failure)
}

func NewPostDeployRunner() Runner {
	return NewPostScriptStepRunner(Stage, "post_deploy", types.RunningPostStep, types.Failure)
}

func NewSmokeRunner() Runner {
	return NewScriptRunner(Stage, "smoke", types.OK, types.TestFailure)
}

func NewDeployRunner() Runner {
	return NewCompoundRunner(
		NewDependencyRunner(Stage, Stage, NewDeployRunner, types.Failure),
		NewPreDeployRunner(),
		NewMainDeployRunner(),
		NewPostDeployRunner(),
		NewStatusCodeRunner(Stage, types.OK),
	)
}
