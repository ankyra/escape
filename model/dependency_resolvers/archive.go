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

package dependency_resolvers

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"io"
	"os"
	"path/filepath"
)

func FromLocalArchive(path *paths.Path, dep *core.Dependency) (bool, error) {
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
	return FromLocalReleaseJson(path, dep)
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
