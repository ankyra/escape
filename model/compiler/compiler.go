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
	. "github.com/ankyra/escape-client/model/interfaces"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"strconv"
)

type CompilerFunc func(*CompilerContext) error

type Compiler struct {
	metadata *core.ReleaseMetadata
	context  Context
	// depends:
	// - archive-test-latest as base
	//
	// => VariableCtx["base"] = "archive-test-v1"
	//
	VariableCtx     map[string]*core.ReleaseMetadata
	CompilerContext *CompilerContext
}

func NewCompiler() *Compiler {
	return &Compiler{
		VariableCtx: map[string]*core.ReleaseMetadata{},
	}
}

func (c *Compiler) Compile(context Context) (*core.ReleaseMetadata, error) {
	context.PushLogSection("Compile")
	plan := context.GetEscapePlan()
	c.context = context
	c.metadata = core.NewEmptyReleaseMetadata()
	c.CompilerContext = NewCompilerContext(c.metadata, context.GetEscapePlan(), context.GetRegistry())

	if plan.GetName() == "" {
		return nil, fmt.Errorf("Missing build name. Add a 'name' field to your Escape plan")
	}
	c.metadata.Name = plan.GetName()
	c.metadata.Description = plan.GetDescription()
	c.metadata.Logo = plan.GetLogo()
	c.metadata.SetProvides(plan.GetProvides())
	c.metadata.SetConsumes(plan.GetConsumes())

	if err := compileExtensions(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileDependencies(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileVersion(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileMetadata(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileScripts(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileInputs(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileOutputs(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileErrands(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileTemplates(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileIncludes(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := compileLogo(c.CompilerContext); err != nil {
		return nil, err
	}
	//if build_fat_package:
	//    self._add_dependencies(escape_config, escape_plan)
	context.PopLogSection()
	return c.metadata, nil
}

func RunScriptForCompileStep(scriptStr string, variableCtx map[string]*core.ReleaseMetadata) (string, error) {
	parsedScript, err := script.ParseScript(scriptStr)
	if err != nil {
		return "", err
	}
	env := map[string]script.Script{}
	for key, metadata := range variableCtx {
		env[key] = metadata.ToScript()
	}
	val, err := parsedScript.Eval(script.NewScriptEnvironmentWithGlobals(env))
	if err != nil {
		return "", err
	}
	if val.Type().IsString() {
		v, err := val.Value()
		if err != nil {
			return "", err
		}
		return v.(string), nil
	}
	if val.Type().IsInteger() {
		v, err := val.Value()
		if err != nil {
			return "", err
		}
		return strconv.Itoa(v.(int)), nil
	}
	return "", fmt.Errorf("Expression '%s' did not return a string value", scriptStr)
}
