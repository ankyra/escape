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

package model

import (
	"errors"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/model/config"
	"github.com/ankyra/escape/model/dependency_resolvers"
	"github.com/ankyra/escape/model/paths"
)

type DependencyResolver struct{}
type ReleaseFetcher struct{}

func (resolver DependencyResolver) Resolve(cfg *config.EscapeConfig, resolveDependencies []*core.DependencyConfig) error {
	path := paths.NewPath()
	for _, dep := range resolveDependencies {
		if err := resolver.resolve(cfg, path, dep); err != nil {
			return err
		}
	}
	return nil
}

func (resolver DependencyResolver) resolve(cfg *config.EscapeConfig, path *paths.Path, dep *core.DependencyConfig) error {
	fetcher := ReleaseFetcher{}
	d, err := core.NewDependencyFromString(dep.ReleaseId)
	if err != nil {
		return err
	}
	err = fetcher.Fetch(cfg, path, d)
	if err != nil {
		return err
	}
	releaseJson := path.UnpackedDepDirectoryReleaseMetadata(d)
	metadata, err := core.NewReleaseMetadataFromFile(releaseJson)
	if err != nil {
		return err
	}
	for _, depDep := range metadata.Depends {
		depPath := path.NewPathForDependency(metadata)
		if err := resolver.resolve(cfg, depPath, depDep); err != nil {
			return err
		}
	}
	for _, extension := range metadata.GetExtensions() {
		depPath := path.NewPathForDependency(metadata)
		if err := resolver.resolve(cfg, depPath, core.NewDependencyConfig(extension)); err != nil {
			return err
		}
	}
	return nil
}

func (ReleaseFetcher) Fetch(cfg *config.EscapeConfig, path *paths.Path, dep *core.Dependency) error {
	fetchers := []func(*config.EscapeConfig, *paths.Path, *core.Dependency) (bool, error){
		localFileReleaseFetcherStrategy,
		archiveReleaseFetcherStrategy,
		escapeServerReleaseFetcherStrategy,
	}
	for _, fetcher := range fetchers {
		ok, err := fetcher(cfg, path, dep)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return errors.New("Could not resolve dependency: " + dep.GetReleaseId())
}

func localFileReleaseFetcherStrategy(cfg *config.EscapeConfig, path *paths.Path, dep *core.Dependency) (bool, error) {
	return dependency_resolvers.FromLocalReleaseJson(path, dep)
}

func archiveReleaseFetcherStrategy(cfg *config.EscapeConfig, path *paths.Path, dep *core.Dependency) (bool, error) {
	return dependency_resolvers.FromLocalArchive(path, dep)
}

func escapeServerReleaseFetcherStrategy(cfg *config.EscapeConfig, path *paths.Path, dep *core.Dependency) (bool, error) {
	if err := path.EnsureDependencyCacheDirectoryExists(dep.Project); err != nil {
		return false, err
	}
	targetFile := path.DependencyDownloadTarget(dep)
	inventory := cfg.GetInventory()

	if err := inventory.DownloadRelease(dep.Project, dep.Name, dep.Version, targetFile); err != nil {
		return false, err
	}
	return archiveReleaseFetcherStrategy(cfg, path, dep)
}
