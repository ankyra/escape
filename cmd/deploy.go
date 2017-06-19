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
		context.SetRootDeploymentName(deployment)
		ctrl := controllers.DeployController{}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}
		if len(args) == 0 {
			if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
				return err
			}
			return ctrl.Deploy(context, parsedExtraVars)
		} else {
			if err := context.LoadLocalState(state, environment); err != nil {
				return err
			}
			for _, arg := range args {
				if err := ctrl.FetchAndDeploy(context, arg, parsedExtraVars); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	setLocalPlanAndStateFlags(deployCmd)
	deployCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
}
