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

package deploy

import (
	"os"
	"testing"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/model/runners"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type testSuite struct{}

var _ = Suite(&testSuite{})

func getRunContext(c *C, stateFile, escapePlan string) *runners.RunnerContext {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState(stateFile, "dev", escapePlan)
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	return runCtx
}

func checkStatus(c *C, runCtx *runners.RunnerContext, code state.StatusCode) {
	deploymentState := runCtx.GetDeploymentState()
	c.Assert(deploymentState.GetStatus(Stage).Code, Equals, state.StatusCode(code))
}

func (s *testSuite) Test_DeployRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/deploy_plan.yml")
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_DeployRunner_commits_version(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
	c.Assert(runCtx.GetDeploymentState().GetVersion(Stage), Equals, "0.0.1")
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_DeployRunner_failing_pre_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetExecStage("pre_deploy", core.NewExecStageForRelativeScript("testdata/failing_test.sh"))
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_DeployRunner_failing_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetExecStage(Stage, core.NewExecStageForRelativeScript("testdata/failing_test.sh"))
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_DeployRunner_missing_post_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetExecStage("post_deploy", core.NewExecStageForRelativeScript("testdata/doesnt_exist.sh"))
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_DeployRunner_variables_are_set_even_if_there_is_no_pre_step(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_no_pre_step_state.json", "testdata/deploy_no_pre_step_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateInputs(Stage, nil)
	c.Assert(deploymentState.GetCalculatedInputs(Stage), HasLen, 0)
	c.Assert(deploymentState.GetUserInputs(Stage), HasLen, 1)
	deploymentState.CommitVersion(Stage, runCtx.GetReleaseMetadata())
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetVersion(Stage), Equals, "0.0.1")
	c.Assert(deploymentState.GetCalculatedInputs(Stage), HasLen, 1)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_DeployRunner_with_dependencies(c *C) {
	os.Chdir("testdata")
	defer os.Chdir("..")
	runCtx := getRunContext(c, "deploy_deps_state.json", "deploy_deps_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateInputs(Stage, nil)
	deploymentState.UpdateOutputs(Stage, nil)
	c.Assert(deploymentState.GetCalculatedInputs(Stage), HasLen, 0)
	c.Assert(deploymentState.GetUserInputs(Stage), HasLen, 1)
	deploymentState.CommitVersion(Stage, runCtx.GetReleaseMetadata())
	err := NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion(Stage), Equals, "0.0.1")
	c.Assert(deploymentState.GetCalculatedInputs(Stage), HasLen, 1)
	c.Assert(deploymentState.GetCalculatedOutputs(Stage), HasLen, 1)
	checkStatus(c, runCtx, state.OK)
}
