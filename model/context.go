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

package model

import (
	"errors"

	"github.com/ankyra/escape-core"
	types "github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model/compiler"
	"github.com/ankyra/escape/model/config"
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape/model/state"
	"github.com/ankyra/escape/util/logger/api"
	"github.com/ankyra/escape/util/logger/consumers"
	"github.com/ankyra/escape/util/logger/loggers"
)

type Context struct {
	EscapeConfig       *config.EscapeConfig
	EscapePlan         *escape_plan.EscapePlan
	ReleaseMetadata    *core.ReleaseMetadata
	EnvironmentState   *types.EnvironmentState
	Logger             api.Logger
	LogConsumers       []api.LogConsumer
	DependencyMetadata map[string]*core.ReleaseMetadata
	RootDeploymentName string
}

func NewContext() *Context {
	ctx := &Context{}
	ctx.EscapeConfig = config.NewEscapeConfig()
	ctx.Logger = loggers.NewLogger([]api.LogConsumer{
		consumers.NewFancyTerminalOutputLogConsumer(),
	})
	ctx.DependencyMetadata = map[string]*core.ReleaseMetadata{}
	return ctx
}

func (c *Context) SetLogCollapse(s bool) {
	consumer := consumers.NewFancyTerminalOutputLogConsumer()
	consumer.CollapseSections = s
	c.Logger = loggers.NewLogger([]api.LogConsumer{consumer})
}

func (c *Context) DisableLogger() {
	consumer := consumers.NewNullLogConsumer()
	c.Logger = loggers.NewLogger([]api.LogConsumer{consumer})
}

func (c *Context) InitFromLocalEscapePlanAndState(state, environment, planPath string) error {
	if environment == "" {
		return errors.New("Missing 'environment'")
	}
	useProfileState := false
	if err := c.LoadLocalState(state, environment, useProfileState); err != nil {
		return err
	}
	if err := c.LoadEscapePlan(planPath); err != nil {
		return err
	}
	if err := c.CompileEscapePlan(); err != nil {
		return err
	}
	return nil
}

func (c *Context) InitReleaseMetadataByReleaseId(releaseId string) error {
	dep := core.NewDependencyConfig(releaseId)
	if err := dep.EnsureConfigIsParsed(); err != nil {
		return err
	}
	metadata, err := c.QueryReleaseMetadata(dep)
	if err != nil {
		return err
	}
	c.ReleaseMetadata = metadata
	return nil
}

func (c *Context) GetLogger() api.Logger {
	return c.Logger
}

func (c *Context) Log(key string, values map[string]string) {
	c.Logger.Log(key, values)
}

func (c *Context) PushLogSection(section string) {
	c.Logger.PushSection(section)
}

func (c *Context) PushLogRelease(release string) {
	c.Logger.PushRelease(release)
}

func (c *Context) PopLogSection() {
	c.Logger.PopSection()
}

func (c *Context) PopLogRelease() {
	c.Logger.PopRelease()
}

func (c *Context) GetInventory() inventory.Inventory {
	return c.EscapeConfig.GetInventory()
}

func (c *Context) GetEscapePlan() *escape_plan.EscapePlan {
	return c.EscapePlan
}

func (c *Context) GetReleaseMetadata() *core.ReleaseMetadata {
	return c.ReleaseMetadata
}

func (c *Context) GetEnvironmentState() *types.EnvironmentState {
	return c.EnvironmentState
}

func (c *Context) GetEscapeConfig() *config.EscapeConfig {
	return c.EscapeConfig
}

func (c *Context) GetRootDeploymentName() string {
	if c.RootDeploymentName == "" {
		return c.ReleaseMetadata.GetVersionlessReleaseId()
	}
	return c.RootDeploymentName
}

func (c *Context) SetRootDeploymentName(name string) {
	c.RootDeploymentName = name
}

func (c *Context) QueryReleaseMetadata(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
	metadata, ok := c.DependencyMetadata[dep.ReleaseId]
	if ok {
		return metadata, nil
	}
	metadata, err := c.GetInventory().QueryReleaseMetadata(dep.Project, dep.Name, dep.GetVersionAsString())
	if err != nil {
		return nil, err
	}
	c.DependencyMetadata[dep.ReleaseId] = metadata
	return metadata, nil
}

func (c *Context) GetDependencyMetadata(dep *core.DependencyConfig) (*core.ReleaseMetadata, error) {
	metadata, ok := c.DependencyMetadata[dep.ReleaseId]
	if ok {
		return metadata, nil
	}
	var err error
	metadata, err = c.fetchDependencyAndReadMetadata(dep)
	if err != nil {
		return nil, err
	}
	c.DependencyMetadata[dep.ReleaseId] = metadata
	return metadata, nil
}

func (c *Context) fetchDependencyAndReadMetadata(depCfg *core.DependencyConfig) (*core.ReleaseMetadata, error) {
	depReleaseId := depCfg.ReleaseId
	c.Log("fetch.start", map[string]string{"dependency": depReleaseId})
	err := DependencyResolver{}.Resolve(c.EscapeConfig, []*core.DependencyConfig{depCfg})
	if err != nil {
		return nil, err
	}
	dep, err := core.NewDependencyFromString(depReleaseId)
	if err != nil {
		return nil, err
	}
	unpacked := paths.NewPath().UnpackedDepDirectoryReleaseMetadata(dep)
	c.Log("fetch.finished", map[string]string{"dependency": depReleaseId})
	return core.NewReleaseMetadataFromFile(unpacked)
}

func (c *Context) LoadEscapeConfig(cfgFile, cfgProfile string) error {
	err := c.EscapeConfig.LoadConfig(cfgFile)
	if err != nil {
		return err
	}

	err = c.EscapeConfig.SetActiveProfile(cfgProfile)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) LoadEscapePlan(cfgFile string) error {
	plan := escape_plan.NewEscapePlan()
	if err := plan.LoadConfig(cfgFile); err != nil {
		return err
	}
	c.EscapePlan = plan
	return nil
}

func (c *Context) CompileEscapePlan() error {
	c.PushLogSection("Compile")
	metadata, err := compiler.Compile(
		c.EscapePlan,
		c.GetInventory(),
		c.GetDependencyMetadata,
		c.QueryReleaseMetadata,
		c.Logger,
	)
	if err != nil {
		return err
	}
	c.ReleaseMetadata = metadata
	c.PopLogSection()
	return nil
}

func (c *Context) LoadLocalState(stateFile, environment string, useProfileState bool) error {
	if useProfileState {
		stateFile = c.EscapeConfig.GetCurrentProfile().GetStatePath()
	}
	envState, err := state.NewLocalStateProvider(stateFile).Load("", environment)
	if err != nil {
		return err
	}
	if envState == nil {
		return errors.New("Empty environment state")
	}
	c.EnvironmentState = envState
	return nil
}

func (c *Context) LoadRemoteState(project, environment string) error {
	apiServer := c.EscapeConfig.GetCurrentProfile().GetApiServer()
	escapeToken := c.EscapeConfig.GetCurrentProfile().GetAuthToken()
	insecureSkipVerify := c.EscapeConfig.GetCurrentProfile().GetInsecureSkipVerify()
	envState, err := state.NewRemoteStateProvider(apiServer, escapeToken, insecureSkipVerify).Load(project, environment)
	if err != nil {
		return err
	}
	if envState == nil {
		return errors.New("Empty environment state")
	}
	c.EnvironmentState = envState
	return nil
}

func (c *Context) LoadReleaseJson() error {
	m, err := core.NewReleaseMetadataFromFile("release.json")
	if err != nil {
		return err
	}
	c.ReleaseMetadata = m
	return nil
}
