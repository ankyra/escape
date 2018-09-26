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

package controllers

import (
	"fmt"
	"io/ioutil"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/ankyra/escape/util"
)

type PlanController struct{}

func (p PlanController) Compile(context *model.Context) {
	fmt.Println(context.GetReleaseMetadata().ToJson())
}

func (p PlanController) Diff(context *model.Context) error {
	metadata := context.GetReleaseMetadata()
	inventory := context.GetInventory()
	previous, err := inventory.QueryReleaseMetadata(metadata.Project, metadata.Name, "latest")
	if err != nil {
		return fmt.Errorf(`Can't show differences with previous version of "%s", because no other version of "%[1]s" exists in the Inventory (%s)`, context.GetRootDeploymentName(), context.GetEscapeConfig().GetCurrentProfile().ApiServer)
	}
	if previous == nil {
		return fmt.Errorf("No previous versions found")
	}
	for _, change := range core.Diff(previous, metadata) {
		fmt.Println(change.ToString())
	}
	return nil
}

func (p PlanController) Format(context *model.Context, outputLocation string) error {
	yaml := context.GetEscapePlan().ToYaml()
	fmt.Print(string(yaml))
	if outputLocation != "" {
		return ioutil.WriteFile(outputLocation, yaml, 0644)
	}
	return nil
}

func (p PlanController) Minify(context *model.Context, outputLocation string) error {
	yaml := context.GetEscapePlan().ToMinifiedYaml()
	fmt.Print(string(yaml))
	if outputLocation != "" {
		return ioutil.WriteFile(outputLocation, yaml, 0644)
	}
	return nil
}

func (p PlanController) Init(context *model.Context, build_id, output_file string, force, minify bool) error {
	if util.PathExists(output_file) && !force {
		return fmt.Errorf("'%s' already exists. Use --force / -f to overwrite.", output_file)
	}
	plan := escape_plan.NewEscapePlan().Init(build_id)
	if minify {
		return ioutil.WriteFile(output_file, plan.ToMinifiedYaml(), 0644)
	}
	return ioutil.WriteFile(output_file, plan.ToInitTemplate(), 0644)
}

func (p PlanController) Get(context *model.Context, field string) error {
	var output string
	switch field {
	case "name":
		output = context.GetEscapePlan().Name
	case "version":
		output = context.GetEscapePlan().Version
	case "description":
		output = context.GetEscapePlan().Description
	case "logo":
		output = context.GetEscapePlan().Logo
	default:
		return fmt.Errorf("This field is currently unsupported by this command.")
	}

	fmt.Println(output)
	return nil

}
