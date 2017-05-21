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

	"github.com/ankyra/escape-client/controllers"
	"github.com/spf13/cobra"
)

var deployStage bool

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manage the Escape state file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command " + args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var showStateDeploymentsCmd = &cobra.Command{
	Use:   "show-deployments",
	Short: "Show the deployments",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		if deployment != "" {
			return controllers.StateController{}.ShowDeployment(context, deployment)
		}
		return controllers.StateController{}.ShowDeployments(context)
	},
}

var showDeploymentCmd = &cobra.Command{
	Use:   "show-deployment",
	Short: "Show a deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deployment == "" {
			return fmt.Errorf("Missing deployment name")
		}
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		context.SetRootDeploymentName(deployment)
		return controllers.StateController{}.ShowDeployment(context, deployment)
	},
}

var createStateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create state for the given escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		context.SetRootDeploymentName(deployment)
		stage := "build"
		if deployStage {
			stage = "deploy"
		}
		return controllers.StateController{}.CreateState(context, stage)
	},
}

var showStateCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		context.SetRootDeploymentName(deployment)
		return controllers.StateController{}.ShowDeployment(context, context.GetRootDeploymentName())
	},
}

func init() {
	RootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(showStateDeploymentsCmd)
	stateCmd.AddCommand(showDeploymentCmd)
	stateCmd.AddCommand(createStateCmd)
	stateCmd.AddCommand(showStateCmd)

	showStateDeploymentsCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	showStateDeploymentsCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	showStateDeploymentsCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "Deployment name (default \"<release name>\")")

	showDeploymentCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	showDeploymentCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	showDeploymentCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "Deployment name (default \"<release name>\")")

	createStateCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	createStateCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	createStateCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location of the Escape plan.")
	createStateCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "Deployment name (default \"<release name>\")")
	createStateCmd.Flags().BoolVarP(&deployStage, "deploy", "", false, "Use deployment instead of build stage")

	showStateCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	showStateCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	showStateCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location of the Escape plan.")
	showStateCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "Deployment name (default \"<release name>\")")
}
