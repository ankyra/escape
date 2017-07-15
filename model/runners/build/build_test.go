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

package build

import (
	"github.com/ankyra/escape-client/model"
	"github.com/ankyra/escape-client/model/runners"
	. "gopkg.in/check.v1"
	"os"
)

func (s *testSuite) Test_BuildRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
    runCtx := getRunContext(c, "testdata/escape_state", "testdata/build_plan.yml")
	c.Assert(NewBuildRunner().Run(runCtx), IsNil)
}

func (s *testSuite) Test_BuildRunner_missing_test_file(c *C) {
    runCtx := getRunContext(c, "testdata/build_state.json", "testdata/build_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_build", "testdata/doesnt_exist.sh")
	c.Assert(NewBuildRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_BuildRunner_failing_script(c *C) {
    runCtx := getRunContext(c, "testdata/build_state.json", "testdata/build_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_build", "testdata/failing_test.sh")
	c.Assert(NewBuildRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_BuildRunner_sets_deployment_status(c *C) {
    runCtx := getRunContext(c, "testdata/build_state.json", "testdata/build_plan.yml")
	c.Assert(NewBuildRunner().Run(runCtx), IsNil)
	deploymentState := runCtx.GetDeploymentState()
	c.Assert(deploymentState.GetVersion("build"), Equals, "0.0.1")
}

func (s *testSuite) Test_BuildRunner_variables_are_set_even_if_there_is_no_pre_step(c *C) {
    runCtx := getRunContext(c, "testdata/build_no_pre_step_state.json", "testdata/build_no_pre_step_plan.yml")

	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateInputs("build", nil)
	c.Assert(deploymentState.GetCalculatedInputs("build"), HasLen, 0)
	c.Assert(deploymentState.GetUserInputs("build"), HasLen, 1)
	deploymentState.CommitVersion("build", runCtx.GetReleaseMetadata())

	c.Assert(NewBuildRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetVersion("build"), Equals, "0.0.1")
	c.Assert(deploymentState.GetCalculatedInputs("build"), HasLen, 1)
}

func (s *testSuite) Test_BuildRunner_has_access_to_previous_outputs(c *C) {
    runCtx := getRunContext(c, "testdata/default_outputs.json", "testdata/default_outputs_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateOutputs("build", map[string]interface{}{
		"variable": "not test",
	})
	c.Assert(deploymentState.GetCalculatedOutputs("build")["variable"], Equals, "not test")
	c.Assert(NewBuildRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetCalculatedOutputs("build")["variable"], Equals, "test")
}

func getRunContext(c *C, stateFile, escapePlan string) runners.RunnerContext {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState(stateFile, "dev", escapePlan)
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx, "build")
	c.Assert(err, IsNil)
    return runCtx
}

