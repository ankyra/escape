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

	core "github.com/ankyra/escape-core"
)

func ScriptFieldError(field string, err error) error {
	return fmt.Errorf("In field %s: %s", field, err.Error())
}

func compileScripts(ctx *CompilerContext) error {
	plan := ctx.Plan
	cases := [][]interface{}{
		[]interface{}{"build", plan.Build},
		[]interface{}{"deploy", plan.Deploy},
		[]interface{}{"destroy", plan.Destroy},
		[]interface{}{"pre_build", plan.PreBuild},
		[]interface{}{"pre_deploy", plan.PreDeploy},
		[]interface{}{"pre_destroy", plan.PreDestroy},
		[]interface{}{"post_build", plan.PostBuild},
		[]interface{}{"post_deploy", plan.PostDeploy},
		[]interface{}{"post_destroy", plan.PostDestroy},
		[]interface{}{"test", plan.Test},
		[]interface{}{"smoke", plan.Smoke},
		[]interface{}{"activate_provider", plan.ActivateProvider},
		[]interface{}{"deactivate_provider", plan.DeactivateProvider},
	}
	for _, script := range cases {
		if err := setStage(ctx, script[0].(string), script[1]); err != nil {
			return ScriptFieldError(script[0].(string), err)
		}
	}
	return nil
}

func setStage(ctx *CompilerContext, field string, script interface{}) error {
	stage, err := core.NewExecStageFromInterface(script)
	if err != nil {
		return err
	}
	if stage == nil {
		return nil
	}
	if stage.RelativeScript != "" {
		ctx.AddFileDigest(stage.RelativeScript)
	}
	ctx.Metadata.SetExecStage(field, stage)
	return nil
}
