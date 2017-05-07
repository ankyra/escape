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
)

type predeploy_runner struct {
}

func NewPreDeployRunner() Runner {
	return &predeploy_runner{}
}

func (p *predeploy_runner) Run(ctx RunnerContext) error {
	return runPreScript(ctx, "deploy", "pre_deploy", false)
}
