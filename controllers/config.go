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
	"fmt"

	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape/util"
)

type ConfigController struct{}

func (ConfigController) ShowProfile(context Context, json bool) *ControllerResult {
	result := &ControllerResult{
		HumanOutput:       fmt.Sprintf("Profile: %s\n", context.GetEscapeConfig().ActiveProfile),
		MarshalableOutput: context.GetEscapeConfig().GetCurrentProfile(),
	}

	configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
	for k, v := range configMap {
		result.HumanOutput = fmt.Sprintf("%s\n%s: %v", result.HumanOutput, k, v)
	}

	return result
}

func (ConfigController) ShowProfileField(context Context, field string) *ControllerResult {
	configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
	if configMap[field] == nil {
		return &ControllerResult{
			Error: fmt.Errorf(`"%s" is not a valid field name`, field),
		}
	}

	return &ControllerResult{
		HumanOutput:       fmt.Sprintf("%s: %v\n", field, configMap[field]),
		MarshalableOutput: configMap[field],
	}
}

func (ConfigController) ActiveProfile(context Context) *ControllerResult {
	return &ControllerResult{
		HumanOutput:       context.GetEscapeConfig().ActiveProfile,
		MarshalableOutput: context.GetEscapeConfig().ActiveProfile,
	}
}

func (ConfigController) ListProfiles(context Context) *ControllerResult {
	result := &ControllerResult{
		MarshalableOutput: []string{},
	}

	i := 0
	for profileName, _ := range context.GetEscapeConfig().Profiles {
		if i == 0 {
			result.HumanOutput = profileName
		} else {
			result.HumanOutput = fmt.Sprintf("%s\n%s", result.HumanOutput, profileName)
		}
		result.MarshalableOutput = append(result.MarshalableOutput.([]string), profileName)
		i++
	}

	return result
}

func (ConfigController) SetProfile(context Context, profile string) *ControllerResult {
	err := context.GetEscapeConfig().SetActiveProfile(profile)
	if err != nil {
		return &ControllerResult{
			Error: err,
		}
	}
	return &ControllerResult{
		HumanOutput:       "Profile has been set",
		MarshalableOutput: "Profile has been set",
		Error:             context.GetEscapeConfig().Save(),
	}
}
