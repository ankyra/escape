package main

import (
	"github.com/ankyra/escape-client/cmd"
	"github.com/spf13/cobra/doc"
	"log"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs/cmd")
	if err != nil {
		log.Fatal(err)
	}
}
