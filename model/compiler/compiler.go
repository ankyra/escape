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
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ankyra/escape-client/model/escape_plan"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/variables"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Compiler struct {
	metadata *core.ReleaseMetadata
	context  Context
	// depends:
	// - archive-test-latest as base
	//
	// => VariableCtx["base"] = "archive-test-v1"
	//
	VariableCtx map[string]*core.ReleaseMetadata
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

	if plan.GetName() == "" {
		return nil, fmt.Errorf("Missing build name. Add a 'name' field to your Escape plan")
	}
	c.metadata.Name = plan.GetName()
	c.metadata.Description = plan.GetDescription()
	c.metadata.Logo = plan.GetLogo()
	c.metadata.Provides = plan.GetProvides()
	c.metadata.Consumes = plan.GetConsumes()

	if err := c.compileExtensions(plan); err != nil {
		return nil, err
	}
	if err := c.compileDependencies(plan.GetDepends()); err != nil {
		return nil, err
	}
	if err := c.compileVersion(plan.GetVersion()); err != nil {
		return nil, err
	}
	if err := c.compileMetadata(plan.GetMetadata()); err != nil {
		return nil, err
	}
	if err := c.compileEscapePlanScriptDigests(plan); err != nil {
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

func (c *Compiler) compileExtensions(plan *escape_plan.EscapePlan) error {
	consumes := map[string]bool{}
	provides := map[string]bool{}
	for _, c := range c.metadata.Consumes {
		consumes[c] = true
	}
	for _, c := range c.metadata.Provides {
		provides[c] = true
	}
	for _, extend := range plan.GetExtends() {
		dep, err := core.NewDependencyFromString(extend)
		if err != nil {
			return err
		}
		if err := c.ResolveVersion(dep, c.context); err != nil {
			return err
		}
		resolvedDep := dep.GetReleaseId()
		versionlessDep := dep.GetVersionlessReleaseId()
		metadata, err := c.context.GetDependencyMetadata(resolvedDep)
		if err != nil {
			return err
		}
		for _, consume := range metadata.GetConsumes() {
			if !consumes[consume] {
				consumes[consume] = true
				c.metadata.Consumes = append(c.metadata.Consumes, consume)
			}
		}
		for _, provide := range metadata.GetProvides() {
			if !provides[provide] {
				provides[provide] = true
				c.metadata.Provides = append(c.metadata.Provides, provide)
			}
		}
		for _, input := range metadata.GetInputs() {
			found := false
			for _, i := range c.metadata.GetInputs() {
				if i.GetId() == input.GetId() {
					found = true
					break
				}
			}
			if !found {
				c.metadata.AddInputVariable(input)
			}
		}
		for _, output := range metadata.GetOutputs() {
			found := false
			for _, i := range c.metadata.GetOutputs() {
				if i.GetId() == output.GetId() {
					found = true
					break
				}
			}
			if !found {
				c.metadata.AddOutputVariable(output)
			}
		}
		for name, newErrand := range metadata.GetErrands() {
			_, exists := c.metadata.Errands[name]
			if exists {
				continue
			}
			newErrand.Script = c.extensionPath(metadata, newErrand.GetScript())
			c.metadata.Errands[name] = newErrand
		}
		for key, val := range metadata.GetMetadata() {
			c.metadata.Metadata[key] = val
		}
		for _, tpl := range metadata.GetTemplates() {
			tpl.File = c.extensionPath(metadata, tpl.File)
			tpl.Target = c.extensionPath(metadata, tpl.Target)
			c.metadata.Templates = append(c.metadata.Templates, tpl)
		}
		for name, stage := range metadata.GetStages() {
			c.metadata.SetStage(name, c.extensionPath(metadata, stage.Script))
		}
		for _, d := range metadata.GetDependencies() {
			found := false
			for _, existing := range plan.Depends {
				if existing == d {
					found = true
				}
			}
			if !found {
				plan.Depends = append(plan.Depends, d)
			}
		}
		c.VariableCtx[versionlessDep] = metadata
		c.metadata.SetVariableInContext(versionlessDep, metadata.GetReleaseId())
	}
	return nil
}

func (c *Compiler) extensionPath(extension *core.ReleaseMetadata, path string) string {
	if path == "" {
		return ""
	}
	return paths.NewPath().ExtensionPath(extension, path)
}

func (c *Compiler) ResolveVersion(d *core.Dependency, context Context) error {
	if d.Version != "latest" && !strings.HasSuffix(d.Version, "@") {
		return nil
	}
	backend := context.GetEscapeConfig().GetCurrentTarget().GetStorageBackend()
	if backend == "escape" {
		metadata, err := context.GetClient().ReleaseQuery(d.GetReleaseId())
		if err != nil {
			return err
		}
		d.Version = metadata.GetVersion()
	} else if backend == "" || backend == "local" {
		return errors.New("Backend not implemented: " + backend)
	} else {
		return errors.New("Unsupported Escape storage backend: " + backend)
	}
	return nil
}

func (c *Compiler) compileDependencies(depends []string) error {

	consumes := map[string]bool{}
	for _, c := range c.metadata.Consumes {
		consumes[c] = true
	}
	result := []string{}
	for _, depend := range depends {
		dep, err := core.NewDependencyFromString(depend)
		if err != nil {
			return err
		}
		if err := c.ResolveVersion(dep, c.context); err != nil {
			return err
		}
		resolvedDep := dep.GetReleaseId()
		versionlessDep := dep.GetVersionlessReleaseId()
		metadata, err := c.context.GetDependencyMetadata(resolvedDep)
		if err != nil {
			return err
		}
		for _, consume := range metadata.GetConsumes() {
			if !consumes[consume] {
				consumes[consume] = true
				c.metadata.Consumes = append(c.metadata.Consumes, consume)
			}
		}
		for _, input := range metadata.GetInputs() {
			found := false
			for _, i := range c.metadata.GetInputs() {
				if i.GetId() == input.GetId() {
					found = true
					break
				}
			}
			if !found && !input.HasDefault() {
				c.metadata.AddInputVariable(input)
			}
		}
		c.VariableCtx[versionlessDep] = metadata
		c.metadata.SetVariableInContext(versionlessDep, metadata.GetReleaseId())
		if dep.GetVariableName() != "" {
			c.VariableCtx[dep.GetVariableName()] = metadata
			c.metadata.SetVariableInContext(dep.GetVariableName(), metadata.GetReleaseId())
		}
		result = append(result, resolvedDep)
	}
	c.metadata.Depends = result
	return nil
}

func (c *Compiler) compileVersion(version string) error {
	_, err := script.ParseScript(version)
	if err != nil {
		return fmt.Errorf("Couldn't parse expression '%s' in version field: %s", version, err.Error())
	}
	str, err := RunScriptForCompileStep(version, c.VariableCtx)
	if err != nil {
		return fmt.Errorf("Couldn't evaluate expression '%s' in version field: %s", version, err.Error())
	}
	version = strings.TrimSpace(str)
	if version == "auto" { // backwards compatibility
		version = "@"
	}
	if err := parsers.ValidateVersion(version); err != nil {
		return err
	}
	client := c.context.GetClient()
	plan := c.context.GetEscapePlan()
	plan.SetVersion(version)
	if strings.HasSuffix(version, "@") {
		prefix := version[:len(version)-1]
		backend := c.context.GetEscapeConfig().GetCurrentTarget().GetStorageBackend()
		if backend == "escape" {
			nextVersion, err := client.NextVersionQuery(plan.GetReleaseId(), prefix)
			if err != nil {
				return err
			}
			version = nextVersion
		} else if backend == "" || backend == "local" {
			return fmt.Errorf("Auto versioning backend %s not implemented. The storage backend can be configured in the escape config (see `escape config show`)", backend)
		} else {
			return fmt.Errorf("Unknown storage backend: " + backend)
		}
	}
	c.metadata.Version = version
	return nil
}

func (c *Compiler) compileEscapePlanScriptDigests(plan *escape_plan.EscapePlan) error {
	paths := []string{
		plan.GetBuild(), plan.GetDeploy(), plan.GetDestroy(),
		plan.GetPreBuild(), plan.GetPreDeploy(), plan.GetPreDestroy(),
		plan.GetPostBuild(), plan.GetPostDeploy(), plan.GetPostDestroy(),
		plan.GetTest(), plan.GetPath(),
	}
	for _, path := range paths {
		if err := c.compileEscapePlanScriptDigest(path); err != nil {
			return err
		}
	}
	c.metadata.Path = plan.GetPath()
	c.metadata.SetStage("build", plan.GetBuild())
	c.metadata.SetStage("deploy", plan.GetDeploy())
	c.metadata.SetStage("destroy", plan.GetDestroy())
	c.metadata.SetStage("pre_build", plan.GetPreBuild())
	c.metadata.SetStage("pre_deploy", plan.GetPreDeploy())
	c.metadata.SetStage("pre_destroy", plan.GetPreDestroy())
	c.metadata.SetStage("post_build", plan.GetPostBuild())
	c.metadata.SetStage("post_deploy", plan.GetPostDeploy())
	c.metadata.SetStage("post_destroy", plan.GetPostDestroy())
	c.metadata.SetStage("test", plan.GetTest())
	c.metadata.SetStage("smoke", plan.GetSmoke())
	return nil
}

func (c *Compiler) compileEscapePlanScriptDigest(path string) error {
	if path != "" {
		return c.addFileDigest(path)
	}
	return nil
}

func (c *Compiler) compileIncludes(includes []string) {
	for _, globPattern := range includes {
		paths, err := filepath.Glob(globPattern)
		if err != nil {
			fmt.Println("Warning: ignoring pattern error: " + err.Error())
			continue
		}
		if paths == nil {
			continue
		}
		for _, path := range paths {
			err = c.addFileDigest(path)
			if err != nil {
				fmt.Println("Ignoring problem with path " + path + ": " + err.Error())
			}
		}
	}
}

func (c *Compiler) addFileDigest(path string) error {
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

func (c *Compiler) compileLogo(logo string) error {
	if logo == "" {
		return nil
	}
	if !util.PathExists(logo) {
		return errors.New("Referenced logo " + logo + " does not exist")
	}
	data, err := ioutil.ReadFile(logo)
	if err != nil {
		return errors.New("Couldn't read logo " + logo + ": " + err.Error())
	}
	c.metadata.Logo = base64.StdEncoding.EncodeToString(data)
	return nil
}

func (c *Compiler) compileInputs(inputs []interface{}) error {
	for _, input := range inputs {
		v, err := c.compileVariable(input)
		if err != nil {
			return errors.New("Error compiling input variable: " + err.Error())
		}
		c.metadata.AddInputVariable(v)
	}
	return nil
}

func (c *Compiler) compileOutputs(outputs []interface{}) error {
	for _, output := range outputs {
		v, err := c.compileVariable(output)
		if err != nil {
			return errors.New("Error compiling output variable: " + err.Error())
		}
		c.metadata.AddOutputVariable(v)
	}
	return nil
}

func (c *Compiler) compileErrands(errands map[string]interface{}) error {
	for name, errandDict := range errands {
		newErrand, err := core.NewErrandFromDict(name, errandDict)
		if err != nil {
			return err
		}
		if err := c.compileEscapePlanScriptDigest(newErrand.GetScript()); err != nil {
			return err
		}
		c.metadata.Errands[name] = newErrand
	}
	return nil
}

func (c *Compiler) compileTemplates(templateList []interface{}) error {
	for _, tpl := range templateList {
		template, err := templates.NewTemplateFromInterface(tpl)
		if err != nil {
			return err
		}
		if template.File == "" {
			return fmt.Errorf("Missing 'file' field in template")
		}
		c.addFileDigest(template.File)
		c.metadata.Templates = append(c.metadata.Templates, template)
	}
	return nil
}

func (c *Compiler) compileVariable(v interface{}) (result *variables.Variable, err error) {
	switch v.(type) {
	case string:
		result = variables.NewVariableFromString(v.(string), "string")
	case map[interface{}]interface{}:
		result, err = variables.NewVariableFromDict(v.(map[interface{}]interface{}))
		if err != nil {
			return nil, err
		}
	default:
		errors.New("Unexpected type")
	}
	if result.Default != nil {
		return c.compileDefault(result)
	}
	return result, nil
}

func (c *Compiler) compileDefault(v *variables.Variable) (*variables.Variable, error) {
	switch v.Default.(type) {
	case int, float64, bool:
		return v, nil
	case string:
		defaultValue := v.Default.(string)
		_, err := script.ParseScript(defaultValue)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse expression '%s' in default field: %s", defaultValue, err.Error())
		}
		str, err := RunScriptForCompileStep(defaultValue, c.VariableCtx)
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
				str, err := RunScriptForCompileStep(k.(string), c.VariableCtx)
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

func (c *Compiler) compileMetadata(metadata map[string]string) error {
	result := map[string]string{}
	for key, val := range metadata {
		str, err := RunScriptForCompileStep(val, c.VariableCtx)
		if err != nil {
			return fmt.Errorf("%s in metadata field.", err.Error())
		}
		result[key] = str
	}
	c.metadata.Metadata = result
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
