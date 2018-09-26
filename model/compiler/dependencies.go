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
	"fmt"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state"
)

func compileDependencies(ctx *CompilerContext) error {
	ctx.Metadata.Depends = []*core.DependencyConfig{}
	deps, err := ctx.Plan.GetDependencies()
	if err != nil {
		return err
	}
	for _, depend := range deps {
		depend, err := compileDependencyConfig(ctx, depend)
		if err != nil {
			return err
		}
		ctx.Metadata.AddDependency(depend.Copy())
	}
	return nil
}

func compileDependencyConfig(ctx *CompilerContext, depend *core.DependencyConfig) (*core.DependencyConfig, error) {
	if err := depend.EnsureConfigIsParsed(); err != nil {
		return nil, err
	}
	metadata, err := resolveVersion(ctx, depend)
	if err != nil {
		return nil, err
	}
	for _, consume := range metadata.GetConsumerConfig(state.DeployStage) {
		found := false
		for provider, _ := range depend.Consumes {
			if provider == consume.VariableName {
				found = true
				break
			}
		}
		if !found {
			consumerForDep := consume.Copy()
			consumerForDep.SkipActivate = true
			consumerForDep.SkipDeactivate = true
			ctx.Metadata.AddConsumes(consumerForDep)
		}
	}
	for _, input := range metadata.Inputs {
		_, mapped := depend.Mapping[input.Id]
		if mapped {
			continue
		}
		_, mappedBuild := depend.BuildMapping[input.Id]
		if mappedBuild && !depend.InScope(state.DeployStage) {
			continue
		}
		_, mappedDeploy := depend.DeployMapping[input.Id]
		if mappedDeploy && !depend.InScope(state.BuildStage) {
			continue
		}
		if mappedBuild && mappedDeploy {
			continue
		}

		if !input.HasDefault() {
			input = input.Copy()
			input.Scopes = depend.Scopes
			input.EvalBeforeDependencies = true
			ctx.Metadata.AddInputVariable(input)
		}
	}
	if err := depend.Validate(ctx.Metadata); err != nil {
		return nil, err
	}
	ctx.VariableCtx[depend.VariableName] = metadata
	return depend, nil
}

func resolveVersion(ctx *CompilerContext, depCfg *core.DependencyConfig) (*core.ReleaseMetadata, error) {
	if depCfg.NeedsResolving() {
		if ctx.ReleaseQuery == nil {
			return nil, fmt.Errorf("Missing release query function")
		}
		metadata, err := ctx.ReleaseQuery(depCfg)
		if err != nil {
			return nil, err
		}
		depCfg.ReleaseId = depCfg.Project + "/" + depCfg.Name + "-v" + metadata.Version
		depCfg.Version = metadata.Version
	}
	if ctx.DependencyFetcher == nil {
		return nil, fmt.Errorf("Missing dependency fetcher")
	}
	metadata, err := ctx.DependencyFetcher(depCfg)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
