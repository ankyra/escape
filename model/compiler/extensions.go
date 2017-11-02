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

	"github.com/ankyra/escape/model/paths"
	"github.com/ankyra/escape-core"
)

func compileExtensions(ctx *CompilerContext) error {
	for _, extend := range ctx.Plan.Extends {
		dep, err := core.NewDependencyFromString(extend)
		if err != nil {
			return err
		}
		metadata, err := resolveVersion(ctx, core.NewDependencyConfig(extend), dep)
		if err != nil {
			return err
		}
		for _, consume := range metadata.Consumes {
			ctx.Metadata.AddConsumes(consume)
		}
		for _, provide := range metadata.GetProvides() {
			ctx.Metadata.AddProvides(provide)
		}
		for _, input := range metadata.Inputs {
			ctx.Metadata.AddInputVariable(input)
		}
		for _, output := range metadata.Outputs {
			ctx.Metadata.AddOutputVariable(output)
		}
		for name, newErrand := range metadata.GetErrands() {
			_, exists := ctx.Metadata.Errands[name]
			if exists {
				continue
			}
			newErrand.Script = extensionPath(metadata, newErrand.Script)
			ctx.Metadata.Errands[name] = newErrand
		}
		for key, val := range metadata.Metadata {
			ctx.Metadata.Metadata[key] = val
		}
		for _, tpl := range metadata.Templates {
			tpl.File = extensionPath(metadata, tpl.File)
			tpl.Target = extensionPath(metadata, tpl.Target)
			ctx.Metadata.Templates = append(ctx.Metadata.Templates, tpl)
		}
		for name, stage := range metadata.Stages {
			ctx.Metadata.SetStage(name, extensionPath(metadata, stage.Script))
		}
		for _, d := range metadata.Depends {
			found := false
			deps, err := ctx.Plan.GetDependencies()
			if err != nil {
				return err
			}
			for _, existing := range deps {
				if existing.ReleaseId == d.ReleaseId {
					found = true
				}
			}
			if !found {
				if err := ctx.Plan.AddDependency(d); err != nil {
					return err
				}
			}
		}
		for key, val := range metadata.GetVariableContext() {
			if ctx.DependencyFetcher == nil {
				return fmt.Errorf("Missing dependency fetcher")
			}
			metadata, err := ctx.DependencyFetcher(core.NewDependencyConfig(val))
			if err != nil {
				return err
			}
			ctx.VariableCtx[key] = metadata
			ctx.Metadata.SetVariableInContext(key, metadata.GetQualifiedReleaseId())
		}
		ctx.VariableCtx[dep.Name] = metadata
		ctx.Metadata.SetVariableInContext(dep.Name, metadata.GetQualifiedReleaseId())
		ctx.Metadata.AddExtension(metadata.GetQualifiedReleaseId())
	}
	return nil
}

func extensionPath(extension *core.ReleaseMetadata, path string) string {
	if path == "" {
		return ""
	}
	return paths.NewPath().ExtensionPath(extension, path)
}
