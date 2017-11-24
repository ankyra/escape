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

	corestate "github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/controllers"
	"github.com/spf13/cobra"
)

var readLocalErrands bool

var errandsCmd = &cobra.Command{
	Use:     "errands",
	Short:   "List and run errands",
	PreRunE: NoExtraArgsPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.UsageFunc()(cmd)
		return nil
	},
}

func ListLocalErrands(state, environment, escapePlanLocation string) error {
	if err := ProcessFlagsForContext(true); err != nil {
		return err
	}
	return controllers.ErrandsController{}.List(context).Print(jsonFlag)
}

func ListDeployedErrands(state, environment, deployment string) error {
	if err := ProcessFlagsForContext(false); err != nil {
		return err
	}
	deplState, exists := context.GetEnvironmentState().Deployments[deployment]
	if !exists {
		return fmt.Errorf("The deployment '%s' could not be found in environment '%s'", deployment, environment)
	}
	releaseId := deplState.GetReleaseId("deploy")
	if err := context.InitReleaseMetadataByReleaseId(releaseId); err != nil {
		return err
	}
	return controllers.ErrandsController{}.List(context).Print(jsonFlag)
}

func ListErrands(cmd *cobra.Command, args []string) error {
	if readLocalErrands {
		return ListLocalErrands(state, environment, escapePlanLocation)
	}
	if len(deployment) == 0 {
		cmd.UsageFunc()(cmd)
		return nil
	}
	return ListDeployedErrands(state, environment, deployment)
}

var errandsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List errands",
	PreRunE: NoExtraArgsPreRunE,
	RunE:    ListErrands,
}

var errand string

var errandsRunCmd = &cobra.Command{
	Use:   "run <errand>",
	Short: "Run an errand",
	RunE: func(cmd *cobra.Command, args []string) error {
		if environment == "" {
			return fmt.Errorf("Missing 'environment'")
		}
		if len(args) != 1 {
			cmd.UsageFunc()(cmd)
			return nil
		}
		if deployment == "" {
			return fmt.Errorf("Missing deployment name")
		}
		context.SetRootDeploymentName(deployment)
		if err := context.LoadLocalState(state, environment); err != nil {
			return err
		}
		parsedExtraVars, err := ParseExtraVars(extraVars)
		if err != nil {
			return err
		}
		errand := args[0]

		if readLocalErrands {
			return RunLocalErrand(state, environment, escapePlanLocation, errand, parsedExtraVars)
		}
		return RunDeployedErrand(deployment, errand, parsedExtraVars)
	},
}

func RunDeployedErrand(deployment, errand string, parsedExtraVars map[string]string) error {
	deplState := context.GetEnvironmentState().Deployments[deployment]
	if deplState == nil {
		return fmt.Errorf("The deployment '%s' could not be found in environment '%s'.", deployment, context.GetEnvironmentState().Name)
	}

	if deplState.GetStageOrCreateNew("deploy").Status.Code != corestate.OK {
		return fmt.Errorf("'%s' has not been deployed in the environment '%s'.", deployment, context.GetEnvironmentState().Name)
	}
	releaseId := deplState.GetReleaseId("deploy")
	if err := context.InitReleaseMetadataByReleaseId(releaseId); err != nil {
		return err
	}
	return controllers.ErrandsController{}.RunRemoteErrand(context, errand, parsedExtraVars)
}

func RunLocalErrand(state, environment, escapePlanLocation, errand string, parsedExtraVars map[string]string) error {
	if err := ProcessFlagsForContext(true); err != nil {
		return err
	}

	deplState := context.GetEnvironmentState().Deployments[deployment]
	if deplState == nil {
		return fmt.Errorf("The deployment '%s' could not be found in environment '%s'.", deployment, context.GetEnvironmentState().Name)
	}

	if deplState.GetStageOrCreateNew("deploy").Status.Code != corestate.OK {
		return fmt.Errorf("'%s' has not been deployed in the environment '%s'. Use 'escape run deploy' to deploy it.", deployment, context.GetEnvironmentState().Name)
	}

	return controllers.ErrandsController{}.Run(context, errand, parsedExtraVars)
}

func init() {
	RootCmd.AddCommand(errandsCmd)
	errandsCmd.AddCommand(errandsListCmd)
	errandsCmd.AddCommand(errandsRunCmd)
	setPlanAndStateFlags(errandsListCmd)
	errandsListCmd.Flags().BoolVarP(&readLocalErrands, "local", "", false, "Read errands from Escape plan instead of deployment")
	errandsListCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "", false, "Output profile in JSON format")

	setPlanAndStateFlags(errandsRunCmd)
	errandsRunCmd.Flags().StringArrayVarP(&extraVars, "extra-vars", "v", []string{}, "Extra variables (format: key=value, key=@value.txt, @values.json)")
	errandsRunCmd.Flags().BoolVarP(&readLocalErrands, "local", "", false, "Read errands from Escape plan instead of deployment")
}
