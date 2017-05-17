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
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
)

type Context interface {
	InitFromLocalEscapePlanAndState(string, string, string) error
	GetDependencyMetadata(string) (*core.ReleaseMetadata, error)
	FetchDependencyAndReadMetadata(string) (*core.ReleaseMetadata, error)
	LoadEscapeConfig(cfgFile string, cfgProfile string) error
	LoadEscapePlan(string) error
	LoadMetadata() error
	LoadLocalState(string, string) error
	Log(key string, values map[string]string)
	PushLogRelease(string)
	PushLogSection(string)
	PopLogSection()
	PopLogRelease()

	GetEscapePlan() *plan.EscapePlan
	GetReleaseMetadata() *core.ReleaseMetadata
	GetProjectState() ProjectState
	GetEnvironmentState() EnvironmentState
	GetEscapeConfig() EscapeConfig
	GetClient() Client
	GetLogger() util.Logger
}
