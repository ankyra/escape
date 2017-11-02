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

var project, application, appVersion string

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Interact with an Escape Registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Unknown command '%s'", args[0])
		}
		cmd.UsageFunc()(cmd)
		return nil
	},
}

var registryQueryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query projects, applications and releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return controllers.RegistryController{}.Query(context, project, application, appVersion)
	},
}

func init() {
	RootCmd.AddCommand(registryCmd)
	registryCmd.AddCommand(registryQueryCommand)

	registryQueryCommand.Flags().StringVarP(&project, "project", "p", "", "The project")
	registryQueryCommand.Flags().StringVarP(&application, "application", "a", "", "The application")
	registryQueryCommand.Flags().StringVarP(&appVersion, "version", "v", "", "The application version")
}
