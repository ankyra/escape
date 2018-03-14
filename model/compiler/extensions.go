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
	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/paths"
)

func compileExtensions(ctx *CompilerContext) error {
	for _, extend := range ctx.Plan.Extends {
		depCfg := core.NewDependencyConfig(extend)
		if err := depCfg.EnsureConfigIsParsed(); err != nil {
			return err
		}
		metadata, err := resolveVersion(ctx, depCfg)
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
			if stage.IsEmpty() {
				continue
			}
			if stage.RelativeScript != "" {
				script := extensionPath(metadata, stage.RelativeScript)
				stage = core.NewExecStageForRelativeScript(script)
			}
			ctx.Metadata.SetExecStage(name, stage)
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
		if err := depCfg.Validate(metadata); err != nil {
			return err
		}
		ctx.VariableCtx[depCfg.VariableName] = metadata
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
