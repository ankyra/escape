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

package interfaces

import (
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model/config"
	plan "github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/util/logger/api"
)

type Context interface {
	InitFromLocalEscapePlanAndState(string, string, string) error
	InitReleaseMetadataByReleaseId(string) error
	GetDependencyMetadata(*core.DependencyConfig) (*core.ReleaseMetadata, error)
	QueryReleaseMetadata(*core.DependencyConfig) (*core.ReleaseMetadata, error)
	LoadEscapeConfig(cfgFile string, cfgProfile string) error
	LoadEscapePlan(string) error
	LoadReleaseJson() error
	CompileEscapePlan() error
	LoadRemoteState(string, string) error
	LoadLocalState(string, string, bool) error
	Log(key string, values map[string]string)
	PushLogRelease(string)
	PushLogSection(string)
	PopLogSection()
	PopLogRelease()

	GetEscapePlan() *plan.EscapePlan
	GetReleaseMetadata() *core.ReleaseMetadata
	GetEnvironmentState() *state.EnvironmentState
	GetEscapeConfig() *config.EscapeConfig
	GetInventory() inventory.Inventory
	SetLogger(api.Logger)
	GetLogger() api.Logger
	GetRootDeploymentName() string
	SetRootDeploymentName(string)
}
