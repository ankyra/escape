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

package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/state"
	. "github.com/ankyra/escape/model/interfaces"
)

type PromoteController struct{}

func (PromoteController) Promote(context Context, state, toEnv, toDeployment, fromEnv, fromDeployment string, extraVars, extraProviders map[string]string, useProfileState, force bool) error {
	if fromDeployment == "" {
		return fmt.Errorf("Missing deployment name.")
	}
	if toEnv == "" {
		return fmt.Errorf("Missing target environment. Use '--to' to define your target environment.")
	}

	if err := context.LoadLocalState(state, fromEnv, useProfileState); err != nil {
		return err
	}
	context.SetRootDeploymentName(fromDeployment)

	releaseId, err := buildReleaseId(context.GetEnvironmentState(), context.GetRootDeploymentName())
	if err != nil {
		return err
	}

	logKey := "promote.state_info"

	context.PushLogSection("Promote")
	context.Log(logKey, map[string]string{
		"deployment":  context.GetRootDeploymentName(),
		"environment": context.GetEnvironmentState().Name,
		"releaseId":   releaseId,
	})

	if toDeployment == "" {
		toDeployment = fromDeployment
	}

	if err := context.LoadLocalState(state, toEnv, useProfileState); err != nil {
		return err
	}
	context.SetRootDeploymentName(toDeployment)

	toReleaseId, err := buildReleaseId(context.GetEnvironmentState(), context.GetRootDeploymentName())
	if err != nil || toReleaseId == "" {
		logKey = "promote.state_info_missing"
	}

	context.Log(logKey, map[string]string{
		"deployment":  context.GetRootDeploymentName(),
		"environment": context.GetEnvironmentState().Name,
		"releaseId":   toReleaseId,
	})

	if !force {
		response, err := confirmationUserInput(getUserInputReader(), fmt.Sprintf("Promote %s from %s (%s) to %s (%s)? [Yn]", releaseId, fromEnv, fromDeployment, toEnv, toDeployment))
		if err != nil {
			return err
		}

		if !response {
			return nil
		}
	}

	context.Log("promote.promoting", map[string]string{
		"releaseId":       releaseId,
		"fromEnvironment": fromEnv,
		"toEnvironment":   context.GetEnvironmentState().Name,
	})

	return DeployController{}.FetchAndDeploy(context, releaseId, extraVars, extraProviders)
}

func getUserInputReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func confirmationUserInput(reader *bufio.Reader, message string) (bool, error) {
	var err error
	fmt.Printf("%s: ", message)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	switch strings.TrimSpace(input) {
	case "n":
		return false, nil
	case "no":
		return false, nil
	default:
		return true, nil
	}
}

func buildReleaseId(env *state.EnvironmentState, deploymentName string) (string, error) {
	deployment := env.Deployments[deploymentName]
	if deployment == nil {
		return "", fmt.Errorf("Deployment %s was not found in the environment %s.", deploymentName, env.Name)
	}

	if deployStage := deployment.Stages["deploy"]; deployStage == nil {
		return "", fmt.Errorf("Deployment %s has not been deployed in the environment %s.", deploymentName, env.Name)
	}

	version := deployment.Stages["deploy"].Version
	if version == "" {
		return "", fmt.Errorf("Deployment %s has not been deployed in the environment %s.", deploymentName, env.Name)
	}

	releaseId := parsers.ReleaseId{
		Name:    deployment.Release,
		Version: version,
	}

	return releaseId.ToString(), nil
}
