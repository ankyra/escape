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
	"github.com/ankyra/escape-client/model"
	"github.com/ankyra/escape-client/model/runners"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) Test_DeployRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_DeployRunner_missing_test_file(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/deploy_state.json", "dev", "testdata/deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("post_deploy", "testdata/doesnt_exist.sh")
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_DeployRunner(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/deploy_state.json", "dev", "testdata/deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	deploymentState, err := runCtx.GetEnvironmentState().GetDeploymentState(runCtx.GetDepends())
	c.Assert(err, IsNil)
	deploymentState.SetVersion("deploy", "")
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)

	c.Assert(deploymentState.GetVersion("deploy"), Equals, "0.0.1")
}

func (s *testSuite) Test_DeployRunner_failing_script(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/deploy_state.json", "dev", "testdata/deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("post_deploy", "testdata/failing_test.sh")
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_DeployRunner_variables_are_set_even_if_there_is_no_pre_step(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/deploy_no_pre_step_state.json", "dev", "testdata/deploy_no_pre_step_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	deploymentState, err := runCtx.GetEnvironmentState().GetDeploymentState(runCtx.GetDepends())
	c.Assert(err, IsNil)
	deploymentState.UpdateInputs("deploy", nil)
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), IsNil)
	c.Assert(*deploymentState.GetUserInputs("deploy"), HasLen, 1)
	deploymentState.SetVersion("deploy", "")
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion("deploy"), Equals, "0.0.1")
	c.Assert(*deploymentState.GetCalculatedInputs("deploy"), HasLen, 1)
}

func (s *testSuite) Test_DeployRunner_with_dependencies(c *C) {
	os.Chdir("testdata")
	defer os.Chdir("..")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("deploy_deps_state.json", "dev", "deploy_deps_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	deploymentState, err := runCtx.GetEnvironmentState().GetDeploymentState(runCtx.GetDepends())
	c.Assert(err, IsNil)
	deploymentState.UpdateInputs("deploy", nil)
	deploymentState.UpdateOutputs("deploy", nil)
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), IsNil)
	c.Assert(*deploymentState.GetUserInputs("deploy"), HasLen, 1)
	deploymentState.SetVersion("deploy", "")
	err = NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion("deploy"), Equals, "0.0.1")
	c.Assert(*deploymentState.GetCalculatedInputs("deploy"), HasLen, 1)
	c.Assert(*deploymentState.GetCalculatedOutputs("deploy"), HasLen, 1)
}
