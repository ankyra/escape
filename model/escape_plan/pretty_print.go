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

package escape_plan

import (
	"bytes"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type prettyPrinter struct {
	IncludeEmpty bool
	Spacing      int
}

var Fields = []string{"name", "version", "description", "license", "logo",
	"extends", "depends",
	"consumes", "build_consumes", "deploy_consumes",
	"provides", "inputs", "build_inputs", "deploy_inputs",
	"outputs", "metadata", "includes", "errands", "downloads",
	"templates", "build_templates", "deploy_templates",
	"pre_build", "build", "post_build", "test",
	"pre_deploy", "deploy", "post_deploy", "smoke",
	"pre_destroy", "destroy", "post_destroy"}

var templateMap = map[string]string{
	"name":             keyValTpl,
	"version":          keyValTpl,
	"description":      keyValTpl,
	"license":          keyValTpl,
	"logo":             keyValTpl,
	"path":             keyValTpl,
	"pre_build":        keyValTpl,
	"build":            keyValTpl,
	"post_build":       keyValTpl,
	"pre_deploy":       keyValTpl,
	"deploy":           keyValTpl,
	"post_deploy":      keyValTpl,
	"pre_destroy":      keyValTpl,
	"destroy":          keyValTpl,
	"post_destroy":     keyValTpl,
	"smoke":            keyValTpl,
	"test":             keyValTpl,
	"depends":          listValTpl,
	"extends":          listValTpl,
	"consumes":         listValTpl,
	"deploy_consumes":  listValTpl,
	"build_consumes":   listValTpl,
	"provides":         listValTpl,
	"includes":         listValTpl,
	"inputs":           listValTpl,
	"build_inputs":     listValTpl,
	"deploy_inputs":    listValTpl,
	"outputs":          listValTpl,
	"templates":        listValTpl,
	"build_templates":  listValTpl,
	"deploy_templates": listValTpl,
	"downloads":        listValTpl,
	"metadata":         mapValTpl,
	"errands":          mapValTpl,
}

type printConf func(*prettyPrinter) *prettyPrinter

func includeEmpty(b bool) printConf {
	return func(p *prettyPrinter) *prettyPrinter {
		p.IncludeEmpty = b
		return p
	}
}

func spacing(i int) printConf {
	return func(p *prettyPrinter) *prettyPrinter {
		p.Spacing = i
		return p
	}
}

func NewPrettyPrinter(cfg ...printConf) *prettyPrinter {
	pp := &prettyPrinter{
		IncludeEmpty: true,
		Spacing:      2,
	}
	for _, c := range cfg {
		pp = c(pp)
	}
	return pp
}

func (e *prettyPrinter) Print(plan *EscapePlan) []byte {
	yamlMap := plan.ToDict()
	writer := bytes.NewBuffer([]byte{})
	ordering := Fields
	for _, key := range ordering {
		val, ok := yamlMap[key]
		if !ok {
			if e.IncludeEmpty {
				val = nil
			} else {
				continue
			}
		}
		prettyPrinted := e.prettyPrintValue(key, val)
		if _, err := writer.Write(prettyPrinted); err != nil {
			panic(err)
		}
		for i := 0; i < e.Spacing; i++ {
			writer.Write([]byte("\n"))
		}
	}
	return writer.Bytes()
}

func (e *prettyPrinter) prettyPrintValue(key string, val interface{}) []byte {
	value, err := yaml.Marshal(val)
	if err != nil {
		panic(err)
	}
	if val == nil {
		value = []byte("")
	}
	tpl := template.New("escape-plan")
	tpl.Funcs(map[string]interface{}{
		"indent": indent,
	})
	tpl, err = tpl.Parse(templateMap[key])
	if err != nil {
		panic(err)
	}
	valueMap := map[string]string{
		"key":   key,
		"value": strings.TrimSpace(string(value)),
	}
	doc := []byte("")
	writer := bytes.NewBuffer(doc)
	if err := tpl.Execute(writer, valueMap); err != nil {
		panic(err)
	}
	return writer.Bytes()
}

func indent(s string) string {
	parts := []string{}
	for _, part := range strings.Split(s, "\n") {
		if part != "" {
			parts = append(parts, "  "+part)
		}
	}
	return strings.Join(parts, "\n")
}

const keyValTpl = `{{ .key }}: {{ if .value }}{{ .value }}{{else}}""{{end}}`
const listValTpl = `{{ .key }}:{{ if eq .value "[]" }} []{{else if eq .value ""}} []{{else}}
{{ .value}}{{end}}`
const mapValTpl = `{{ .key }}:{{ if eq .value "{}" }} {}{{else if eq .value ""}} {}{{else}}
{{ indent .value }}{{end}}`
