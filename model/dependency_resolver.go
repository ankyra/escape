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
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DependencyResolver struct{}
type ReleaseFetcher struct{}

func (resolver DependencyResolver) Resolve(cfg EscapeConfig, resolveDependencies []string) error {
	path := paths.NewPath()
	for _, dep := range resolveDependencies {
		if err := resolver.resolve(cfg, path, dep); err != nil {
			return err
		}
	}
	return nil
}

func (resolver DependencyResolver) resolve(cfg EscapeConfig, path Paths, dep string) error {
	fetcher := ReleaseFetcher{}
	d, err := core.NewDependencyFromString(dep)
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
	for _, depDep := range metadata.GetDependencies() {
		depPath := path.NewPathForDependency(metadata)
		if err := resolver.resolve(cfg, depPath, depDep); err != nil {
			return err
		}
	}
	for _, extension := range metadata.GetExtends() {
		depPath := path.NewPathForDependency(metadata)
		if err := resolver.resolve(cfg, depPath, extension); err != nil {
			return err
		}
	}
	return nil
}

func (ReleaseFetcher) Fetch(cfg EscapeConfig, path Paths, dep *core.Dependency) error {
	fetchers := []func(EscapeConfig, Paths, *core.Dependency) (bool, error){
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

func localFileReleaseFetcherStrategy(cfg EscapeConfig, path Paths, dep *core.Dependency) (bool, error) {
	releaseJson := path.UnpackedDepDirectoryReleaseMetadata(dep)
	if util.PathExists(releaseJson) {
		contents, err := ioutil.ReadFile(releaseJson)
		if err != nil {
			return false, errors.New("Couldn't read dependency release.json: " + err.Error())
		}
		escapePlan := map[string]interface{}{}
		err = json.Unmarshal(contents, &escapePlan)
		if err != nil {
			return false, errors.New("Couldn't unmarshal dependency release.json: " + err.Error())
		}
		version, ok := escapePlan["version"]
		if !ok {
			util.RemoveTree(path.UnpackedDepDirectory(dep))
			return false, nil
		}
		if version.(string) != dep.GetVersion() {
			util.RemoveTree(path.UnpackedDepDirectory(dep))
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func archiveReleaseFetcherStrategy(cfg EscapeConfig, path Paths, dep *core.Dependency) (bool, error) {
	localArchive := path.DependencyReleaseArchive(dep)
	buildDirArchive := path.DependencyDownloadTarget(dep)
	if !util.PathExists(localArchive) && !util.PathExists(buildDirArchive) {
		return false, nil
	}
	if util.PathExists(buildDirArchive) {
		localArchive = buildDirArchive
	}
	fp, err := os.Open(localArchive)
	if err != nil {
		return false, errors.New("Couldn't open archive " + localArchive + ": " + err.Error())
	}
	defer fp.Close()

	gzf, err := gzip.NewReader(fp)
	if err != nil {
		return false, errors.New("Couldn't read gzip archive " + localArchive + ": " + err.Error())
	}

	tarReader := tar.NewReader(gzf)
	targetDir := path.DepTypeDirectory(dep)
	finalDir := path.UnpackedDepDirectory(dep)
	path.EnsureDependencyTypeDirectoryExists(dep)
	if util.PathExists(finalDir) {
		err := util.RemoveTree(finalDir)
		if err != nil {
			return false, errors.New("Failed to remove tree " + targetDir + ": " + err.Error())
		}
	}
	err = UnpackTarReader(tarReader, targetDir)
	if err != nil {
		return false, errors.New("Failed to unpack " + localArchive + ": " + err.Error())
	}
	unpackedDir := filepath.Join(targetDir, dep.GetReleaseId())
	if !util.PathExists(unpackedDir) {
		return false, errors.New("Expected path " + unpackedDir + " does not exist")
	}
	err = os.Rename(unpackedDir, finalDir)
	if err != nil {
		return false, err
	}
	return localFileReleaseFetcherStrategy(cfg, path, dep)
}

func escapeServerReleaseFetcherStrategy(cfg EscapeConfig, path Paths, dep *core.Dependency) (bool, error) {
	backend := cfg.GetCurrentTarget().GetStorageBackend()
	if backend != "escape" {
		return false, nil
	}
	if err := path.EnsureDependencyCacheDirectoryExists(); err != nil {
		return false, err
	}
	targetFile := path.DependencyDownloadTarget(dep)
	client := cfg.GetClient()
	if err := client.DownloadRelease(dep.GetReleaseId(), targetFile); err != nil {
		return false, err
	}
	return archiveReleaseFetcherStrategy(cfg, path, dep)
}

func UnpackTarReader(tarReader *tar.Reader, targetDir string) error {
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			util.MkdirRecursively(name)
		case tar.TypeReg, tar.TypeRegA:
			dir, _ := filepath.Split(name)
			if err := util.MkdirRecursively(dir); err != nil {
				return errors.New("Failed to make directory " + dir + ": " + err.Error())
			}
			out, err := os.Create(name)
			if err != nil {
				return errors.New("Couldn't create file: " + name + ": " + err.Error())
			}
			if header.Size != 0 {
				_, err = io.Copy(out, tarReader)
				if err != nil {
					return errors.New("Couldn't write to " + name + ": " + err.Error())
				}
			}
			out.Close()
			os.Chmod(name, os.FileMode(header.Mode))
		default:
			return errors.New("Unsupported type for tar header")
		}
	}
	return nil
}
