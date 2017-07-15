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

func getRunContext(c *C, stateFile, escapePlan string) runners.RunnerContext {
	ctx := model.NewContext()
	ctx.DisableLogger()
	err := ctx.InitFromLocalEscapePlanAndState(stateFile, "dev", escapePlan)
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx, "deploy")
	c.Assert(err, IsNil)
	return runCtx
}

func (s *testSuite) Test_DeployRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/deploy_plan.yml")
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
}

func (s *testSuite) Test_DeployRunner_commits_version(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
	c.Assert(runCtx.GetDeploymentState().GetVersion("deploy"), Equals, "0.0.1")
}

func (s *testSuite) Test_DeployRunner_failing_pre_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("pre_deploy", "testdata/failing_test.sh")
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_DeployRunner_failing_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("deploy", "testdata/failing_test.sh")
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_DeployRunner_missing_post_deploy_file(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_state.json", "testdata/deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_deploy", "testdata/doesnt_exist.sh")
	c.Assert(NewDeployRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_DeployRunner_variables_are_set_even_if_there_is_no_pre_step(c *C) {
	runCtx := getRunContext(c, "testdata/deploy_no_pre_step_state.json", "testdata/deploy_no_pre_step_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateInputs("deploy", nil)
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), HasLen, 0)
	c.Assert(deploymentState.GetUserInputs("deploy"), HasLen, 1)
	deploymentState.CommitVersion("deploy", runCtx.GetReleaseMetadata())
	c.Assert(NewDeployRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetVersion("deploy"), Equals, "0.0.1")
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), HasLen, 1)
}

func (s *testSuite) Test_DeployRunner_with_dependencies(c *C) {
	os.Chdir("testdata")
	defer os.Chdir("..")
	runCtx := getRunContext(c, "deploy_deps_state.json", "deploy_deps_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateInputs("deploy", nil)
	deploymentState.UpdateOutputs("deploy", nil)
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), HasLen, 0)
	c.Assert(deploymentState.GetUserInputs("deploy"), HasLen, 1)
	deploymentState.CommitVersion("deploy", runCtx.GetReleaseMetadata())
	err := NewDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion("deploy"), Equals, "0.0.1")
	c.Assert(deploymentState.GetCalculatedInputs("deploy"), HasLen, 1)
	c.Assert(deploymentState.GetCalculatedOutputs("deploy"), HasLen, 1)
}
