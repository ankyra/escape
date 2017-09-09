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
	"github.com/ankyra/escape-client/controllers"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the last built release using a local state file.",
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

func init() {
	RootCmd.AddCommand(deployCmd)
	setPlanAndStateFlags(deployCmd)
	deployCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	deployCmd.Flags().StringArrayVarP(&extraProviders, "extra-providers", "p", []string{}, "Extra providers (format: provider=deployment, provider=@deployment.txt, @values.json)")
}
