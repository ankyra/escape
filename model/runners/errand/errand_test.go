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

package errand

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

func (s *testSuite) Test_ErrandRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/errand_state.json", "dev", "testdata/errand_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	errand := ctx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.SetScript("")
	err = NewErrandRunner(errand).Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_ErrandRunner_missing_test_file(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/errand_state.json", "dev", "testdata/errand_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	errand := ctx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.SetScript("testdata/doesnt_exist.sh")
	err = NewErrandRunner(errand).Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_ErrandRunner_missing_deployment_state(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/errand_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	errand := ctx.GetReleaseMetadata().GetErrands()["my-errand"]
	err = NewErrandRunner(errand).Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Deployment state 'name' (version 0.0.1) could not be found")
}

func (s *testSuite) Test_ErrandRunner(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/errand_state.json", "dev", "testdata/errand_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	errand := ctx.GetReleaseMetadata().GetErrands()["my-errand"]
	err = NewErrandRunner(errand).Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_ErrandRunner_failing_script(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/errand_state.json", "dev", "testdata/errand_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	errand := ctx.GetReleaseMetadata().GetErrands()["my-errand"]
	errand.SetScript("testdata/failing_test.sh")
	err = NewErrandRunner(errand).Run(runCtx)
	c.Assert(err, Not(IsNil))
}
