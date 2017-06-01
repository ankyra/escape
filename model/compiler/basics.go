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
)

func compileBasicFields(ctx *CompilerContext) error {
	if ctx.Plan.Name == "" {
		return fmt.Errorf("Missing build name. Add a 'name' field to your Escape plan")
	}
	ctx.Metadata.Name = ctx.Plan.Name
	ctx.Metadata.Description = ctx.Plan.Description
	ctx.Metadata.Logo = ctx.Plan.Logo
	ctx.Metadata.SetProvides(ctx.Plan.Provides)
	ctx.Metadata.SetConsumes(ctx.Plan.Consumes)
	project := ctx.Project
	if project == "" {
		project = "_"
	}
	ctx.Metadata.Project = project
	return nil
}