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

package runners

import (
	"github.com/ankyra/escape-client/model"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) Test_NewRunnerContext_fails_if_metadata_is_missing(c *C) {
	ctx := model.NewContext()
	_, err := NewRunnerContext(ctx)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing metadata in context. This is a bug in Escape.")
}

func (s *testSuite) Test_NewRunnerContext(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(runCtx, Not(IsNil))
	c.Assert(runCtx.GetEnvironmentState(), Equals, ctx.GetEnvironmentState())
	c.Assert(runCtx.GetReleaseMetadata(), Equals, ctx.GetReleaseMetadata())
	c.Assert(runCtx.Logger(), Equals, ctx.GetLogger())
	c.Assert(runCtx.GetDepends(), DeepEquals, []string{"name"})
	c.Assert(runCtx.GetDeploymentState(), IsNil)
}

func (s *testSuite) Test_GetScriptEnvironment_fails_if_deployment_state_is_missing(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(runCtx, Not(IsNil))
	_, err = runCtx.GetScriptEnvironment("deploy")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing deployment state in context. This is a bug in Escape.")
}

func (s *testSuite) Test_GetScriptEnvironment_no_depends(c *C) {
	os.RemoveAll("testdata/escape_state")
	ctx := model.NewContext()
	err := ctx.InitFromLocalEscapePlanAndState("testdata/escape_state", "dev", "testdata/plan.yml")
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx)
	c.Assert(runCtx, Not(IsNil))
	depl, err := ctx.GetEnvironmentState().GetDeploymentState([]string{"name"})
	c.Assert(err, IsNil)
	runCtx.SetDeploymentState(depl)
	scriptEnv, err := runCtx.GetScriptEnvironment("deploy")
	c.Assert(err, IsNil)
	c.Assert(scriptEnv, Not(IsNil))
}
