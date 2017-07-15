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
    "os"
	. "gopkg.in/check.v1"
)

func (s *testSuite) Test_NewScriptStep(c *C) {
    runCtx := getRunContext(c, "testdata/helper_state.json", "testdata/helper.yml")
    shouldBeDeployed := true
    step := NewScriptStep(runCtx, "deploy", "pre_build", shouldBeDeployed)
    c.Assert(step.ShouldBeDeployed, Equals, shouldBeDeployed)
    c.Assert(step.Stage, Equals, "deploy")
    c.Assert(step.Step, Equals, "pre_build")
    c.Assert(step.Inputs, IsNil)
    c.Assert(step.LoadOutputs, Equals, shouldBeDeployed)
    c.Assert(step.ScriptPath, Equals, "")
    c.Assert(step.Commit, IsNil)
    c.Assert(step.ModifiesOutputVariables, Equals, false)
}

func (s *testSuite) Test_NewScriptStep_inits_scriptpath(c *C) {
    runCtx := getRunContext(c, "testdata/helper_state.json", "testdata/helper.yml")
    runCtx.GetReleaseMetadata().SetStage("pre_build", "yo.sh")
    step := NewScriptStep(runCtx, "deploy", "pre_build", false)
    c.Assert(step.ScriptPath, Equals, "yo.sh")
}

func (s *testSuite) Test_NewScriptStep_initScript_returns_abs_path(c *C) {
    runCtx := getRunContext(c, "testdata/helper_state.json", "testdata/helper.yml")
    runCtx.GetReleaseMetadata().SetStage("pre_build", "testdata/prebuild.sh")
    step := NewScriptStep(runCtx, "deploy", "pre_build", false)
    scriptPath, err := step.initScript(runCtx)
    c.Assert(err, IsNil)
    cwd, err := os.Getwd()
    c.Assert(err, IsNil)
    c.Assert(scriptPath, Equals, cwd + "/testdata/prebuild.sh")
}

func (s *testSuite) Test_NewScriptStep_initScript_fails_if_script_doesnt_exist(c *C) {
    runCtx := getRunContext(c, "testdata/helper_state.json", "testdata/helper.yml")
    runCtx.GetReleaseMetadata().SetStage("pre_build", "doesnt_exist.sh")
    step := NewScriptStep(runCtx, "deploy", "pre_build", false)
    _, err := step.initScript(runCtx)
    c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_ReadOutputsFromFile(c *C) {
    outputs, err := readOutputsFromFile("testdata/outputs.json")
    c.Assert(err, IsNil)
    c.Assert(outputs, HasLen, 1)
    c.Assert(outputs["output"], Equals, "test")
}

func (s *testSuite) Test_ReadOutputsFromFile_empty_if_file_doesnt_exist(c *C) {
    outputs, err := readOutputsFromFile("testdata/doesnt_exist.json")
    c.Assert(err, IsNil)
    c.Assert(outputs, HasLen, 0)
}

func (s *testSuite) Test_ReadOutputsFromFile_empty_if_file_empty(c *C) {
    outputs, err := readOutputsFromFile("testdata/emptyfile.json")
    c.Assert(err, IsNil)
    c.Assert(outputs, HasLen, 0)
}

func (s *testSuite) Test_ReadOutputsFromFile_fails_if_file_cant_be_read(c *C) {
    os.Chmod("testdata/cantread.json", 0)
    _, err := readOutputsFromFile("testdata/cantread.json")
    c.Assert(err, Not(IsNil))
    os.Chmod("testdata/cantread.json", 0644)
}

func (s *testSuite) Test_ReadOutputsFromFile_fails_if_invalid_json(c *C) {
    _, err := readOutputsFromFile("testdata/invalid.json")
    c.Assert(err, Not(IsNil))
}
