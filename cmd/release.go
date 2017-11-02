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
	"github.com/ankyra/escape/controllers"
	"github.com/spf13/cobra"
)

var skipTests, skipCache, skipPush, skipDestroyBuild, skipDeploy, skipSmoke, skipDestroyDeploy, skipDestroy bool

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release (build, test, package, push)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}
		parsedExtraProviders, err := ParseExtraVars(extraProviders)
		if err != nil {
			return err
		}
		return controllers.ReleaseController{}.Release(context, uber, skipBuild, skipTests,
			skipCache, skipPush, skipDestroyBuild, skipDeploy, skipSmoke, skipDestroyDeploy, skipDestroy, force, parsedExtraVars, parsedExtraProviders)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
	setPlanAndStateFlags(releaseCmd)

	releaseCmd.Flags().BoolVarP(&uber, "uber", "u", false, "Build an uber package containing all dependencies")
	releaseCmd.Flags().BoolVarP(&skipBuild, "skip-build", "", false, "Skip build")
	releaseCmd.Flags().BoolVarP(&skipTests, "skip-tests", "", false, "Skip tests")
	releaseCmd.Flags().BoolVarP(&skipCache, "skip-cache", "", false, "Skip caching the release")
	releaseCmd.Flags().BoolVarP(&skipPush, "skip-push", "", false, "Skip push")
	releaseCmd.Flags().BoolVarP(&skipDeploy, "skip-deploy", "", false, "Skip deploy")
	releaseCmd.Flags().BoolVarP(&skipSmoke, "skip-smoke", "", false, "Skip smoke tests")
	releaseCmd.Flags().BoolVarP(&skipDestroy, "skip-destroy", "", false, "Skip destroy steps")
	releaseCmd.Flags().BoolVarP(&skipDestroyBuild, "skip-build-destroy", "", false, "Skip build destroy step")
	releaseCmd.Flags().BoolVarP(&skipDestroyDeploy, "skip-deploy-destroy", "", false, "Skip deploy destroy step")
	releaseCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
	releaseCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	releaseCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")
}
