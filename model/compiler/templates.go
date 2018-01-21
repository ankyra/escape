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

	"github.com/ankyra/escape-core/templates"
)

func compileTemplates(ctx *CompilerContext) error {
	for _, tpl := range ctx.Plan.Templates {
		template, err := compileTemplate(ctx, tpl, "")
		if err != nil {
			return err
		}
		ctx.Metadata.Templates = append(ctx.Metadata.Templates, template)
	}
	for _, tpl := range ctx.Plan.BuildTemplates {
		template, err := compileTemplate(ctx, tpl, "build")
		if err != nil {
			return err
		}
		ctx.Metadata.Templates = append(ctx.Metadata.Templates, template)
	}
	for _, tpl := range ctx.Plan.DeployTemplates {
		template, err := compileTemplate(ctx, tpl, "deploy")
		if err != nil {
			return err
		}
		ctx.Metadata.Templates = append(ctx.Metadata.Templates, template)
	}
	return nil
}

func compileTemplate(ctx *CompilerContext, tpl interface{}, scope string) (*templates.Template, error) {
	template, err := templates.NewTemplateFromInterface(tpl)
	if err != nil {
		return nil, err
	}
	if template.File == "" {
		return nil, fmt.Errorf("Missing 'file' field in template")
	}
	if scope != "" {
		template.Scopes = []string{scope}
	}
	mapping := template.Mapping
	for _, scope := range template.Scopes {
		for _, i := range ctx.Metadata.GetInputs(scope) {
			_, exists := mapping[i.Id]
			if !exists {
				mapping[i.Id] = "$this.inputs." + i.Id
			}
		}
	}
	extraVars := []string{"branch", "description", "logo", "name",
		"revision", "id", "version", "repository", "release", "versionless_release"}
	for _, v := range extraVars {
		_, exists := mapping[v]
		if !exists {
			mapping[v] = "$this." + v
		}
	}
	ctx.AddFileDigest(template.File)
	return template, nil
}
