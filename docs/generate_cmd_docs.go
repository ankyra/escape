package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/ankyra/escape/cmd"
	"github.com/ankyra/escape/model/escape_plan"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: 2017-11-11 00:00:00
title: "%s"
slug: %s
type: "docs"
toc: true
---
`

const planHeader = `---
date: 2017-11-11 00:00:00
title: "The Escape Plan"
slug: escape-plan
type: "docs"
toc: true
---

The Escape Plan describes a package. 

Field | Type | Description
------|------|-------------
`

var typeMap = map[string][]string{
	"extends":          []string{"[string]", "Extensions"},
	"depends":          []string{"[string]", "Dependencies"},
	"consumes":         []string{"[string]", "Consumers"},
	"build_consumes":   []string{"[string]", "Consumers"},
	"deploy_consumes":  []string{"[string]", "Consumers"},
	"provides":         []string{"[string]", "Consumers"},
	"inputs":           []string{"[string]", "Variables"},
	"build_inputs":     []string{"[string]", "Variables"},
	"deploy_inputs":    []string{"[string]", "Variables"},
	"outputs":          []string{"[string]", "Variables"},
	"metadata":         []string{"{}"},
	"includes":         []string{"[]string"},
	"errands":          []string{"Errands"},
	"downloads":        []string{"Downloads"},
	"templates":        []string{"Templates"},
	"build_templates":  []string{"Templates"},
	"deploy_templates": []string{"Templates"},
}

var typeLinks = map[string]string{
	"Extensions":   "extensions",
	"Dependencies": "dependencies",
	"Consumers":    "providers-and-consumers",
	"Variables":    "input-and-output-variables",
	"Errands":      "errands",
	"Downloads":    "downloads",
	"Templates":    "templates",
}

func GeneratePlanDocs() {
	result := planHeader
	for _, field := range escape_plan.Fields {
		typ := ""
		types, ok := typeMap[field]
		if !ok {
			typ = "`string`"
		}
		for _, t := range types {
			typLink, ok := typeLinks[t]
			if ok {
				typ += "[" + t + "](/docs/" + typLink + "/) "
			} else {
				typ += "`" + t + "` "
			}
		}
		desc := string(escape_plan.GetDoc(field))
		desc = strings.TrimSpace(desc)
		result += "|" + field + "|" + typ + "|"
		for _, line := range strings.Split(desc, "\n") {
			if strings.HasPrefix(line, "#") {
				line = line[1:]
			}
			line = strings.TrimSpace(line)
			if line == "" {
				result += "\n|||"
			} else {
				result += line + " "
			}
		}
		result += "\n"
	}
	if err := ioutil.WriteFile("./docs/escape-plan.md", []byte(result), 0644); err != nil {
		panic(err)
	}
}

func main() {
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "../" + strings.ToLower(base) + "/"
	}
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "./docs/cmd", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
	GeneratePlanDocs()
}
