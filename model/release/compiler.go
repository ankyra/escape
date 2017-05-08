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

package release

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/script"
	"github.com/ankyra/escape-client/util"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Compiler struct {
	metadata *releaseMetadata
	context  Context
	// depends:
	// - archive-test-latest as base
	//
	// => VariableCtx["base"] = "archive-test-v1"
	//
	VariableCtx map[string]ReleaseMetadata
	Consumers   map[string]bool
	ResolveType ReleaseTypeResolver
}

func NewCompiler(typeResolver ReleaseTypeResolver) *Compiler {
	return &Compiler{
		VariableCtx: map[string]ReleaseMetadata{},
		Consumers:   map[string]bool{},
		ResolveType: typeResolver,
	}
}

func (c *Compiler) Compile(context Context) (ReleaseMetadata, error) {
	context.PushLogSection("Compile")
	plan := context.GetEscapePlan()
	c.context = context
	c.metadata = NewEmptyReleaseMetadata().(*releaseMetadata)

	c.metadata.Name = plan.GetBuild()
	c.metadata.Type = plan.GetType()
	c.metadata.Description = plan.GetDescription()
	c.metadata.Logo = plan.GetLogo()
	c.metadata.Provides = plan.GetProvides()
	c.metadata.Consumes = plan.GetConsumes()

	if err := c.compileConsumers(plan.GetConsumes()); err != nil {
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
	c.compileIncludes(plan.GetIncludes())
	if err := c.compileLogo(plan.GetLogo()); err != nil {
		return nil, err
	}
	if err := c.compileReleaseTypeExtras(); err != nil {
		return nil, err
	}
	//if build_fat_package:
	//    self._add_dependencies(escape_config, escape_plan)
	context.PopLogSection()
	return c.metadata, nil

}

func (c *Compiler) compileConsumers(consumes []string) error {
	for _, consumer := range consumes {
		c.Consumers[consumer] = true
	}
	c.metadata.Consumes = consumes
	return nil
}

func (c *Compiler) compileDependencies(depends []string) error {

	consumes := map[string]bool{}
	for _, c := range c.metadata.Consumes {
		consumes[c] = true
	}
	result := []string{}
	for _, depend := range depends {
		dep, err := NewDependencyFromString(depend)
		if err != nil {
			return err
		}
		dep = dep.(*dependency)
		if err := dep.ResolveVersion(c.context); err != nil {
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
	version = strings.TrimSpace(version)
	if version == "auto" { // backwards compatibility
		version = "@"
	}
	client := c.context.GetClient()
	plan := c.context.GetEscapePlan()
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
			return errors.New("Not implemented")
		} else {
			return errors.New("Unknown storage backend: " + backend)
		}
	}
	c.metadata.Version = version
	return nil
}

func (c *Compiler) compileEscapePlanScriptDigests(plan EscapePlan) error {
	paths := []string{
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
		variable, err := c.compileVariable(input)
		if err != nil {
			return errors.New("Error compiling input variable: " + err.Error())
		}
		c.metadata.AddInputVariable(variable)
	}
	return nil
}

func (c *Compiler) compileOutputs(outputs []interface{}) error {
	for _, output := range outputs {
		variable, err := c.compileVariable(output)
		if err != nil {
			return errors.New("Error compiling output variable: " + err.Error())
		}
		c.metadata.AddOutputVariable(variable)
	}
	return nil
}

func (c *Compiler) compileErrands(errands map[string]interface{}) error {
	for name, errandDict := range errands {
		errandIface, err := NewErrandFromDict(name, errandDict)
		if err != nil {
			return err
		}
		newErrand := errandIface.(*errand)
		if newErrand.Script == "" {
			return errors.New("Errand " + newErrand.Name + " is missing a script")
		}
		if err := c.compileEscapePlanScriptDigest(newErrand.Script); err != nil {
			return err
		}
		c.metadata.Errands[name] = newErrand
	}
	return nil
}

func (c *Compiler) compileVariable(v interface{}) (Variable, error) {
	var result *variable
	switch v.(type) {
	case string:
		result = NewVariableFromString(v.(string), "string").(*variable)
	case map[interface{}]interface{}:
		resultIface, err := NewVariableFromDict(v.(map[interface{}]interface{}))
		if err != nil {
			return nil, err
		}
		result = resultIface.(*variable)
	default:
		errors.New("Unexpected type")
	}
	if result.Default != nil {
		return c.compileDefault(result)
	}
	return result, nil
}

func (c *Compiler) compileDefault(v *variable) (Variable, error) {
	switch v.Default.(type) {
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

func (c *Compiler) compileReleaseTypeExtras() error {
	plan := c.context.GetEscapePlan()
	releaseType, err := c.ResolveType(c.metadata.GetType())
	if err != nil {
		return err
	}
	return releaseType.CompileMetadata(plan, c.metadata)
}

func RunScriptForCompileStep(scriptStr string, variableCtx map[string]ReleaseMetadata) (string, error) {
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
		return string(v.(int)), nil
	}
	return "", fmt.Errorf("Expression '%s' did not return a string value", scriptStr)
}
