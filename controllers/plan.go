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

package controllers

import (
	"fmt"
	"github.com/ankyra/escape-client/model/escape_plan"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
)

type PlanController struct{}

func (p PlanController) Compile(context Context) {
	fmt.Println(context.GetReleaseMetadata().ToJson())
}

func (p PlanController) Format(context Context, outputLocation string) error {
	yaml := context.GetEscapePlan().ToYaml()
	fmt.Print(string(yaml))
	if outputLocation != "" {
		return ioutil.WriteFile(outputLocation, yaml, 0644)
	}
	return nil
}

func (p PlanController) Minify(context Context, outputLocation string) error {
	yaml := context.GetEscapePlan().ToMinifiedYaml()
	fmt.Print(string(yaml))
	if outputLocation != "" {
		return ioutil.WriteFile(outputLocation, yaml, 0644)
	}
	return nil
}

func (p PlanController) Init(context Context, build_id, output_file string, force bool) error {
	plan := escape_plan.NewEscapePlan().Init(build_id)
	if util.PathExists(output_file) && !force {
		return fmt.Errorf("'%s' already exists. Use --force / -f to overwrite.", output_file)
	}
	return ioutil.WriteFile(output_file, plan.ToYaml(), 0644)
}
