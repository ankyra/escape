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

package types

import (
	"encoding/json"
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/release"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
	"os"
	"strconv"
)

type TerraformReleaseType struct{}

func (a *TerraformReleaseType) GetType() string {
	return "terraform"
}

func (a *TerraformReleaseType) InitEscapePlan(plan EscapePlan) {
	plan.SetPath(plan.GetBuild() + ".tf")
}

func (a *TerraformReleaseType) CompileMetadata(plan EscapePlan, metadata ReleaseMetadata) error {
	for _, i := range metadata.GetOutputs() {
		found, err := checkExistingVariable(i, "terraform_state", "string")
		if err != nil {
			return err
		}
		if found {
			return nil
		}
	}
	v := release.NewVariableFromString("terraform_state", "string")
	defaultValue := "{}"
	v.SetDefault(&defaultValue)
	v.SetSensitive(true)
	v.SetVisible(false)
	v.SetDescription("Terraform state, managed by Escape")
	metadata.AddOutputVariable(v)
	return nil
}

func (t *TerraformReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	inputs := *ctx.GetBuildInputs()
	env := t.buildEnvironment(ctx.GetReleaseMetadata(), inputs)
	stdout, err := t.runTerraform([]string{"get"}, env, ctx.Logger())
	if err != nil {
		return nil, fmt.Errorf("Could not fetch terraform modules: %s", err.Error())
	}
	stateFile, err := t.writeStateFile(inputs)
	if err != nil {
		return nil, fmt.Errorf("Could not write terraform state file: %s", err.Error())
	}
	stdout, err = t.runTerraform([]string{"apply", "-state", stateFile}, env, ctx.Logger())
	if err != nil {
		return nil, fmt.Errorf("Could not apply terraform state: %s", err.Error())
	}
	stdout, err = t.runTerraform([]string{"output", "-state", stateFile, "-json"}, env, ctx.Logger())
	outputs := map[string]interface{}{}
	if err == nil {
		err := json.Unmarshal([]byte(stdout), &outputs)
		if err != nil {
			return nil, fmt.Errorf("Could not read json outputs: %s", err.Error())
		}
	}
	result := map[string]interface{}{}
	for key, value := range outputs {
		switch value.(type) {
		case map[string]interface{}:
			mapValue := value.(map[string]interface{})
			v, ok := mapValue["value"]
			if ok {
				result[key] = v
			}
		}
	}
	newState, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read terraform state file: %s", err.Error())
	}
	result["terraform_state"] = string(newState)
	return &result, os.Remove(stateFile)
}

func (t *TerraformReleaseType) Destroy(ctx RunnerContext) error {
	inputs_ := ctx.GetBuildInputs()
	var inputs map[string]interface{}
	if inputs_ == nil {
		inputs = map[string]interface{}{}
	} else {
		inputs = *inputs_
	}
	//            applog("destroy.terraform_start", release=release_id)
	stateFile, err := t.writeStateFile(inputs)
	if err != nil {
		return fmt.Errorf("Could not write terraform state file: %s", err.Error())
	}
	env := t.buildEnvironment(ctx.GetReleaseMetadata(), inputs)
	_, err = t.runTerraform([]string{"get"}, env, ctx.Logger())
	if err != nil {
		return fmt.Errorf("Could not fetch terraform modules: %s", err.Error())
	}
	_, err = t.runTerraform([]string{"destroy", "-force", "-state", stateFile}, env, ctx.Logger())
	if err != nil {
		return fmt.Errorf("Could not fetch terraform modules: %s", err.Error())
	}
	_, err = ioutil.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("Could not read terraform state file: %s", err.Error())
	}
	//                environment_state.add_release_inputs_and_outputs(metadata, inputs, outputs, deps)
	//            applog("destroy.terraform_finished", release=release_id)
	return os.Remove(stateFile)
}

func (t *TerraformReleaseType) buildEnvironment(metadata ReleaseMetadata, inputs map[string]interface{}) []string {
	env := os.Environ()
	for k, v := range inputs {
		stringVal := ""
		switch v.(type) {
		case string:
			stringVal = v.(string)
		case float64:
			stringVal = strconv.Itoa(int(v.(float64)))
		case int:
			stringVal = strconv.Itoa(v.(int))
		case []interface{}:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			stringVal = string(jsonBytes)
		default:
			panic(fmt.Sprintf("Unexpected type '%T' for variable %s", v, k))
		}
		env = append(env, "TF_VAR_"+k+"="+stringVal)
	}
	return env
}

func (t *TerraformReleaseType) writeStateFile(inputs map[string]interface{}) (string, error) {
	tmp, err := ioutil.TempFile("", "escape_terraform_state_")
	if err != nil {
		return "", err
	}
	previous, ok := inputs["PREVIOUS_OUTPUT_terraform_state"]
	if !ok {
		return tmp.Name(), os.Remove(tmp.Name())
	}
	switch previous.(type) {
	case string:
		err := ioutil.WriteFile(tmp.Name(), []byte(previous.(string)), 0644)
		if err != nil {
			return "", err
		}
		return tmp.Name(), nil
	}
	return "", fmt.Errorf("Expecting string,")
}

func (t *TerraformReleaseType) runTerraform(args []string, env []string, log Logger) (string, error) {
	cmd := []string{"terraform"}
	cmd = append(cmd, args...)
	p := util.NewProcessRecorder()
	return p.Record(cmd, env, log)
}
