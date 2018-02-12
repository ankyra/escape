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

package cmd

import (
	"fmt"

	"github.com/ankyra/escape/controllers"
	"github.com/spf13/cobra"
)

var refresh bool
var skipDeployment bool
var uber bool

var skipBuild, skipTests bool
var skipCache, skipPush bool
var skipDeploy, skipSmoke bool
var skipDestroyBuild, skipDestroyDeploy, skipDestroy bool
var skipIfExists bool
var toEnv, toDeployment string

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run Escape steps: build, converge, deploy, package, release, smoke, test",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var runBuildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Build the Escape plan using a local state file.",
	PreRunE: NoExtraArgsPreRunE,
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
		return controllers.BuildController{}.Build(context, uber, parsedExtraVars, parsedExtraProviders)
	},
}

var runConvergeCmd = &cobra.Command{
	Use:     "converge",
	Short:   "Bring the environment into its desired state",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(false); err != nil {
			return err
		}
		return controllers.ConvergeController{}.Converge(context, refresh)
	},
}

var runDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a release unit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		loadLocalEscapePlan := len(args) == 0
		if err := ProcessFlagsForContext(loadLocalEscapePlan); err != nil {
			return err
		}

		ctrl := controllers.DeployController{}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}

		parsedExtraProviders, err := ParseExtraVars(extraProviders)
		if err != nil {
			return err
		}
		if loadLocalEscapePlan {
			return ctrl.Deploy(context, parsedExtraVars, parsedExtraProviders)
		} else {
			for _, arg := range args {
				if err := ctrl.FetchAndDeploy(context, arg, parsedExtraVars, parsedExtraProviders); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

var runDestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Destroy the deployment of the current release in the local state file.",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		return controllers.DestroyController{}.Destroy(context, !skipBuild, !skipDeployment)
	},
}

var runPackageCmd = &cobra.Command{
	Use:     "package",
	Short:   "Create a package",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		return controllers.PackageController{}.Package(context, force)
	},
}

var runReleaseCmd = &cobra.Command{
	Use:     "release",
	Short:   "Release (build, test, package, push)",
	PreRunE: NoExtraArgsPreRunE,
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
			skipCache, skipPush, skipDestroyBuild, skipDeploy, skipSmoke, skipDestroyDeploy, skipDestroy, skipIfExists, force, parsedExtraVars, parsedExtraProviders)
	},
}

var runSmokeCmd = &cobra.Command{
	Use:     "smoke",
	Short:   "Run smoke tests using a local state file.",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		return controllers.SmokeController{}.Smoke(context)
	},
}

var runTestCmd = &cobra.Command{
	Use:     "test",
	Short:   "Run tests using a local state file.",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(true); err != nil {
			return err
		}
		return controllers.TestController{}.Test(context)
	},
}

var runPromoteCmd = &cobra.Command{
	Use:     "promote",
	Short:   "Run a promotion of package from one environment to another",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if remoteState != "" {
			return fmt.Errorf("Currently not supported with remote state")
		}

		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}

		parsedExtraProviders, err := ParseExtraVars(extraProviders)
		if err != nil {
			return err
		}

		return controllers.PromoteController{}.Promote(context, state, toEnv, toDeployment, environment, deployment, parsedExtraVars, parsedExtraProviders, useProfileState, force)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.AddCommand(runBuildCmd)
	setPlanAndStateFlags(runBuildCmd)
	runBuildCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	runBuildCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")

	runCmd.AddCommand(runConvergeCmd)
	setPlanAndStateFlags(runConvergeCmd)
	runConvergeCmd.Flags().BoolVarP(&refresh, "refresh", "", false, "Redeploy 'ok' deployments")

	runCmd.AddCommand(runDeployCmd)
	setPlanAndStateFlags(runDeployCmd)
	runDeployCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	runDeployCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")

	runCmd.AddCommand(runDestroyCmd)
	setPlanAndStateFlags(runDestroyCmd)
	runDestroyCmd.Flags().BoolVarP(&skipDeployment, "skip-deployment", "", false, "Don't destroy the deployment.")
	runDestroyCmd.Flags().BoolVarP(&skipBuild, "skip-build", "", false, "Don't destroy the build")

	runCmd.AddCommand(runPackageCmd)
	setPlanAndStateFlags(runPackageCmd)
	runPackageCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
	runPackageCmd.Flags().BoolVarP(&uber, "uber", "u", false, "Build an uber package containing all dependencies")

	runCmd.AddCommand(runReleaseCmd)
	setPlanAndStateFlags(runReleaseCmd)
	runReleaseCmd.Flags().BoolVarP(&uber, "uber", "u", false, "Build an uber package containing all dependencies")
	runReleaseCmd.Flags().BoolVarP(&skipBuild, "skip-build", "", false, "Skip build")
	runReleaseCmd.Flags().BoolVarP(&skipTests, "skip-tests", "", false, "Skip tests")
	runReleaseCmd.Flags().BoolVarP(&skipCache, "skip-cache", "", false, "Skip caching the release")
	runReleaseCmd.Flags().BoolVarP(&skipPush, "skip-push", "", false, "Skip push")
	runReleaseCmd.Flags().BoolVarP(&skipDeploy, "skip-deploy", "", false, "Skip deploy")
	runReleaseCmd.Flags().BoolVarP(&skipSmoke, "skip-smoke", "", false, "Skip smoke tests")
	runReleaseCmd.Flags().BoolVarP(&skipDestroy, "skip-destroy", "", false, "Skip destroy steps")
	runReleaseCmd.Flags().BoolVarP(&skipDestroyBuild, "skip-build-destroy", "", false, "Skip build destroy step")
	runReleaseCmd.Flags().BoolVarP(&skipDestroyDeploy, "skip-deploy-destroy", "", false, "Skip deploy destroy step")
	runReleaseCmd.Flags().BoolVarP(&skipIfExists, "skip-if-exists", "", false, "Skip all the steps if the version that would be released already exists in the Inventory")
	runReleaseCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")
	runReleaseCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	runReleaseCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")

	runCmd.AddCommand(runSmokeCmd)
	setPlanAndStateFlags(runSmokeCmd)

	runCmd.AddCommand(runTestCmd)
	setPlanAndStateFlags(runTestCmd)

	runCmd.AddCommand(runPromoteCmd)
	setPlanAndStateFlags(runPromoteCmd)
	runPromoteCmd.Flags().StringVarP(&toEnv, "to", "", "", "The logical environment to promote to")
	runPromoteCmd.Flags().StringVarP(&toDeployment, "to-deployment", "", "", "The deployment name to promote to (default is the package's \"project/name\")")
	runPromoteCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	runPromoteCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	runPromoteCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")
}
