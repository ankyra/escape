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

package errand

import (
	"os"
	"testing"

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

func (s *testSuite) Test_ErrandRunner(c *C) {
	runCtx := getRunContext(c, "testdata/errand_state.json", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	extraVars := map[string]interface{}{
		"errand_variable": "yo",
	}
	c.Assert(NewErrandRunner(errand, extraVars).Run(runCtx), IsNil)
}

func (s *testSuite) Test_ErrandRunner_no_script_defined(c *C) {
	runCtx := getRunContext(c, "testdata/errand_state.json", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.Script = ""
	c.Assert(NewErrandRunner(errand, nil).Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_ErrandRunner_missing_test_file(c *C) {
	runCtx := getRunContext(c, "testdata/errand_state.json", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.Script = "testdata/doesnt_exist.sh"
	c.Assert(NewErrandRunner(errand, nil).Run(runCtx), Not(IsNil))
}

func (s *testSuite) Test_ErrandRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	runCtx := getRunContext(c, "testdata/escape_state", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	err := NewErrandRunner(errand, nil).Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment state '_/name' for release 'name-v0.0.1' could not be found\n\nYou may need to run `escape run deploy name-v0.0.1` to resolve this issue")
}

func (s *testSuite) Test_ErrandRunner_fails_if_errand_variable_is_missing(c *C) {
	runCtx := getRunContext(c, "testdata/errand_state.json", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	err := NewErrandRunner(errand, nil).Run(runCtx)
	c.Assert(err.Error(), Equals, "Missing value for variable 'errand_variable'")
}

func (s *testSuite) Test_ErrandRunner_failing_script(c *C) {
	runCtx := getRunContext(c, "testdata/errand_state.json", "testdata/errand_plan.yml")
	errand := runCtx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.Script = "testdata/failing_test.sh"
	c.Assert(NewErrandRunner(errand, nil).Run(runCtx), Not(IsNil))
}
