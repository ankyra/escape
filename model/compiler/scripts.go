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
	"os"
	"path/filepath"
	"strings"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/util"
)

func ScriptFieldError(field string, err error) error {
	return fmt.Errorf("In field %s: %s", field, err.Error())
}

func RelativeScriptOutsideOfBaseDirError(script string) error {
	return fmt.Errorf("The file '%s' is outside of this package's base directory and can't be added.", script)
}

func ScriptDoesNotExistError(script, str string) error {
	return fmt.Errorf("The path to script '%s' does not exist (in '%s')", script, str)
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
	}
	for _, script := range cases {
		if err := setStage(ctx, script[0].(string), script[1]); err != nil {
			return ScriptFieldError(script[0].(string), err)
		}
	}
	return nil
}

func setStage(ctx *CompilerContext, field string, script interface{}) error {
	var stage *core.ExecStage
	metadata := ctx.Metadata

	if script == nil {
		return nil
	}

	switch script.(type) {
	case string:
		str := script.(string)
		if str == "" {
			return nil
		}
		parts := strings.Fields(str)
		firstArg := parts[0]

		// e.g. /bin/ls -al
		if filepath.IsAbs(firstArg) {
			stage = &core.ExecStage{
				Cmd:  firstArg,
				Args: parts[1:],
			}
		} else {
			if util.PathExists(firstArg) {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				path, err := filepath.Abs(firstArg)
				if err != nil {
					return err
				}
				rel, err := filepath.Rel(cwd, path)
				if err != nil {
					return err
				}
				if strings.HasPrefix(rel, "..") {
					return RelativeScriptOutsideOfBaseDirError(firstArg)
				}
				ctx.AddFileDigest(firstArg)
				stage = core.NewExecStageForRelativeScript(str)
			} else if strings.HasPrefix(firstArg, ".") {
				return ScriptDoesNotExistError(firstArg, str)
			} else {
				stage = &core.ExecStage{
					Cmd:  firstArg,
					Args: parts[1:],
				}
			}
		}
	case map[interface{}]interface{}:
		returnedStage, err := core.NewExecStageFromDict(script.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		if returnedStage.RelativeScript != "" {
			ctx.AddFileDigest(returnedStage.RelativeScript)
		}
		stage = returnedStage
	default:
		return fmt.Errorf("Expecting dict or string type. Got '%T'", script)
	}
	if stage == nil {
	}
	metadata.SetExecStage(field, stage)
	return nil
}
