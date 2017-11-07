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

var refresh bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Escape steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command '%s'", args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var runBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Escape plan using a local state file.",
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
	Use:   "converge",
	Short: "Bring the environment into its desired state",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(false); err != nil {
			return err
		}
		return controllers.ConvergeController{}.Converge(context, refresh)
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

}
