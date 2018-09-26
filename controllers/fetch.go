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
	"os"

	"github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/model/paths"
)

type FetchController struct{}

func (FetchController) Fetch(context *model.Context, releaseIds []string) error {
	for _, releaseId := range releaseIds {
		err := model.DependencyResolver{}.Resolve(context.GetEscapeConfig(), []*core.DependencyConfig{core.NewDependencyConfig(releaseId)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f FetchController) ResolveFetchAndLoad(context *model.Context, releaseId string) error {
	// TODO cd into temp directory ?
	parsed := core.NewDependencyConfig(releaseId)
	if err := parsed.EnsureConfigIsParsed(); err != nil {
		return err
	}
	if parsed.NeedsResolving() {
		metadata, err := context.QueryReleaseMetadata(parsed)
		if err != nil {
			return err
		}
		parsed.Version = metadata.Version
		metadata.Project = parsed.Project // inventory needs to be updated to latest core
		releaseId = metadata.GetQualifiedReleaseId()
	}
	if err := f.Fetch(context, []string{releaseId}); err != nil {
		return err
	}
	root := paths.NewPath().UnpackedDepCfgDirectory(parsed)
	err := os.Chdir(root)
	if err != nil {
		return err
	}
	return context.LoadReleaseJson()
}
