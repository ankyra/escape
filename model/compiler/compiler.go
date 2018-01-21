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
	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/util"
)

type CompilerFunc func(*CompilerContext) error

func Compile(plan *escape_plan.EscapePlan,
	reg inventory.Inventory,
	depFetcher func(*core.DependencyConfig) (*core.ReleaseMetadata, error),
	releaseQuery func(*core.Dependency) (*core.ReleaseMetadata, error),
	logger util.Logger) (*core.ReleaseMetadata, error) {

	ctx := NewCompilerContextWithLogger(plan, reg, logger)
	ctx.DependencyFetcher = depFetcher
	ctx.ReleaseQuery = releaseQuery
	compilerSteps := []CompilerFunc{
		compileBasicFields,
		compileConsumers,
		compileGit,
		compileExtensions,
		compileDependencies,
		compileVersion,
		compileMetadata,
		compileScripts,
		compileInputs,
		compileOutputs,
		compileErrands,
		compileTemplates,
		compileIncludes,
		compileLogo,
	}
	for _, step := range compilerSteps {
		if err := step(ctx); err != nil {
			return nil, err
		}
	}
	return ctx.Metadata, ctx.Metadata.Validate()
}
