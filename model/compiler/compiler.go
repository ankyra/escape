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
	"crypto/md5"
	"errors"
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	c.CompilerContext = NewCompilerContext(c.metadata, context.GetEscapePlan())

	if plan.GetName() == "" {
		return nil, fmt.Errorf("Missing build name. Add a 'name' field to your Escape plan")
	}
	c.metadata.Name = plan.GetName()
	c.metadata.Description = plan.GetDescription()
	c.metadata.Logo = plan.GetLogo()
	c.metadata.SetProvides(plan.GetProvides())
	c.metadata.SetConsumes(plan.GetConsumes())

	if err := c.compileExtensions(plan); err != nil {
		return nil, err
	}
	if err := c.compileDependencies(plan.GetDepends()); err != nil {
		return nil, err
	}
	if err := c.compileVersion(plan.GetVersion()); err != nil {
		return nil, err
	}
	if err := c.compileMetadata(plan.Metadata); err != nil {
		return nil, err
	}
	if err := compileScripts(c.CompilerContext); err != nil {
		return nil, err
	}
	if err := c.compileInputs(plan.GetInputs()); err != nil {
		return nil, err
	}
	if err := c.compileOutputs(plan.GetOutputs()); err != nil {
		return nil, err
	}
	if err := c.compileErrands(plan.GetErrands()); err != nil {
		return nil, err
	}
	if err := c.compileTemplates(plan.GetTemplates()); err != nil {
		return nil, err
	}
	c.compileIncludes(plan.GetIncludes())
	if err := c.compileLogo(plan.GetLogo()); err != nil {
		return nil, err
	}
	//if build_fat_package:
	//    self._add_dependencies(escape_config, escape_plan)
	context.PopLogSection()
	return c.metadata, nil

}

func (c *Compiler) ResolveVersion(d *core.Dependency, context Context) error {
	if d.Version != "latest" && !strings.HasSuffix(d.Version, "@") {
		return nil
	}
	project := c.context.GetEscapeConfig().GetCurrentTarget().GetProject()
	versionQuery := d.GetVersion()
	if versionQuery != "latest" {
		versionQuery = "v" + versionQuery
	}
	metadata, err := context.GetRegistry().QueryReleaseMetadata(project, d.GetBuild(), versionQuery)
	if err != nil {
		return err
	}
	d.Version = metadata.Version
	return nil
}

func (c *Compiler) addFileDigest(path string) error {
	if path == "" {
		return nil
	}
	if !util.PathExists(path) {
		return errors.New("File " + path + " was referenced in the escape plan, but it doesn't exist")
	}
	if util.IsDir(path) {
		return c.addDirectoryFileDigests(path)
	}
	f, err := os.Open(path)
	if err != nil {
		return errors.New("Couldn't open " + path + ": " + err.Error())
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	hexDigest := fmt.Sprintf("%x", h.Sum(nil))
	c.metadata.AddFileWithDigest(path, hexDigest)
	return nil
}

func (c *Compiler) addDirectoryFileDigests(path string) error {
	if !util.IsDir(path) {
		return errors.New("Not a directory: " + path)
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.New("Could not read directory " + path + ": " + err.Error())
	}
	for _, fileInfo := range fileInfos {
		target := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			c.addDirectoryFileDigests(target)
		} else {
			c.addFileDigest(target)
		}
	}
	return nil
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
