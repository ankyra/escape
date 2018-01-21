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
	"github.com/ankyra/escape/controllers"
	"github.com/spf13/cobra"
)

var project, application, appVersion string

var inventoryCmd = &cobra.Command{
	Use:     "inventory",
	Short:   "Interact with an Escape Inventory",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var inventoryQueryCommand = &cobra.Command{
	Use:     "query",
	Short:   "Query projects, applications and releases",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		result := controllers.InventoryController{}.Query(context, project, application, appVersion)

		return result.Print(jsonFlag)
	},
}

func init() {
	RootCmd.AddCommand(inventoryCmd)
	inventoryCmd.AddCommand(inventoryQueryCommand)

	inventoryQueryCommand.Flags().StringVarP(&project, "project", "p", "", "The project")
	inventoryQueryCommand.Flags().StringVarP(&application, "application", "a", "", "The application")
	inventoryQueryCommand.Flags().StringVarP(&appVersion, "version", "v", "", "The application version")
	inventoryQueryCommand.PersistentFlags().BoolVarP(&jsonFlag, "json", "", false, "Output profile in JSON format")
}
