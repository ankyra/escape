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

var showStateCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the Escape state file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		return controllers.StateController{}.Show(context)
	},
}

var showStateDeploymentsCmd = &cobra.Command{
	Use:   "show-deployments",
	Short: "Show the deployments",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		return controllers.StateController{}.ShowDeployments(context)
	},
}

var showStateDeploymentCmd = &cobra.Command{
	Use:   "show-deployment",
	Short: "Show a deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Missing deployment name")
		}
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		return controllers.StateController{}.ShowDeployment(context, args[0])
	},
}

var createDeploymentCmd = &cobra.Command{
	Use:   "create-deployment",
	Short: "Create a deployment for a given escape plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		return controllers.StateController{}.CreateDeployment(context)
	},
}

func init() {
	RootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(showStateCmd)
	stateCmd.AddCommand(showStateDeploymentsCmd)
	stateCmd.AddCommand(showStateDeploymentCmd)
	stateCmd.AddCommand(createDeploymentCmd)

	showStateCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")

	showStateDeploymentsCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	showStateDeploymentsCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")

	showStateDeploymentCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	showStateDeploymentCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")

	createDeploymentCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	createDeploymentCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	createDeploymentCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")
}
