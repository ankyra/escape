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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/runners"
)

func NewPreDeployRunner() Runner {
	return runners.NewPreScriptStepRunner("deploy", "pre_deploy")
}
func NewMainDeployRunner() Runner {
	return runners.NewMainStepRunner("deploy", "deploy")
}

func NewPostDeployRunner() Runner {
	return runners.NewPostScriptStepRunner("deploy", "post_deploy")
}

func NewSmokeRunner() Runner {
	return runners.NewScriptRunner("deploy", "smoke")
}

func NewDeployRunner() Runner {
	return runners.NewCompoundRunner(
		runners.NewDependencyRunner("deploy", NewDeployRunner),
		NewPreDeployRunner(),
		NewMainDeployRunner(),
		NewPostDeployRunner(),
	)
}
