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
	"github.com/ankyra/escape-core/state"
	. "gopkg.in/check.v1"
	"os"
)

func (s *testSuite) Test_SmokeRunner(c *C) {
	runCtx := getRunContext(c, "testdata/smoke_state.json", "testdata/smoke_plan.yml")
	c.Assert(NewSmokeRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_SmokeRunner_no_test_script_defined(c *C) {
	runCtx := getRunContext(c, "testdata/smoke_state.json", "testdata/plan.yml")
	c.Assert(NewSmokeRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_SmokeRunner_missing_smoke_file(c *C) {
	runCtx := getRunContext(c, "testdata/smoke_state.json", "testdata/plan.yml")
	runCtx.GetReleaseMetadata().SetStage("smoke", "testdata/doesnt_exist.sh")
	c.Assert(NewSmokeRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.TestFailure)
}

func (s *testSuite) Test_SmokeRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/smoke_plan.yml")
	err := NewSmokeRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment state '_/name' for release 'name-v0.0.1' could not be found")
	checkStatus(c, runCtx, state.TestFailure)
}

func (s *testSuite) Test_SmokeRunner_failing_test(c *C) {
	runCtx := getRunContext(c, "testdata/smoke_state.json", "testdata/smoke_plan.yml")
	runCtx.GetReleaseMetadata().SetStage("smoke", "testdata/failing_test.sh")
	c.Assert(NewSmokeRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.TestFailure)
}
