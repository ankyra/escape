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

var deployStage bool
var extraVars, extraProviders []string

var stateCmd = &cobra.Command{
	Use:     "state",
	Short:   "Manage the Escape state file",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var listDeploymentsCmd = &cobra.Command{
	Use:     "list-deployments",
	Short:   "Show the deployments",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		result := controllers.StateController{}.ListDeployments(context)

		return result.Print(jsonFlag)
	},
}

var showDeploymentCmd = &cobra.Command{
	Use:     "show-deployment",
	Short:   "Show a deployment",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if deployment == "" {
			if err := ProcessFlagsForContext(true); err != nil {
				return err
			}
			return controllers.StateController{}.ShowDeployment(context, context.GetRootDeploymentName())
		}
		if err := ProcessFlagsForContext(false); err != nil {
			return err
		}
		return controllers.StateController{}.ShowDeployment(context, deployment)
	},
}

var showProvidersCmd = &cobra.Command{
	Use:     "show-providers",
	Short:   "Show the providers available in the environment",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ProcessFlagsForContext(false); err != nil {
			return err
		}
		result := controllers.StateController{}.ShowProviders(context)

		return result.Print(jsonFlag)
	},
}

var createStateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create state for the given escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		useEscapePlan := len(args) == 0
		if err := ProcessFlagsForContext(useEscapePlan); err != nil {
			return err
		}
		stage := "build"
		if deployStage {
			stage = "deploy"
		}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}
		parsedExtraProviders, err := ParseExtraVars(extraProviders)
		if err != nil {
			return err
		}
		if !useEscapePlan {
			if err := context.InitReleaseMetadataByReleaseId(args[0]); err != nil {
				return err
			}
		}
		return controllers.StateController{}.CreateState(context, stage, parsedExtraVars, parsedExtraProviders)
	},
}

func init() {
	RootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(listDeploymentsCmd)
	stateCmd.AddCommand(showDeploymentCmd)
	stateCmd.AddCommand(showProvidersCmd)
	stateCmd.AddCommand(createStateCmd)

	setEscapeStateLocationFlag(listDeploymentsCmd)
	setEscapeStateEnvironmentFlag(listDeploymentsCmd)
	listDeploymentsCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "", false, "Output profile in JSON format")

	setEscapeStateLocationFlag(showDeploymentCmd)
	setEscapeStateEnvironmentFlag(showDeploymentCmd)
	setEscapeDeploymentFlag(showDeploymentCmd)

	setEscapeStateLocationFlag(showProvidersCmd)
	setEscapeStateEnvironmentFlag(showProvidersCmd)
	showProvidersCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "", false, "Output profile in JSON format")

	setPlanAndStateFlags(createStateCmd)
	createStateCmd.Flags().BoolVarP(&deployStage, "deploy", "", false, "Use deployment instead of build stage")
	createStateCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	createStateCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")
}
