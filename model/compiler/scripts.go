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

func compileScripts(ctx *CompilerContext) error {
	plan := ctx.Plan
	metadata := ctx.Metadata
	paths := []string{
		plan.Build, plan.Deploy, plan.Destroy,
		plan.PreBuild, plan.PreDeploy, plan.PreDestroy,
		plan.PostBuild, plan.PostDeploy, plan.PostDestroy,
		plan.Test, plan.Smoke,
	}
	for _, path := range paths {
		if err := ctx.AddFileDigest(path); err != nil {
			return err
		}
	}
	metadata.SetStage("build", plan.Build)
	metadata.SetStage("deploy", plan.Deploy)
	metadata.SetStage("destroy", plan.Destroy)
	metadata.SetStage("pre_build", plan.PreBuild)
	metadata.SetStage("pre_deploy", plan.PreDeploy)
	metadata.SetStage("pre_destroy", plan.PreDestroy)
	metadata.SetStage("post_build", plan.PostBuild)
	metadata.SetStage("post_deploy", plan.PostDeploy)
	metadata.SetStage("post_destroy", plan.PostDestroy)
	metadata.SetStage("test", plan.Test)
	metadata.SetStage("smoke", plan.Smoke)
	return nil
}
