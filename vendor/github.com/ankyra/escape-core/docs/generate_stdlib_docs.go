package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ankyra/escape-core/script"
)

type Type struct {
	Methods map[string]string
}

func (t *Type) AddMethod(f script.StdlibFunc) {
	if t.Methods == nil {
		t.Methods = map[string]string{}
	}
	header := fmt.Sprintf("%s(%s)", f.Id, f.Args)
	t.Methods[header] = f.Doc
}

func main() {
	class := map[string]*Type{}

	for _, f := range script.Stdlib {
		cls, found := class[f.ActsOn]
		if !found {
			cls = &Type{}
			class[f.ActsOn] = cls
		}
		cls.AddMethod(f)
	}

	s := `---
title: "Escape Standard Library Reference"
slug: scripting-language-stdlib 
type: "reference"
toc: true
---

<style>
h2 {
  font-size: 0.8em;
  font-family: mono;
  background: #4B9CD3;
  padding: 5px;
}
</style>

Standard library functions for the [Escape Scripting Language](../scripting-language/)

`
	for cls, typ := range class {
		if cls == "" {
			s = fmt.Sprintf("%s\n# Unary functions\n\n", s)
		} else {
			s = fmt.Sprintf("%s\n# Functions acting on %s\n\n", s, cls)
		}
		for sig, doc := range typ.Methods {
			s = fmt.Sprintf("%s## %s\n\n%s\n\n", s, sig, doc)
		}
	}
	os.Mkdir("docs/generated/", 0755)
	ioutil.WriteFile("docs/generated/stdlib.md", []byte(s), 0644)
}
