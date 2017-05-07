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

package types

import (
	. "github.com/ankyra/escape-client/model/interfaces"
)

type AnsibleReleaseType struct{}

func (a *AnsibleReleaseType) GetType() string {
	return "ansible"
}

func (a *AnsibleReleaseType) InitEscapePlan(plan EscapePlan) {
	plan.SetIncludes([]string{
		"defaults/*",
		"files/*",
		"handlers/*",
		"meta/*",
		"tasks/*",
		"templates/*",
		"tests/*",
		"vars/*",
		"README.md",
	})
}

func (a *AnsibleReleaseType) CompileMetadata(plan EscapePlan, release ReleaseMetadata) error {
	return nil
}

func (a *AnsibleReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	return nil, nil
}

func (a *AnsibleReleaseType) Destroy(ctx RunnerContext) error {
	return nil
}
