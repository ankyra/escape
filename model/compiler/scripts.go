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

func compileScripts(ctx *CompilerContext) error {
	plan := ctx.Plan
	metadata := ctx.Metadata
	paths := []string{
		plan.GetBuild(), plan.GetDeploy(), plan.GetDestroy(),
		plan.GetPreBuild(), plan.GetPreDeploy(), plan.GetPreDestroy(),
		plan.GetPostBuild(), plan.GetPostDeploy(), plan.GetPostDestroy(),
		plan.GetTest(), plan.GetSmoke(),
	}
	for _, path := range paths {
		if err := ctx.AddFileDigest(path); err != nil {
			return err
		}
	}
	metadata.SetStage("build", plan.GetBuild())
	metadata.SetStage("deploy", plan.GetDeploy())
	metadata.SetStage("destroy", plan.GetDestroy())
	metadata.SetStage("pre_build", plan.GetPreBuild())
	metadata.SetStage("pre_deploy", plan.GetPreDeploy())
	metadata.SetStage("pre_destroy", plan.GetPreDestroy())
	metadata.SetStage("post_build", plan.GetPostBuild())
	metadata.SetStage("post_deploy", plan.GetPostDeploy())
	metadata.SetStage("post_destroy", plan.GetPostDestroy())
	metadata.SetStage("test", plan.GetTest())
	metadata.SetStage("smoke", plan.GetSmoke())
	return nil
}
