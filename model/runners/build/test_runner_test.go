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

func (s *testSuite) Test_TestRunner(c *C) {
	runCtx := getRunContext(c, "testdata/test_state.json", "testdata/test_plan.yml")
	c.Assert(NewTestRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_TestRunner_no_test_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/test_state.json", "testdata/plan.yml")
	c.Assert(NewTestRunner().Run(runCtx), IsNil)
	checkStatus(c, runCtx, state.OK)
}

func (s *testSuite) Test_TestRunner_missing_test_file(c *C) {
	runCtx := getRunContext(c, "testdata/test_state.json", "testdata/plan.yml")
	runCtx.GetReleaseMetadata().SetStage("test", "testdata/doesnt_exist.sh")
	c.Assert(NewTestRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.TestFailure)
}

func (s *testSuite) Test_TestRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/test_plan.yml")
	err := NewTestRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Build state '_/name' for release 'name-v0.0.1' could not be found")
	checkStatus(c, runCtx, state.TestFailure)
}

func (s *testSuite) Test_TestRunner_failing_test(c *C) {
	runCtx := getRunContext(c, "testdata/test_state.json", "testdata/failing_test_plan.yml")
	c.Assert(NewTestRunner().Run(runCtx), Not(IsNil))
	checkStatus(c, runCtx, state.TestFailure)
}
