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

var errandsCmd = &cobra.Command{
	Use:   "errands",
	Short: "List and run errands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command '%s'", args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var errandsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List errands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
			return err
		}
		return controllers.ErrandsController{}.List(context)
	},
}

var errand string

var errandsRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an errand",
	RunE: func(cmd *cobra.Command, args []string) error {
		if environment == "" {
			return fmt.Errorf("Missing 'environment'")
		}
		if len(args) != 1 {
			return fmt.Errorf("Expecting one errand")
		}
		errand := args[0]
		err := context.LoadLocalState(state, environment)
		if err != nil {
			return err
		}
		err = context.LoadEscapePlan(escapePlanLocation)
		if err != nil {
			return err
		}
		err = context.LoadMetadata()
		if err != nil {
			return err
		}
		return controllers.ErrandsController{}.Run(context, errand)
	},
}

func init() {
	RootCmd.AddCommand(errandsCmd)
	errandsCmd.AddCommand(errandsListCmd)
	errandsCmd.AddCommand(errandsRunCmd)

	errandsListCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	errandsListCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	errandsListCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")

	errandsRunCmd.Flags().StringVarP(&state, "state", "s", "escape_state.json", "Location of the Escape state file")
	errandsRunCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "The logical environment to target")
	errandsRunCmd.Flags().StringVarP(&escapePlanLocation, "input", "i", "escape.yml", "The location onf the Escape plan.")
}
