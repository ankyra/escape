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
	"github.com/ankyra/escape-core"
)

func compileDependencies(ctx *CompilerContext) error {
	result := []string{}
	for _, depend := range ctx.Plan.GetDepends() {
		dep, err := core.NewDependencyFromString(depend)
		if err != nil {
			return err
		}
		metadata, err := resolveVersion(ctx, dep)
		if err != nil {
			return err
		}
		resolvedDep := dep.GetReleaseId()
		versionlessDep := dep.GetVersionlessReleaseId()
		for _, consume := range metadata.GetConsumes() {
			ctx.Metadata.AddConsumes(consume)
		}
		for _, input := range metadata.GetInputs() {
			if !input.HasDefault() {
				ctx.Metadata.AddInputVariable(input)
			}
		}
		ctx.VariableCtx[versionlessDep] = metadata
		ctx.Metadata.SetVariableInContext(versionlessDep, metadata.GetReleaseId())

		if dep.GetVariableName() != "" {
			ctx.VariableCtx[dep.GetVariableName()] = metadata
			ctx.Metadata.SetVariableInContext(dep.GetVariableName(), metadata.GetReleaseId())
		}

		result = append(result, resolvedDep)
	}
	ctx.Metadata.SetDependencies(result)
	return nil
}

func resolveVersion(ctx *CompilerContext, d *core.Dependency) (*core.ReleaseMetadata, error) {
	versionQuery := d.GetVersion()
	if versionQuery != "latest" {
		versionQuery = "v" + versionQuery
	}
	metadata, err := ctx.Registry.QueryReleaseMetadata(d.Project, d.GetName(), versionQuery)
	if err != nil {
		return nil, err
	}
	d.Version = metadata.Version
	return metadata, nil
}
