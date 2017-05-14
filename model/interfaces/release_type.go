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

package interfaces

import (
	plan "github.com/ankyra/escape-client/model/escape_plan"
)

type ReleaseType interface {
	GetType() string
	InitEscapePlan(*plan.EscapePlan)
	CompileMetadata(*plan.EscapePlan, ReleaseMetadata) error
	Run(RunnerContext) (*map[string]interface{}, error)
	Destroy(RunnerContext) error
}

type ReleaseTypeResolver func(string) (ReleaseType, error)
