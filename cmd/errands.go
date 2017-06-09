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
		context.SetRootDeploymentName(deployment)
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
			return fmt.Errorf("Expecting errand")
		}
		context.SetRootDeploymentName(deployment)
		if deployment != "" {
			if err := context.LoadLocalState(state, environment); err != nil {
				return err
			}
			deplState := context.GetEnvironmentState().GetOrCreateDeploymentState(deployment)
			releaseId := deplState.GetReleaseId("deploy")
			// todo create temp dir?
			if err := context.InitReleaseMetadataByReleaseId(releaseId); err != nil {
				return err
			}
			// todo: cd into directory
		} else {
			if err := context.InitFromLocalEscapePlanAndState(state, environment, escapePlanLocation); err != nil {
				return err
			}
		}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}
		errand := args[0]
		return controllers.ErrandsController{}.Run(context, errand, parsedExtraVars)
	},
}

func init() {
	RootCmd.AddCommand(errandsCmd)
	errandsCmd.AddCommand(errandsListCmd)
	errandsCmd.AddCommand(errandsRunCmd)
	setLocalPlanAndStateFlags(errandsListCmd)

	setLocalPlanAndStateFlags(errandsRunCmd)
	errandsRunCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
}
