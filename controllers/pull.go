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
	core "github.com/ankyra/escape-core"
	. "github.com/ankyra/escape/model/interfaces"
)

type PullController struct{}

func (PullController) PullReleases(context Context, packages []string) error {
	for _, pkg := range packages {
		depCfg := core.NewDependencyConfig(pkg)
		if err := depCfg.EnsureConfigIsParsed(); err != nil {
			return err
		}
		if depCfg.NeedsResolving() {
			metadata, err := context.QueryReleaseMetadata(depCfg)
			if err != nil {
				return err
			}
			depCfg = core.NewDependencyConfig(metadata.GetQualifiedReleaseId())
		}
		if _, err := context.GetDependencyMetadata(depCfg); err != nil {
			return err
		}
	}
	return nil
}
