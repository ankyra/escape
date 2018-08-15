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

package runners

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape/util"
)

type environmentBuilder struct {
	Environ []string
}

func NewEnvironmentBuilderWithEnv(env []string) *environmentBuilder {
	return &environmentBuilder{
		Environ: env,
	}
}

func NewEnvironmentBuilder() *environmentBuilder {
	return NewEnvironmentBuilderWithEnv(os.Environ())
}
func NewEmptyEnvEnvironmentBuilder() *environmentBuilder {
	return NewEnvironmentBuilderWithEnv([]string{})
}

func (e *environmentBuilder) GetEnviron() []string {
	result := make([]string, len(e.Environ))
	copy(result, e.Environ)
	return result
}

func (e *environmentBuilder) GetInputsForPreStep(ctx *RunnerContext, stage string) (map[string]interface{}, error) {
	calculatedInputs := map[string]interface{}{}
	inputs := ctx.GetDeploymentState().GetPreStepInputs(stage)
	scriptEnv, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return nil, err
	}
	for _, inputVar := range ctx.GetReleaseMetadata().GetInputs(stage) {
		val, err := inputVar.GetValue(&inputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		calculatedInputs[inputVar.Id] = val
	}
	return prepInputs(ctx, stage, &calculatedInputs, false)
}

func (e *environmentBuilder) GetPreDependencyInputs(ctx *RunnerContext, stage string) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}
	for key, val := range ctx.GetDeploymentState().GetUserInputs(stage) {
		inputs[key] = val
	}
	scriptEnv, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return nil, err
	}
	for _, inputVar := range ctx.GetReleaseMetadata().GetInputs(stage) {
		if inputVar.EvalBeforeDependencies {
			val, err := inputVar.GetValue(&inputs, scriptEnv)
			if err != nil {
				return nil, err
			}
			inputs[inputVar.Id] = val
		}
	}
	return inputs, nil
}

func (e *environmentBuilder) GetInputsForDependency(ctx *RunnerContext, parentStage string, mapping map[string]interface{}, parentInputs map[string]interface{}) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}
	scriptEnv, err := ctx.GetScriptEnvironmentForPreDependencyStep(parentStage)
	if err != nil {
		return nil, err
	}
	globals := script.ExpectDictAtom((*scriptEnv)["$"])
	this := script.ExpectDictAtom(globals["this"])
	this["inputs"] = script.ShouldLift(parentInputs)
	for key, val := range mapping {
		strVal, ok := val.(string)
		if !ok {
			inputs[key] = val
		} else {
			val, err := script.ParseAndEvalToGoValue(strVal, scriptEnv)
			if err != nil {
				return nil, err
			}
			inputs[key] = val
		}
	}
	return inputs, nil
}

func (e *environmentBuilder) GetInputsForErrand(ctx *RunnerContext, errand *core.Errand, extraVars map[string]interface{}) (map[string]interface{}, error) {
	deplState := ctx.GetDeploymentState()
	inputs := deplState.GetCalculatedInputs("deploy")
	for key, val := range extraVars {
		inputs[key] = val
	}
	result, err := prepInputs(ctx, "deploy", &inputs, true)
	if err != nil {
		return nil, err
	}
	scriptEnv, err := ctx.GetScriptEnvironment("deploy")
	if err != nil {
		return nil, err
	}
	for _, inputVar := range errand.GetInputs() {
		val, err := inputVar.GetValue(&inputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		result[inputVar.Id] = val
	}
	return result, nil
}

func (e *environmentBuilder) GetOutputs(ctx *RunnerContext, stage string) (map[string]interface{}, error) {
	metadata := ctx.GetReleaseMetadata()
	buildOutputs := ctx.GetBuildOutputs()
	result := map[string]interface{}{}
	scriptEnv, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return nil, err
	}
	outputVariables := metadata.GetOutputs(stage)
	for _, outputVar := range outputVariables {
		val, err := outputVar.GetValue(&buildOutputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		result[outputVar.Id] = val
	}
	for key, _ := range buildOutputs {
		if _, expected := result[key]; !expected {
			fmt.Printf("Warning: received unexpected output variable '%s'\n", key)
		}
	}
	return result, nil
}

func addToEnvironmentWithKeyPrefix(env []string, values map[string]interface{}, prefix string) []string {
	stringValues := util.InterfaceMapToStringMap(&values, prefix)
	for key, val := range stringValues {
		envEntry := key + "=" + val
		env = append(env, envEntry)
	}
	return env
}

func (e *environmentBuilder) GetEscapeEnvironmentVariables() map[string]interface{} {
	return map[string]interface{}{
		"ESCAPE_PLATFORM":  runtime.GOOS,
		"ESCAPE_ARCH":      runtime.GOARCH,
		"PYTHONUNBUFFERED": "1",
	}
}

func (e *environmentBuilder) MergeInputsWithOsEnvironment(ctx *RunnerContext) []string {
	result := e.GetEnviron()
	inputs := ctx.GetBuildInputs()
	escapeEnv := e.GetEscapeEnvironmentVariables()
	result = addToEnvironmentWithKeyPrefix(result, inputs, "INPUT_")
	result = addToEnvironmentWithKeyPrefix(result, escapeEnv, "")
	return result
}

func (e *environmentBuilder) MergeInputsAndOutputsWithOsEnvironment(ctx *RunnerContext) []string {
	result := e.GetEnviron()
	inputs := ctx.GetBuildInputs()
	outputs := ctx.GetBuildOutputs()
	escapeEnv := e.GetEscapeEnvironmentVariables()
	result = addToEnvironmentWithKeyPrefix(result, inputs, "INPUT_")
	result = addToEnvironmentWithKeyPrefix(result, outputs, "OUTPUT_")
	result = addToEnvironmentWithKeyPrefix(result, escapeEnv, "")
	return result
}

func addValues(result, values *map[string]interface{}, prefix string) {
	if values == nil {
		return
	}
	for key, val := range *values {
		if !strings.HasPrefix(key, "PREVIOUS_") {
			(*result)[prefix+key] = val
		}
	}
}

func prepInputs(ctx *RunnerContext, stage string, inputs *map[string]interface{}, isErrand bool) (map[string]interface{}, error) {
	metadata := ctx.GetReleaseMetadata()
	deplState := ctx.GetDeploymentState()
	result := map[string]interface{}{}
	for key, val := range metadata.Metadata {
		result["METADATA_"+key] = val
	}
	calcInputs := deplState.GetCalculatedInputs(stage)
	calcOutputs := deplState.GetCalculatedOutputs(stage)
	if isErrand {
		addValues(&result, &calcInputs, "")
		addValues(&result, &calcOutputs, "OUTPUT_")
	} else {
		addValues(&result, &calcInputs, "PREVIOUS_")
		addValues(&result, &calcOutputs, "PREVIOUS_OUTPUT_")
	}
	addValues(&result, inputs, "")
	return result, nil
}
