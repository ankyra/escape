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

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape/model"
)

type TagController struct{}

func (t TagController) TagRelease(context *model.Context, releaseID, tag string) error {
	parsed, err := parsers.ParseQualifiedReleaseId(releaseID)
	if err != nil {
		return err
	}
	if !parsers.IsValidTag(tag) {
		return fmt.Errorf("The tag '%s' is not allowed.", tag)
	}
	return context.GetInventory().TagRelease(parsed.Project, parsed.Name, parsed.Version, tag)
}
