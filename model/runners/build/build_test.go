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
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/build_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewBuildRunner().Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_BuildRunner_missing_test_file(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/build_state.json", "dev", "testdata/build_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("post_build", "testdata/doesnt_exist.sh")
	err = NewBuildRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_BuildRunner(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/build_state.json", "dev", "testdata/build_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewBuildRunner().Run(runCtx)
	c.Assert(err, IsNil)

	deploymentState, err := runCtx.GetEnvironmentState().GetDeploymentState(runCtx.GetDepends())
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion("build"), Equals, "0.0.1")
}

func (s *testSuite) Test_BuildRunner_failing_script(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/build_state.json", "dev", "testdata/build_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("post_build", "testdata/failing_test.sh")
	err = NewBuildRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_BuildRunner_variables_are_set_even_if_there_is_no_pre_step(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/build_no_pre_step_state.json", "dev", "testdata/build_no_pre_step_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	deploymentState, err := runCtx.GetEnvironmentState().GetDeploymentState(runCtx.GetDepends())
	c.Assert(err, IsNil)
	deploymentState.UpdateInputs("build", nil)
	c.Assert(deploymentState.GetCalculatedInputs("build"), IsNil)
	c.Assert(*deploymentState.GetUserInputs("build"), HasLen, 1)
	deploymentState.SetVersion("build", "")
	err = NewBuildRunner().Run(runCtx)
	c.Assert(err, IsNil)
	c.Assert(deploymentState.GetVersion("build"), Equals, "0.0.1")
	c.Assert(*deploymentState.GetCalculatedInputs("build"), HasLen, 1)
}
