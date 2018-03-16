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

package compiler

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/script"
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/util"
	"github.com/ankyra/escape/util/logger/api"
)

type CompilerContext struct {
	Metadata          *core.ReleaseMetadata
	Plan              *escape_plan.EscapePlan
	VariableCtx       map[string]*core.ReleaseMetadata
	DependencyFetcher func(*core.DependencyConfig) (*core.ReleaseMetadata, error)
	ReleaseQuery      func(*core.DependencyConfig) (*core.ReleaseMetadata, error)
	Inventory         inventory.Inventory
	Logger            api.Logger
}

func NewCompilerContext(plan *escape_plan.EscapePlan, inventory inventory.Inventory) *CompilerContext {
	return &CompilerContext{
		Metadata:    core.NewEmptyReleaseMetadata(),
		Plan:        plan,
		VariableCtx: map[string]*core.ReleaseMetadata{},
		Inventory:   inventory,
	}
}

func NewCompilerContextWithLogger(plan *escape_plan.EscapePlan, inventory inventory.Inventory, logger api.Logger) *CompilerContext {
	ctx := NewCompilerContext(plan, inventory)
	ctx.Logger = logger
	return ctx
}

func (c *CompilerContext) AddFileDigest(path string) error {
	if path == "" {
		return nil
	}
	if !util.PathExists(path) {
		return fmt.Errorf("File '%s' was referenced in the escape plan, but it doesn't exist", path)
	}
	if util.IsDir(path) {
		return c.addDirectoryFileDigests(path)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Couldn't open '%s': %s", path, err.Error())
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	hexDigest := fmt.Sprintf("%x", h.Sum(nil))
	c.Metadata.AddFileWithDigest(path, hexDigest)
	return nil
}

func (c *CompilerContext) addDirectoryFileDigests(path string) error {
	if !util.IsDir(path) {
		return fmt.Errorf("Not a directory: %s", path)
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Could not read directory '%s': %s", path, err.Error())
	}
	for _, fileInfo := range fileInfos {
		target := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			c.addDirectoryFileDigests(target)
		} else {
			c.AddFileDigest(target)
		}
	}
	return nil
}

func (c *CompilerContext) RunScriptForCompileStep(scriptStr string) (string, error) {
	env := map[string]script.Script{}
	for key, metadata := range c.VariableCtx {
		env[key] = metadata.ToScript()
	}
	scriptEnv := script.NewScriptEnvironmentWithGlobals(env)
	return script.ParseAndEvalToString(scriptStr, scriptEnv)
}
