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

var jsonFlag bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the escape client configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command '%s'", args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var configProfileCmd = &cobra.Command{
	Use:   "profile <profile field name>",
	Short: "Show the currently active Escape profile",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var result *controllers.ControllerResult
		if len(args) < 1 {
			result = controllers.ConfigController{}.ShowProfile(context, jsonFlag)
		} else {
			result = controllers.ConfigController{}.ShowProfileField(context, args[0])
		}

		return result.Print(jsonFlag)
	},
}

var configActiveProfileCmd = &cobra.Command{
	Use:   "active-profile",
	Short: "Show the currently active profile name",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		result := controllers.ConfigController{}.ActiveProfile(context)

		result.Print(jsonFlag)
	},
}

var configListProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List the currently available Escape profiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		result := controllers.ConfigController{}.ListProfiles(context)

		result.Print(jsonFlag)
	},
}

var configSetProfileCmd = &cobra.Command{
	Use:   "set-profile <profile name>",
	Short: "Set the active Escape profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result *controllers.ControllerResult
		if len(args) < 1 {
			result = controllers.ConfigController{}.SetProfile(context, cfgProfile)
		} else {
			result = controllers.ConfigController{}.SetProfile(context, args[0])
		}

		return result.Print(jsonFlag)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configProfileCmd)
	configCmd.AddCommand(configListProfilesCmd)
	configCmd.AddCommand(configSetProfileCmd)
	configCmd.AddCommand(configActiveProfileCmd)

	configCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "", false, "Output profile in JSON format")
}
