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
	"github.com/ankyra/escape-core/templates"
)

func compileTemplates(ctx *CompilerContext) error {
	for _, tpl := range ctx.Plan.Templates {
		template, err := templates.NewTemplateFromInterface(tpl)
		if err != nil {
			return err
		}
		if template.File == "" {
			return fmt.Errorf("Missing 'file' field in template")
		}
		mapping := template.Mapping
		for _, i := range ctx.Metadata.GetInputs() {
			_, exists := mapping[i.Id]
			if !exists {
				mapping[i.Id] = "$this.inputs." + i.Id
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
		ctx.Metadata.Templates = append(ctx.Metadata.Templates, template)
	}
	return nil
}
