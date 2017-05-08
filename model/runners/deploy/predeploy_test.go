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
)

func (s *testSuite) Test_PreDeployRunner_no_script_defined(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewPreDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_PreDeployRunner_missing_test_file(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/pre_deploy_state.json", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("pre_deploy", "testdata/doesnt_exist.sh")
	err = NewPreDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_PreDeployRunner_missing_deployment_state(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/pre_deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewPreDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing value for variable 'variable'")
}

func (s *testSuite) Test_PreDeployRunner(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/pre_deploy_state.json", "dev", "testdata/pre_deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	err = NewPreDeployRunner().Run(runCtx)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_PreDeployRunner_failing_script(c *C) {
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/pre_deploy_state.json", "dev", "testdata/pre_deploy_plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := runners.NewRunnerContext(ctx)
	c.Assert(err, IsNil)
	ctx.GetReleaseMetadata().SetStage("pre_deploy", "testdata/failing_test.sh")
	err = NewPreDeployRunner().Run(runCtx)
	c.Assert(err, Not(IsNil))
}
