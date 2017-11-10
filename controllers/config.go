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

func (ConfigController) ShowProfile(context Context, json bool) {
	if json {
		fmt.Println(context.GetEscapeConfig().GetCurrentProfile().ToJson())
	} else {
		fmt.Printf("Profile: %s\n\n", context.GetEscapeConfig().ActiveProfile)

		configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
		for k, v := range configMap {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}

func (ConfigController) ShowProfileField(context Context, field string) error {
	configMap := util.StructToMapStringInterface(*context.GetEscapeConfig().GetCurrentProfile(), "json")
	if configMap[field] != nil {
		fmt.Printf("%s: %v\n", field, configMap[field])
		return nil
	}

	return fmt.Errorf(`"%s" is not a valid field name`, field)
}

func (ConfigController) ActiveProfile(context Context) {
	fmt.Println(context.GetEscapeConfig().ActiveProfile)
}

func (ConfigController) ListProfiles(context Context) {
	for profileName, _ := range context.GetEscapeConfig().Profiles {
		fmt.Println(profileName)
	}
}

func (ConfigController) SetProfile(context Context) {
	context.GetEscapeConfig().Save()
}
