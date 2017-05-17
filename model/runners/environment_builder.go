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
	"fmt"
	"os"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
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

func (e *environmentBuilder) GetInputsForPreStep(ctx RunnerContext, stage string) (*map[string]interface{}, error) {
	calculatedInputs := map[string]interface{}{}
	inputs := ctx.GetDeploymentState().GetPreStepInputs(stage)
	scriptEnv, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return nil, err
	}
	for _, inputVar := range ctx.GetReleaseMetadata().GetInputs() {
		val, err := inputVar.GetValue(inputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		calculatedInputs[inputVar.GetId()] = val
	}
	return prepInputs(ctx, stage, &calculatedInputs)
}

func (e *environmentBuilder) GetInputsForErrand(ctx RunnerContext, errand *core.Errand) (*map[string]interface{}, error) {
	deplState := ctx.GetDeploymentState()
	inputs := deplState.GetCalculatedInputs("deploy")
	result, err := prepInputs(ctx, "deploy", inputs)
	if err != nil {
		return nil, err
	}
	if errand.GetInputs() == nil {
		return result, nil
	}
	scriptEnv, err := ctx.GetScriptEnvironment("deploy")
	if err != nil {
		return nil, err
	}
	for _, inputVar := range errand.GetInputs() {
		val, err := inputVar.GetValue(inputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		(*result)[inputVar.GetId()] = val
	}
	return result, nil
}

func (e *environmentBuilder) GetOutputs(ctx RunnerContext, stage string) (*map[string]interface{}, error) {
	metadata := ctx.GetReleaseMetadata()
	buildOutputs := ctx.GetBuildOutputs()
	result := map[string]interface{}{}
	if len(metadata.GetOutputs()) == 0 {
		if buildOutputs != nil && len(*buildOutputs) > 0 {
			for key, _ := range *buildOutputs {
				fmt.Printf("Warning: received unexpected output variable '%s'\n", key)
			}
		}
		return &result, nil
	}
	scriptEnv, err := ctx.GetScriptEnvironment(stage)
	if err != nil {
		return nil, err
	}
	for _, outputVar := range metadata.GetOutputs() {
		val, err := outputVar.GetValue(buildOutputs, scriptEnv)
		if err != nil {
			return nil, err
		}
		result[outputVar.GetId()] = val
	}
	return &result, nil
}

func addToEnvironmentWithKeyPrefix(env []string, values *map[string]interface{}, prefix string) []string {
	stringValues := util.InterfaceMapToStringMap(values, prefix)
	for key, val := range stringValues {
		envEntry := key + "=" + val
		env = append(env, envEntry)
	}
	return env
}

func (e *environmentBuilder) MergeInputsWithOsEnvironment(ctx RunnerContext) []string {
	result := e.GetEnviron()
	inputs := ctx.GetBuildInputs()
	return addToEnvironmentWithKeyPrefix(result, inputs, "INPUT_")
}

func (e *environmentBuilder) MergeInputsAndOutputsWithOsEnvironment(ctx RunnerContext) []string {
	result := e.GetEnviron()
	inputs := ctx.GetBuildInputs()
	outputs := ctx.GetBuildOutputs()
	result = addToEnvironmentWithKeyPrefix(result, inputs, "INPUT_")
	result = addToEnvironmentWithKeyPrefix(result, outputs, "OUTPUT_")
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

func prepInputs(ctx RunnerContext, stage string, inputs *map[string]interface{}) (*map[string]interface{}, error) {
	metadata := ctx.GetReleaseMetadata()
	deplState := ctx.GetDeploymentState()
	result := map[string]interface{}{}
	for key, val := range metadata.GetMetadata() {
		result["METADATA_"+key] = val
	}
	addValues(&result, deplState.GetCalculatedInputs(stage), "PREVIOUS_")
	addValues(&result, deplState.GetCalculatedOutputs(stage), "PREVIOUS_OUTPUT_")
	addValues(&result, inputs, "")
	return &result, nil
}
