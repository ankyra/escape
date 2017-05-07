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

var skipTests, skipCache, skipPush bool

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release (build, test, package, push)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		return controllers.ReleaseController{}.Release(context, uber, skipTests, skipCache, skipPush, force)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	releaseCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	releaseCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")
	releaseCmd.Flags().BoolVarP(&uber, "uber", "u", false, "Build an uber package containing all dependencies")
	releaseCmd.Flags().BoolVarP(&skipTests, "skip-tests", "", false, "Skip tests")
	releaseCmd.Flags().BoolVarP(&skipCache, "skip-cache", "", false, "Skip caching the release")
	releaseCmd.Flags().BoolVarP(&skipPush, "skip-push", "", false, "Skip push")
	releaseCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
}
