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
)

func getRunContext(c *C, stateFile, escapePlan string) RunnerContext {
	ctx := model.NewContext()
	ctx.DisableLogger()
	err := ctx.InitFromLocalEscapePlanAndState(stateFile, "dev", escapePlan)
	c.Assert(err, IsNil)
	runCtx, err := NewRunnerContext(ctx, "deploy")
	c.Assert(err, IsNil)
	return runCtx
}

func (s *testSuite) Test_GetInputsForPreStep(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	inputs, err := NewEmptyEnvEnvironmentBuilder().GetInputsForPreStep(runCtx, "deploy")
	c.Assert(err, IsNil)
	c.Assert(inputs, HasLen, 4)
	c.Assert(inputs["input_variable"], DeepEquals, "testinput")
	c.Assert(inputs["PREVIOUS_input_variable"], DeepEquals, "previous testinput")
	c.Assert(inputs["METADATA_key"], DeepEquals, "value")
	c.Assert(inputs["PREVIOUS_OUTPUT_output_variable"], DeepEquals, "testoutput")
}

func (s *testSuite) Test_GetInputsForPreStep_calculated_inputs(c *C) {
	runCtx := getRunContext(c, "testdata/env_calculated_inputs.json", "testdata/env_calculated_inputs.yml")
	inputs, err := NewEmptyEnvEnvironmentBuilder().GetInputsForPreStep(runCtx, "deploy")
	c.Assert(err, IsNil)
	c.Assert(inputs, HasLen, 5)
	c.Assert(inputs["input_variable"], DeepEquals, "0.0.1")
	c.Assert(inputs["magic_version"], DeepEquals, "0.0.1")
	c.Assert(inputs["override"], DeepEquals, "override")
	c.Assert(inputs["METADATA_key"], DeepEquals, "value")
	c.Assert(inputs["PREVIOUS_OUTPUT_output_variable"], DeepEquals, "testoutput")
}

func (s *testSuite) Test_GetInputsForErrand(c *C) {
	runCtx := getRunContext(c, "testdata/errand.json", "testdata/errand.yml")
	errand := runCtx.GetReleaseMetadata().Errands["my-errand"]
	extraVars := map[string]string{}
	inputs, err := NewEmptyEnvEnvironmentBuilder().GetInputsForErrand(runCtx, errand, extraVars)
	c.Assert(err, IsNil)
	c.Assert(inputs, HasLen, 4)
	c.Assert(inputs["input_variable"], DeepEquals, "override")
	c.Assert(inputs["magic_version"], DeepEquals, "0.0.1")
	c.Assert(inputs["override"], DeepEquals, "not overridden")
	c.Assert(inputs["OUTPUT_output_variable"], DeepEquals, "testoutput")
}

func (s *testSuite) Test_GetInputsForErrand_default_inputs(c *C) {
	runCtx := getRunContext(c, "testdata/errand.json", "testdata/errand.yml")
	errand := runCtx.GetReleaseMetadata().Errands["my-errand"]
	errand.Inputs = nil
	extraVars := map[string]string{}
	inputs, err := NewEmptyEnvEnvironmentBuilder().GetInputsForErrand(runCtx, errand, extraVars)
	c.Assert(err, IsNil)
	c.Assert(inputs, HasLen, 2)
	c.Assert(inputs["input_variable"], DeepEquals, "override")
	c.Assert(inputs["OUTPUT_output_variable"], DeepEquals, "testoutput")
}

func (s *testSuite) Test_GetInputsForErrand_extravars_override(c *C) {
	runCtx := getRunContext(c, "testdata/errand.json", "testdata/errand.yml")
	errand := runCtx.GetReleaseMetadata().Errands["my-errand"]
	extraVars := map[string]string{
		"input_variable": "aight",
		"override":       "yo",
	}
	inputs, err := NewEmptyEnvEnvironmentBuilder().GetInputsForErrand(runCtx, errand, extraVars)
	c.Assert(err, IsNil)
	c.Assert(inputs, HasLen, 4)
	c.Assert(inputs["input_variable"], DeepEquals, "aight")
	c.Assert(inputs["magic_version"], DeepEquals, "0.0.1")
	c.Assert(inputs["override"], DeepEquals, "yo")
	c.Assert(inputs["OUTPUT_output_variable"], DeepEquals, "testoutput")
}

func (s *testSuite) Test_MergeInputsWithOsEnvironment(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	inputs := map[string]interface{}{"input_variable": "yo"}
	runCtx.SetBuildInputs(inputs)

	unit := NewEnvironmentBuilderWithEnv([]string{"test=test"})
	c.Assert(unit.GetEnviron(), DeepEquals, []string{"test=test"})
	env := unit.MergeInputsWithOsEnvironment(runCtx)
	c.Assert(env, HasLen, 2)
	c.Assert(env, DeepEquals, []string{"test=test", "INPUT_input_variable=yo"})
}

func (s *testSuite) Test_MergeInputsAndOutputsWithOsEnvironment(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	inputs := map[string]interface{}{"input_variable": "yo"}
	outputs := map[string]interface{}{"output_variable": "yo"}
	runCtx.SetBuildInputs(inputs)
	runCtx.SetBuildOutputs(outputs)

	unit := NewEnvironmentBuilderWithEnv([]string{"test=test"})
	c.Assert(unit.GetEnviron(), DeepEquals, []string{"test=test"})
	env := unit.MergeInputsAndOutputsWithOsEnvironment(runCtx)
	c.Assert(env, HasLen, 3)
	c.Assert(env, DeepEquals, []string{"test=test", "INPUT_input_variable=yo", "OUTPUT_output_variable=yo"})
}

func (s *testSuite) Test_GetOutputs(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	runCtx.SetBuildOutputs(map[string]interface{}{"output_variable": "test"})
	outputs, err := NewEmptyEnvEnvironmentBuilder().GetOutputs(runCtx, "deploy")
	c.Assert(err, IsNil)
	c.Assert(outputs, HasLen, 1)
	c.Assert(outputs["output_variable"], DeepEquals, "test")
}

func (s *testSuite) Test_GetOutputs_only_warns_if_unexpected_output_is_given(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	runCtx.SetBuildOutputs(map[string]interface{}{"output_variable": "test", "unexpected": "no worries"})
	outputs, err := NewEmptyEnvEnvironmentBuilder().GetOutputs(runCtx, "deploy")
	c.Assert(err, IsNil)
	c.Assert(outputs, HasLen, 1)
	c.Assert(outputs["output_variable"], DeepEquals, "test")
}

func (s *testSuite) Test_GetOutputs_fails_if_outputs_not_set(c *C) {
	runCtx := getRunContext(c, "testdata/env_state.json", "testdata/env_test_plan.yml")
	_, err := NewEmptyEnvEnvironmentBuilder().GetOutputs(runCtx, "deploy")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Missing value for variable 'output_variable'")
}

func (s *testSuite) Test_AddToEnvironmentWithKeyPrefix_empty_values(c *C) {
	env := []string{}
	newEnv := addToEnvironmentWithKeyPrefix(env, nil, "PREFIX_")
	c.Assert(newEnv, DeepEquals, env)
}

func (s *testSuite) Test_AddToEnvironmentWithKeyPrefix_empty_env(c *C) {
	values := map[string]interface{}{
		"test":       "string",
		"other_test": 12,
		"list_test":  []interface{}{"test", "test2"},
	}
	newEnv := addToEnvironmentWithKeyPrefix(nil, values, "PREFIX_")
	c.Assert(newEnv, HasLen, 3)
	var testFound, otherFound, listFound bool
	for _, e := range newEnv {
		if e == "PREFIX_test=string" {
			testFound = true
		} else if e == "PREFIX_other_test=12" {
			otherFound = true
		} else if e == "PREFIX_list_test=[\"test\",\"test2\"]" {
			listFound = true
		}
	}
	c.Assert(testFound, Equals, true)
	c.Assert(otherFound, Equals, true)
	c.Assert(listFound, Equals, true)
}

func (s *testSuite) Test_AddToEnvironmentWithKeyPrefix_unsupported_type(c *C) {
	values := map[string]interface{}{
		"test": map[string]interface{}{},
	}
	c.Assert(func() { addToEnvironmentWithKeyPrefix(nil, values, "PREFIX_") }, PanicMatches,
		`Type '.*' not supported \(key: 'test'\). This is a bug in Escape.`)
}
