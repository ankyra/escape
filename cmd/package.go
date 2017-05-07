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
	"github.com/ankyra/escape-client/controllers"
	"github.com/spf13/cobra"
)

var uber bool

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Create a package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		return controllers.PackageController{}.Package(context, force)
	},
}

func init() {
	RootCmd.AddCommand(packageCmd)

	packageCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")
	packageCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
	packageCmd.Flags().BoolVarP(&uber, "uber", "u", false, "Build an uber package containing all dependencies")
}
