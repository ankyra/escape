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

package compiler

import (
	"fmt"

	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/variables"
)

func compileInputs(ctx *CompilerContext) error {
	for _, input := range ctx.Plan.Inputs {
		v, err := compileVariable(ctx, input)
		if err != nil {
			return fmt.Errorf("Error compiling input variable: %s", err.Error())
		}
		ctx.Metadata.AddInputVariable(v)
	}
	for _, input := range ctx.Plan.BuildInputs {
		v, err := compileVariable(ctx, input)
		if err != nil {
			return fmt.Errorf("Error compiling 'build_inputs' variable: %s", err.Error())
		}
		v.Scopes = []string{"build"}
		ctx.Metadata.AddInputVariable(v)
	}
	for _, input := range ctx.Plan.DeployInputs {
		v, err := compileVariable(ctx, input)
		if err != nil {
			return fmt.Errorf("Error compiling 'deploy_inputs' variable: %s", err.Error())
		}
		v.Scopes = []string{"deploy"}
		ctx.Metadata.AddInputVariable(v)
	}
	return compileDependencyVariableMapping(ctx)
}

func compileOutputs(ctx *CompilerContext) error {
	for _, output := range ctx.Plan.Outputs {
		v, err := compileVariable(ctx, output)
		if err != nil {
			return fmt.Errorf("Error compiling output variable: %s", err.Error())
		}
		ctx.Metadata.AddOutputVariable(v)
	}
	return nil
}

func compileVariable(ctx *CompilerContext, v interface{}) (result *variables.Variable, err error) {
	result, err = variables.NewVariableFromInterface(v)
	if err != nil {
		return nil, err
	}
	if result.Default != nil {
		return compileDefault(ctx, result)
	}
	return result, nil
}

func compileDefault(ctx *CompilerContext, v *variables.Variable) (*variables.Variable, error) {
	switch v.Default.(type) {
	case int, float64, bool:
		return v, nil
	case string:
		defaultValue := v.Default.(string)
		_, err := script.ParseScript(defaultValue)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse expression '%s' in default field: %s", defaultValue, err.Error())
		}
		str, err := ctx.RunScriptForCompileStep(defaultValue)
		if err == nil {
			v.Default = &str
		}
		return v, nil
	case []interface{}:
		values := []interface{}{}
		for _, k := range v.Default.([]interface{}) {
			switch k.(type) {
			case string:
				_, err := script.ParseScript(k.(string))
				if err != nil {
					return nil, fmt.Errorf("Couldn't parse expression '%s' in default field: %s", k.(string), err.Error())
				}
				str, err := ctx.RunScriptForCompileStep(k.(string))
				if err == nil {
					values = append(values, str)
				} else {
					values = append(values, k)
				}
			}
		}
		v.Default = values
		return v, nil
	}
	return nil, fmt.Errorf("Unexpected type '%T' for default field of variable '%s'", v.Default, v.Id)
}

func compileDependencyVariableMapping(ctx *CompilerContext) error {
	for _, depend := range ctx.Metadata.Depends {
		for _, variable := range ctx.Metadata.Inputs {
			depend.AddVariableMapping(variable.Scopes, variable.Id, "$this.inputs."+variable.Id)
		}
	}
	return nil
}
