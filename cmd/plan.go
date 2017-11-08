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

package cmd

import (
	"fmt"

	"github.com/ankyra/escape/controllers"
	"github.com/spf13/cobra"
)

var releaseName string
var outputPath string
var force, includeDocs bool

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Manage the Escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command '%s'", args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if releaseName == "" {
			return fmt.Errorf("Missing 'name' parameter")
		}
		return controllers.PlanController{}.Init(context, releaseName, outputPath, force, includeDocs)
	},
}

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format an existing Escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := context.LoadEscapePlan(escapePlanLocation)
		if err != nil {
			return err
		}
		return controllers.PlanController{}.Format(context, escapePlanLocation)
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff compiled Escape plan against latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		return controllers.PlanController{}.Diff(context)
	},
}

var minifyCmd = &cobra.Command{
	Use:   "minify",
	Short: "Minify an existing Escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := context.LoadEscapePlan(escapePlanLocation)
		if err != nil {
			return err
		}
		return controllers.PlanController{}.Minify(context, escapePlanLocation)
	},
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the Escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		controllers.PlanController{}.Compile(context)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(planCmd)
	planCmd.AddCommand(initCmd)
	planCmd.AddCommand(fmtCmd)
	planCmd.AddCommand(minifyCmd)
	planCmd.AddCommand(previewCmd)
	planCmd.AddCommand(diffCmd)

	initCmd.Flags().StringVarP(&releaseName, "name", "n", "", "The release name (eg. hello-world)")
	initCmd.Flags().StringVarP(&outputPath, "output", "o", "escape.yml", "The output location")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
	initCmd.Flags().BoolVarP(&includeDocs, "include-docs", "", false, "Include all fields and their relevant documentation")

	setPlanAndStateFlags(previewCmd)
	setPlanAndStateFlags(diffCmd)
	setEscapePlanLocationFlag(fmtCmd)
	setEscapePlanLocationFlag(minifyCmd)
}
