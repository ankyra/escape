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

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
)

type ReleaseController struct{}

func (r ReleaseController) Release(context Context, buildFatPackage, skipTests, skipCache, skipPush, forceOverwrite bool) error {
	context.PushLogRelease(context.GetReleaseMetadata().GetReleaseId())
	context.PushLogSection("Release")
	context.Log("release.start", nil)
	if err := (BuildController{}).Build(context, buildFatPackage); err != nil {
		return err
	}
	if !skipTests {
		if err := (TestController{}).Test(context); err != nil {
			return err
		}
	}
	if err := (PackageController{}).Package(context, forceOverwrite); err != nil {
		return err
	}
	if !skipCache {
		if err := r.cacheRelease(context, forceOverwrite); err != nil {
			return err
		}
	}
	if !skipPush {
		if err := (PushController{}).Push(context, buildFatPackage); err != nil {
			return err
		}
	}
	context.Log("release.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}

func (r ReleaseController) cacheRelease(context Context, forceOverwrite bool) error {
	path := paths.NewPath()
	metadata := context.GetReleaseMetadata()
	packagePath := path.ReleaseLocation(metadata)
	if err := path.EnsureDependencyCacheDirectoryExists(); err != nil {
		return err
	}
	userPackageCachePath := path.DependencyDownloadTarget(metadata.ToDependency())
	if util.PathExists(userPackageCachePath) && !forceOverwrite {
		return fmt.Errorf("Release already exists in local release cache: %s. Use --force / -f to overwrite", userPackageCachePath)
	}
	if err := util.CopyFile(packagePath, userPackageCachePath); err != nil {
		return err
	}
	return nil
}
