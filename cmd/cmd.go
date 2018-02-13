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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var state, environment, deployment, escapePlanLocation, remoteState string
var useProfileState bool

func LoadState() error {
	if remoteState != "" {
		if err := context.LoadRemoteState(remoteState, environment); err != nil {
			return err
		}
	}
	return context.LoadLocalState(state, environment, useProfileState)
}

func ProcessFlagsForContext(loadLocalEscapePlan bool) error {
	if environment == "" {
		return fmt.Errorf("Missing 'environment'")
	}
	if err := LoadState(); err != nil {
		return err
	}
	if loadLocalEscapePlan {
		if err := context.LoadEscapePlan(escapePlanLocation); err != nil {
			return err
		}
		if err := context.CompileEscapePlan(); err != nil {
			return err
		}
	}
	context.SetRootDeploymentName(deployment)
	return nil
}

func setEscapePlanLocationFlag(c *cobra.Command) {
	c.Flags().StringVarP(&escapePlanLocation,
		"plan", "", "escape.yml",
		"The location of the Escape plan.",
	)
}

func setEscapeStateLocationFlag(c *cobra.Command) {
	c.Flags().StringVarP(&state,
		"state", "s", "escape_state.json",
		"Location of the Escape state file (ignored when --remote-state is set)",
	)
	c.Flags().BoolVarP(&useProfileState,
		"use-profile-state", "", false,
		"Instead of using the Escape state file specified in --state, read the 'state_path' value from the configuration profile.")
}

func setEscapeStateEnvironmentFlag(c *cobra.Command) {
	c.Flags().StringVarP(&environment,
		"environment", "e", "dev",
		"The logical environment to target",
	)
}

func setEscapeDeploymentFlag(c *cobra.Command) {
	c.Flags().StringVarP(&deployment,
		"deployment", "d", "",
		"Deployment name (default is the package's \"project/name\")",
	)
}

func setEscapeRemoteStateFlag(c *cobra.Command) {
	c.Flags().StringVarP(&remoteState,
		"remote-state", "r", "",
		"Use remote state project.")
}

func setPlanAndStateFlags(c *cobra.Command) {
	setEscapePlanLocationFlag(c)
	setEscapeStateLocationFlag(c)
	setEscapeStateEnvironmentFlag(c)
	setEscapeDeploymentFlag(c)
	setEscapeRemoteStateFlag(c)
}

func ParseExtraVars(extraVars []string) (result map[string]string, err error) {
	result = map[string]string{}
	for _, extraVar := range extraVars {
		err = fmt.Errorf("Invalid extra variable format '%s'", extraVar)
		parts := strings.Split(extraVar, "=")
		if len(parts) == 0 {
			return nil, err
		}
		key := parts[0]
		value := strings.Join(parts[1:], "=")
		if value == "" {
			if strings.HasPrefix(key, "@") {
				v, err := ioutil.ReadFile(key[1:])
				if err != nil {
					return nil, fmt.Errorf("Coulnd't read file '%s': %s", key[1:], err.Error())
				}
				unmarshalled := map[string]interface{}{}
				err = json.Unmarshal(v, &unmarshalled)
				if err != nil {
					return nil, fmt.Errorf("Coulnd't read file '%s' into JSON map: %s", key[1:], err.Error())
				}
				for key, val := range unmarshalled {
					switch val.(type) {
					case string:
						result[key] = val.(string)
					default:
						marshalled, err := json.Marshal(val)
						if err != nil {
							return nil, err
						}
						result[key] = string(marshalled)
					}
				}
			} else {
				result[key] = ""
			}
		} else if strings.HasPrefix(value, "@") {
			v, err := ioutil.ReadFile(value[1:])
			if err != nil {
				return nil, fmt.Errorf("Coulnd't read file '%s': %s", value[1:], err.Error())
			}
			result[key] = string(v)
		} else {
			result[key] = value
		}
	}
	return result, nil
}

func NoExtraArgsPreRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Unknown command '%s'", args[0])
	}
	return nil
}
