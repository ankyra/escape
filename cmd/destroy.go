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

var skipDeployment, skipBuild bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the deployment of the current release in the local state file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		return controllers.DestroyController{}.Destroy(context, !skipBuild, !skipDeployment)
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	destroyCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	destroyCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")
	destroyCmd.Flags().BoolVarP(&skipDeployment, "skip-deployment", "", false, "Don't destroy the deployment.")
	destroyCmd.Flags().BoolVarP(&skipBuild, "skip-build", "", false, "Don't destroy the build")
}
