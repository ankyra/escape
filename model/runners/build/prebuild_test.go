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
	"github.com/ankyra/escape-client/model/state/types"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) Test_PreBuildRunner(c *C) {
	runCtx := getRunContext(c, "testdata/pre_build_state.json", "testdata/pre_build_plan.yml")
	c.Assert(NewPreBuildRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, types.RunningPreStep)
}

func (s *testSuite) Test_PreBuildRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/plan.yml")
	c.Assert(NewPreBuildRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, types.RunningPreStep)
}

func (s *testSuite) Test_PreBuildRunner_missing_test_file(c *C) {
	runCtx := getRunContext(c, "testdata/pre_build_state.json", "testdata/plan.yml")
	runCtx.GetReleaseMetadata().SetStage("pre_build", "testdata/doesnt_exist.sh")
	c.Assert(NewPreBuildRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, types.Failure)
}

func (s *testSuite) Test_PreBuildRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/pre_build_plan.yml")
	err := NewPreBuildRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing value for variable 'variable'")
	checkStatus(c, runCtx, types.Failure)
}

func (s *testSuite) Test_PreBuildRunner_sets_version(c *C) {
	runCtx := getRunContext(c, "testdata/pre_build_state.json", "testdata/pre_build_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.CommitVersion(Stage, runCtx.GetReleaseMetadata())
	c.Assert(NewPreBuildRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetVersion(Stage), Equals, "0.0.1")
	checkStatus(c, runCtx, types.RunningPreStep)
}

func (s *testSuite) Test_PreBuildRunner_failing_script(c *C) {
	runCtx := getRunContext(c, "testdata/pre_build_state.json", "testdata/pre_build_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("pre_build", "testdata/failing_test.sh")
	c.Assert(NewPreBuildRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, types.Failure)
}

func (s *testSuite) Test_PreBuildRunner_propagates_default_updates(c *C) {
	runCtx := getRunContext(c, "testdata/pre_build_default_state.json", "testdata/pre_build_plan.yml")
	runCtx.GetDeploymentState().UpdateInputs(Stage, nil)
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage), HasLen, 0)
	runCtx.GetReleaseMetadata().Stages["pre_build"].Script = ""
	runCtx.GetReleaseMetadata().Inputs[0].Default = "test"

	c.Assert(NewPreBuildRunner().Run(runCtx), IsNil)
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage), HasLen, 1)
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage)["variable"], Equals, "test")

	runCtx.GetReleaseMetadata().Inputs[0].Default = "another test"
	c.Assert(NewPreBuildRunner().Run(runCtx), IsNil)
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage), HasLen, 2)
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage)["variable"], Equals, "another test")
	c.Assert(runCtx.GetDeploymentState().GetCalculatedInputs(Stage)["PREVIOUS_variable"], Equals, "test")
	checkStatus(c, runCtx, types.RunningPreStep)
}
