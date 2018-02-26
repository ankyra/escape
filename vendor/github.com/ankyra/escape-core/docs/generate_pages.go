package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type Page struct {
	Name       string
	Slug       string
	SrcFile    string
	StructName string
}

var Pages = map[string]Page{
	"consumer":  Page{"Providers and Consumers", "providers-and-consumers", "consumer.go", "ConsumerConfig"},
	"depends":   Page{"Dependencies", "dependencies", "dependency_config.go", "DependencyConfig"},
	"downloads": Page{"Downloads", "downloads", "download_config.go", "DownloadConfig"},
	"errands":   Page{"Errands", "errands", "errand.go", "Errand"},
	"extends":   Page{"Extensions", "extensions", "extension_config.go", "ExtensionConfig"},
	"templates": Page{"Templates", "templates", "templates/templates.go", "Template"},
	"variables": Page{"Input and Output Variables", "input-and-output-variables", "variables/variable.go", "Variable"},
}

const PageHeader = `---
date: 2017-11-11 00:00:00
title: "%s"
slug: %s
type: "docs"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/%s
---

%s

Field | Type | Description
------|------|-------------
%s
`

func GetJsonFieldFromTag(tag string) string {
	for _, s := range strings.Split(tag, " ") {
		s = strings.Trim(s, "`")
		if strings.HasPrefix(s, "json:\"") {
			s = s[6 : len(s)-1]
			parts := strings.Split(s, ",")
			return parts[0]
		}
	}
	return ""
}

func ParseType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return ParseType(t.X) + "." + t.Sel.String() // probably wrong
	case *ast.ArrayType:
		return "[" + ParseType(t.Elt) + "]"
	case *ast.StarExpr:
		return ParseType(t.X)
	case *ast.MapType:
		return "{" + ParseType(t.Key) + ":" + ParseType(t.Value) + "}"
	case *ast.InterfaceType:
		return "any"
	default:
		fmt.Printf("%T\n", t)
		panic("type not supported in documentation: ")
	}
	return ""
}

func StructTable(page Page, topLevelDoc string, s *ast.TypeSpec) string {
	structType := s.Type.(*ast.StructType)
	result := ""
	for _, field := range structType.Fields.List {
		tag := GetJsonFieldFromTag(field.Tag.Value)
		typ := ParseType(field.Type)
		result += "|" + tag + "|`" + typ + "`|"
		doc := strings.TrimSpace(field.Doc.Text())
		if doc != "" {
			for _, line := range strings.Split(doc, "\n") {
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
		}
		result += "\n"
	}
	return fmt.Sprintf(PageHeader, page.Name, page.Slug, page.SrcFile, topLevelDoc, result)
}

func GenerateStructDocs(f *ast.File, page Page) string {
	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if s, ok := spec.(*ast.TypeSpec); ok {
					switch s.Type.(type) {
					case *ast.StructType:
						if s.Name.String() == page.StructName {
							return StructTable(page, gen.Doc.Text(), s)
						}
					}
				}
			}
		}
	}
	return ""
}

func main() {
	os.Mkdir("docs/generated/", 0755)
	for _, page := range Pages {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, page.SrcFile, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		str := GenerateStructDocs(f, page)
		filename := "docs/generated/" + page.Slug + ".md"
		fmt.Println("Writing ", filename)
		ioutil.WriteFile(filename, []byte(str), 0644)
	}
}
