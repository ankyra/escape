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

package controllers

import (
	"fmt"

	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/util"
)

type ConfigController struct{}

func (ConfigController) ShowProfile(context *model.Context, json bool) *ControllerResult {
	result := NewControllerResult()

	result.HumanOutput.AddLine("Profile: %s", context.GetEscapeConfig().ActiveProfile)

	configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
	result.HumanOutput.AddMap(configMap)

	result.MarshalableOutput = context.GetEscapeConfig().GetCurrentProfile()

	return result
}

func (ConfigController) ShowProfileField(context *model.Context, field string) *ControllerResult {
	result := NewControllerResult()

	configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
	if configMap[field] == nil {
		result.Error = fmt.Errorf(`'%s' is not a valid field name`, field)
		return result
	}

	result.HumanOutput.AddLine("%s: %v", field, configMap[field])
	result.MarshalableOutput = configMap[field]

	return result
}

func (ConfigController) ActiveProfile(context *model.Context) *ControllerResult {
	result := NewControllerResult()

	result.HumanOutput.AddLine(context.GetEscapeConfig().ActiveProfile)
	result.MarshalableOutput = context.GetEscapeConfig().ActiveProfile

	return result
}

func (ConfigController) ListProfiles(context *model.Context) *ControllerResult {
	result := NewControllerResult()

	profileNames := []interface{}{}
	for profileName, _ := range context.GetEscapeConfig().Profiles {
		profileNames = append(profileNames, profileName)
	}

	result.HumanOutput.AddList(profileNames)
	result.MarshalableOutput = profileNames

	return result
}

func (ConfigController) SetProfile(context *model.Context, profile string) *ControllerResult {
	result := NewControllerResult()

	err := context.GetEscapeConfig().SetActiveProfile(profile)
	if err != nil {
		return &ControllerResult{
			Error: err,
		}
	}
	result.HumanOutput.AddLine("Profile has been set")
	result.MarshalableOutput = "Profile has been set"
	result.Error = context.GetEscapeConfig().Save()

	return result
}

func (ConfigController) CreateProfile(context *model.Context, targetName string) *ControllerResult {
	result := NewControllerResult()

	if targetName == "" {
		result.Error = fmt.Errorf("Missing profile name")
		return result
	}

	result.Error = context.GetEscapeConfig().NewProfile(targetName)
	if result.Error != nil {
		return result
	}

	result.HumanOutput.AddLine("Profile '%s' has been created.", targetName)
	result.MarshalableOutput = "Profile '" + targetName + "' has been created."
	result.Error = context.GetEscapeConfig().Save()

	return result
}
