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

package model

import (
	"errors"

	"github.com/ankyra/escape-client/model/compiler"
	"github.com/ankyra/escape-client/model/escape_plan"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/model/state"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
)

type context struct {
	EscapeConfig       EscapeConfig
	EscapePlan         *escape_plan.EscapePlan
	ReleaseMetadata    *core.ReleaseMetadata
	ProjectState       ProjectState
	EnvironmentState   EnvironmentState
	Logger             util.Logger
	LogConsumers       []util.LogConsumer
	DependencyMetadata map[string]*core.ReleaseMetadata
}

func NewContext() Context {
	ctx := &context{}
	ctx.EscapeConfig = NewEscapeConfig(ctx)
	ctx.Logger = util.NewLogger([]util.LogConsumer{
		util.NewFancyTerminalOutputLogConsumer(),
	})
	ctx.DependencyMetadata = map[string]*core.ReleaseMetadata{}
	return ctx
}

func (c *context) InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation string) error {
	if environment == "" {
		return errors.New("Missing 'environment'")
	}
	if err := c.LoadLocalState(state, environment); err != nil {
		return err
	}
	if err := c.LoadEscapePlan(escapePlanLocation); err != nil {
		return err
	}
	if err := c.LoadMetadata(); err != nil {
		return err
	}
	return nil
}

func (c *context) GetLogger() util.Logger {
	return c.Logger
}
func (c *context) Log(key string, values map[string]string) {
	c.Logger.Log(key, values)
}
func (c *context) PushLogSection(section string) {
	c.Logger.PushSection(section)
}
func (c *context) PushLogRelease(release string) {
	c.Logger.PushRelease(release)
}
func (c *context) PopLogSection() {
	c.Logger.PopSection()
}
func (c *context) PopLogRelease() {
	c.Logger.PopRelease()
}
func (c *context) GetClient() Client {
	return c.EscapeConfig.GetClient()
}
func (c *context) GetEscapePlan() *escape_plan.EscapePlan {
	return c.EscapePlan
}
func (c *context) GetReleaseMetadata() *core.ReleaseMetadata {
	return c.ReleaseMetadata
}
func (c *context) GetProjectState() ProjectState {
	return c.ProjectState
}
func (c *context) GetEnvironmentState() EnvironmentState {
	return c.EnvironmentState
}
func (c *context) GetEscapeConfig() EscapeConfig {
	return c.EscapeConfig
}

func (c *context) GetDependencyMetadata(depReleaseId string) (*core.ReleaseMetadata, error) {
	metadata, ok := c.DependencyMetadata[depReleaseId]
	if ok {
		return metadata, nil
	}
	var err error
	metadata, err = c.FetchDependencyAndReadMetadata(depReleaseId)
	if err != nil {
		return nil, err
	}
	c.DependencyMetadata[depReleaseId] = metadata
	return metadata, nil

}

func (c *context) FetchDependencyAndReadMetadata(depReleaseId string) (*core.ReleaseMetadata, error) {
	c.Log("fetch.start", map[string]string{"dependency": depReleaseId})
	err := DependencyResolver{}.Resolve(c.EscapeConfig, []string{depReleaseId})
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

func (c *context) LoadEscapeConfig(cfgFile, cfgProfile string) error {
	return c.EscapeConfig.LoadConfig(cfgFile, cfgProfile)
}

func (c *context) LoadEscapePlan(cfgFile string) error {
	plan := escape_plan.NewEscapePlan()
	if err := plan.LoadConfig(cfgFile); err != nil {
		return err
	}
	c.EscapePlan = plan
	return nil
}

func (c *context) LoadMetadata() error {
	metadata, err := compiler.NewCompiler().Compile(c)
	if err != nil {
		return err
	}
	c.ReleaseMetadata = metadata
	return nil
}

func (c *context) LoadLocalState(cfgFile, environment string) error {
	p, err := state.NewProjectStateFromFile(cfgFile)
	if err != nil {
		return err
	}
	c.ProjectState = p
	c.EnvironmentState = p.GetEnvironmentStateOrMakeNew(environment)
	if c.EnvironmentState == nil {
		return errors.New("Empty environment state")
	}
	return nil
}
