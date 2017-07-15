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
	. "gopkg.in/check.v1"
	"os"
)

func (s *testSuite) Test_PostDeployRunner(c *C) {
	runCtx := getRunContext(c, "testdata/post_deploy_state.json", "testdata/post_deploy_plan.yml")
	outputs := map[string]interface{}{"output_variable": "testinput"}
	runCtx.GetDeploymentState().UpdateOutputs("deploy", outputs)
	c.Assert(NewPostDeployRunner().Run(runCtx), IsNil)
}

func (s *testSuite) Test_PostDeployRunner_no_script_defined(c *C) {
	runCtx := getRunContext(c, "testdata/post_deploy_state.json", "testdata/plan.yml")
	c.Assert(NewPostDeployRunner().Run(runCtx), IsNil)
}

func (s *testSuite) Test_PostDeployRunner_missing_test_file(c *C) {
	runCtx := getRunContext(c, "testdata/post_deploy_state.json", "testdata/plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_deploy", "testdata/doesnt_exist.sh")
	c.Assert(NewPostDeployRunner().Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_PostDeployRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/post_deploy_plan.yml")
	err := NewPostDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment state '_/name' for release 'name-v0.0.1' could not be found")
}

func (s *testSuite) Test_PostDeployRunner_failing_script(c *C) {
	runCtx := getRunContext(c, "testdata/post_deploy_state.json", "testdata/post_deploy_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("post_deploy", "testdata/failing_test.sh")
	c.Assert(NewPostDeployRunner().Run(runCtx), Not(IsNil))
}
