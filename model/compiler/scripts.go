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

import core "github.com/ankyra/escape-core"

func compileScripts(ctx *CompilerContext) error {
	plan := ctx.Plan
	setStage(ctx, "build", plan.Build)
	setStage(ctx, "deploy", plan.Deploy)
	setStage(ctx, "destroy", plan.Destroy)
	setStage(ctx, "pre_build", plan.PreBuild)
	setStage(ctx, "pre_deploy", plan.PreDeploy)
	setStage(ctx, "pre_destroy", plan.PreDestroy)
	setStage(ctx, "post_build", plan.PostBuild)
	setStage(ctx, "post_deploy", plan.PostDeploy)
	setStage(ctx, "post_destroy", plan.PostDestroy)
	setStage(ctx, "test", plan.Test)
	setStage(ctx, "smoke", plan.Smoke)
	return nil
}

func setStage(ctx *CompilerContext, field string, script string) {
	if script == "" {
		return
	}
	metadata := ctx.Metadata
	ctx.AddFileDigest(script)
	stage := core.NewExecStageForRelativeScript(script)
	metadata.SetExecStage(field, stage)
}
