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

package runners

import (
	"errors"
	"fmt"
	"os"
	"testing"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) Test_NewRunnerContext_fails_if_metadata_is_missing(c *C) {
	ctx := model.NewContext()
	_, err := NewRunnerContext(ctx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing metadata in context. This is a bug in Escape.")
}

func (s *testSuite) Test_NewRunnerContext(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	ctx.RootDeploymentName = "test-name"
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(runCtx, Not(IsNil))
	c.Assert(runCtx.GetEnvironmentState(), Equals, ctx.GetEnvironmentState())
	c.Assert(runCtx.GetReleaseMetadata(), Equals, ctx.GetReleaseMetadata())
	c.Assert(runCtx.Logger(), Equals, ctx.GetLogger())
	c.Assert(runCtx.GetRootDeploymentName(), Equals, "test-name")
	c.Assert(runCtx.GetDeploymentState().GetName(), Equals, "test-name")
	c.Assert(runCtx.GetDeploymentState().GetReleaseId("deploy"), Equals, "_/name-v")
}

func (s *testSuite) Test_GetScriptEnvironment_no_depends(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	c.Assert(runCtx, Not(IsNil))
	scriptEnv, err := runCtx.GetScriptEnvironment("deploy")
	c.Assert(err, IsNil)
	c.Assert(scriptEnv, Not(IsNil))
}

func (s *testSuite) Test_NewContextForDependency(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	metadata := core.NewReleaseMetadata("test", "1.0")
	consumers := map[string]string{}

	depl, err := runCtx.deploymentState.GetDeploymentOrMakeNew("build",
		metadata.GetVersionlessReleaseId())

	c.Assert(err, IsNil)
	depRunCtx, err := runCtx.NewContextForDependency("build", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, IsNil)
	c.Assert(depRunCtx.GetEnvironmentState(), Equals, runCtx.environmentState)
	c.Assert(depRunCtx.GetReleaseMetadata(), Equals, metadata)
	c.Assert(depRunCtx.deploymentState.Name, Equals, "_/test")
	c.Assert(depRunCtx.GetRootDeploymentName(), Equals, "_/name")
	c.Assert(depRunCtx.GetRootDeploymentName(), Equals, runCtx.GetRootDeploymentName())
	c.Assert(depRunCtx.deploymentState.GetRootDeploymentStage(), Equals, "build")
	c.Assert(depRunCtx.GetDeploymentState(), Equals, depl)
	c.Assert(depRunCtx.GetPath(), DeepEquals, runCtx.path.NewPathForDependency(metadata))
	c.Assert(depRunCtx.Logger(), Equals, runCtx.logger)
	c.Assert(depRunCtx.context, Equals, ctx)
}

func (s *testSuite) Test_NewContextForDependency_with_consumers(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	consumers := map[string]string{
		"provider1": "otherdepl",
	}
	depRunCtx, err := runCtx.NewContextForDependency("build", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, IsNil)
	c.Assert(depRunCtx.GetDeploymentState().GetProviders("deploy")["provider1"], Equals, "otherdepl")
}

func (s *testSuite) Test_NewContextForDependency_evaluates_consumers(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/consumer_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	runCtx.GetDeploymentState().SetProvider("deploy", "provider1", "somedepl")
	runCtx.toScriptEnvironment = func(d *state.DeploymentState, metadata *core.ReleaseMetadata, stage string, context state.DeploymentResolver) (*script.ScriptEnvironment, error) {
		m := map[string]script.Script{
			"$": script.LiftDict(map[string]script.Script{
				"provider1": script.LiftDict(map[string]script.Script{
					"deployment": script.LiftString("otherdepl"),
				}),
			}),
		}
		env := script.NewScriptEnvironmentFromMap(m)
		return env, nil
	}
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	consumers := map[string]string{
		"provider1": "$provider1.deployment",
	}
	depRunCtx, err := runCtx.NewContextForDependency("deploy", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, IsNil)
	c.Assert(depRunCtx.GetDeploymentState().GetProviders("deploy")["provider1"], Equals, "otherdepl")
}

func (s *testSuite) Test_NewContextForDependency_fails_if_toScriptEnviroment_fails_if_evaluation_fails(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/consumer_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	runCtx.GetDeploymentState().SetProvider("deploy", "provider1", "somedepl")
	runCtx.toScriptEnvironment = func(d *state.DeploymentState, metadata *core.ReleaseMetadata, stage string, context state.DeploymentResolver) (*script.ScriptEnvironment, error) {
		env := script.NewScriptEnvironment()
		return env, nil
	}
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	consumers := map[string]string{
		"provider1": "$provider1.deployment",
	}
	_, err = runCtx.NewContextForDependency("deploy", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, DeepEquals, errors.New("Failed to evaluate '$provider1.deployment': Field '$' was not found in environment."))
}

func (s *testSuite) Test_NewContextForDependency_fails_if_toScriptEnviroment_fails(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/consumer_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	runCtx.GetDeploymentState().SetProvider("deploy", "provider1", "somedepl")
	runCtx.toScriptEnvironment = func(d *state.DeploymentState, metadata *core.ReleaseMetadata, stage string, context state.DeploymentResolver) (*script.ScriptEnvironment, error) {
		return nil, errors.New("Nope")
	}
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	consumers := map[string]string{
		"provider1": "$provider1.deployment",
	}
	_, err = runCtx.NewContextForDependency("deploy", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, DeepEquals, errors.New("Nope"))
}

func (s *testSuite) Test_NewContextForDependency_fails_if_missing_consumer(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	metadata := core.NewReleaseMetadata("test", "1.0")
	metadata.Consumes = []*core.ConsumerConfig{
		core.NewConsumerConfig("provider1"),
	}
	consumers := map[string]string{}
	_, err = runCtx.NewContextForDependency("deploy", metadata.GetVersionlessReleaseId(), metadata, consumers)
	c.Assert(err, Not(IsNil))
	c.Assert(err, DeepEquals, fmt.Errorf("Missing provider of type 'provider1'. This can be configured using the -p / --extra-provider flag."))
}
