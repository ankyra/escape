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
	"github.com/ankyra/escape-core/state"
	. "gopkg.in/check.v1"
	"os"
)

func (s *testSuite) Test_PostBuildRunner(c *C) {
	runCtx := getRunContext(c, "testdata/post_build_state.json", "testdata/post_build_plan.yml")
	c.Assert(NewPostBuildRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.RunningPostStep)
}

func (s *testSuite) Test_PostBuildRunner_no_script_defined(c *C) {
	runCtx := getRunContext(c, "testdata/post_build_state.json", "testdata/plan.yml")
	c.Assert(NewPostBuildRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.RunningPostStep)
}

func (s *testSuite) Test_PostBuildRunner_missing_test_file(c *C) {
	runCtx := getRunContext(c, "testdata/post_build_state.json", "testdata/plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_build", "testdata/doesnt_exist.sh")
	c.Assert(NewPostBuildRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_PostBuildRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/post_build_plan.yml")
	err := NewPostBuildRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Build state '_/name' for release 'name-v0.0.1' could not be found")
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_PostBuildRunner_failing_script(c *C) {
	runCtx := getRunContext(c, "testdata/post_build_state.json", "testdata/post_build_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_build", "testdata/failing_test.sh")
	c.Assert(NewPostBuildRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.Failure)
}

func (s *testSuite) Test_PostBuildRunner_default_outputs_dont_calculate(c *C) {
	runCtx := getRunContext(c, "testdata/default_outputs.json", "testdata/default_outputs_plan.yml")
	deploymentState := runCtx.GetDeploymentState()
	deploymentState.UpdateOutputs(Stage, map[string]interface{}{
		"variable": "not test",
	})
	c.Assert(deploymentState.GetCalculatedOutputs(Stage)["variable"], Equals, "not test")
	c.Assert(NewPostBuildRunner().Run(runCtx), IsNil)
	c.Assert(deploymentState.GetCalculatedOutputs(Stage)["variable"], Equals, "not test")
	checkStatus(c, runCtx, state.RunningPostStep)
}
