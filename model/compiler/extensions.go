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
	"github.com/ankyra/escape-client/model/escape_plan"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-core"
)

func (c *Compiler) compileExtensions(plan *escape_plan.EscapePlan) error {
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
			c.metadata.AddConsumes(consume)
		}
		for _, provide := range metadata.GetProvides() {
			c.metadata.AddProvides(provide)
		}
		for _, input := range metadata.GetInputs() {
			c.metadata.AddInputVariable(input)
		}
		for _, output := range metadata.GetOutputs() {
			c.metadata.AddOutputVariable(output)
		}
		for name, newErrand := range metadata.GetErrands() {
			_, exists := c.metadata.Errands[name]
			if exists {
				continue
			}
			newErrand.Script = c.extensionPath(metadata, newErrand.GetScript())
			c.metadata.Errands[name] = newErrand
		}
		for key, val := range metadata.Metadata {
			c.metadata.Metadata[key] = val
		}
		for _, tpl := range metadata.Templates {
			tpl.File = c.extensionPath(metadata, tpl.File)
			tpl.Target = c.extensionPath(metadata, tpl.Target)
			c.metadata.Templates = append(c.metadata.Templates, tpl)
		}
		for name, stage := range metadata.Stages {
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
		for key, val := range metadata.GetVariableContext() {
			metadata, err := c.context.GetDependencyMetadata(val)
			if err != nil {
				return err
			}
			c.VariableCtx[key] = metadata
			c.metadata.SetVariableInContext(key, metadata.GetReleaseId())
		}
		c.VariableCtx[versionlessDep] = metadata
		c.metadata.SetVariableInContext(versionlessDep, metadata.GetReleaseId())
		c.metadata.AddExtension(metadata.GetReleaseId())
	}
	return nil
}

func (c *Compiler) extensionPath(extension *core.ReleaseMetadata, path string) string {
	if path == "" {
		return ""
	}
	return paths.NewPath().ExtensionPath(extension, path)
}
