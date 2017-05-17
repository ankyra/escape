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
	"errors"
	"github.com/ankyra/escape-client/model/paths"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Archiver struct{}

func NewReleaseArchiver() *Archiver {
	return &Archiver{}
}

func (a *Archiver) Archive(metadata *core.ReleaseMetadata, forceOverwrite bool) error {
	//    applog("archive.start", release=metadata.get_full_build_id())
	if err := buildReleaseAndTargetDirectories(metadata); err != nil {
		return err
	}
	path := paths.NewPath()
	scratchSpace := path.ScratchSpaceDirectory(metadata)
	releaseJsonPath := path.ScratchSpaceReleaseMetadata(metadata)
	if err := metadata.WriteJsonFile(releaseJsonPath); err != nil {
		return err
	}
	for _, dir := range metadata.GetDirectories() {
		scratchDir := filepath.Join(scratchSpace, dir)
		util.MkdirRecursively(scratchDir)
	}
	if err := copyFiles(scratchSpace, metadata.GetFiles()); err != nil {
		return err
	}
	return a.buildTarArchive(metadata, forceOverwrite)
	//    applog("archive.finished", release=metadata.get_full_build_id(), path=os.path.realpath(tar))
}

func (a *Archiver) buildTarArchive(metadata *core.ReleaseMetadata, forceOverwrite bool) error {
	path := paths.NewPath()
	scratchSpace := path.ScratchSpaceDirectory(metadata)
	packageId := metadata.GetReleaseId()
	packageGzip := packageId + ".tgz"
	target := path.ReleaseLocation(metadata)
	if util.PathExists(target) {
		if !forceOverwrite {
			return errors.New("File '" + target + "' already exists. Use --force / -f to ignore.")
		}
		os.Remove(target)
	}
	packageCwd, err := filepath.Abs(filepath.Join(scratchSpace, ".."))
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(packageCwd)
	if err := buildGzip(packageId, packageGzip); err != nil {
		return err
	}
	os.Chdir(currentDir)
	pkg := filepath.Join(packageCwd, packageGzip)
	return os.Rename(pkg, target)
}

func buildReleaseAndTargetDirectories(metadata *core.ReleaseMetadata) error {
	path := paths.NewPath()
	scratchSpace := path.ScratchSpaceDirectory(metadata)
	if util.PathExists(scratchSpace) {
		if err := util.RemoveTree(scratchSpace); err != nil {
			return err
		}
	}
	if err := path.EnsureEscapeDirectoryExists(); err != nil {
		return err
	}
	if err := path.EnsureScratchSpaceDirectoryExists(metadata); err != nil {
		return err
	}
	if err := path.EnsureReleaseTargetDirectoryExists(); err != nil {
		return err
	}
	return nil
}

func copyFiles(baseDir string, files map[string]string) error {
	for file := range files {
		dest := filepath.Join(baseDir, file)
		if err := util.CopyFile(file, dest); err != nil {
			return err
		}
	}
	return nil
}

func buildGzip(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	dirs := []string{src}
	for len(dirs) != 0 {
		dir := dirs[0]
		dirs = dirs[1:]
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, fileInfo := range files {
			path := filepath.Join(dir, fileInfo.Name())
			if fileInfo.IsDir() {
				dirs = append(dirs, path)
			} else {
				in, err := os.Open(path)
				if err != nil {
					return err
				}
				hdr := &tar.Header{
					Name:    path,
					Size:    fileInfo.Size(),
					Mode:    int64(fileInfo.Mode()),
					ModTime: fileInfo.ModTime(),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if _, err := io.Copy(tw, in); err != nil {
					return err
				}
				if err := in.Close(); err != nil {
					return err
				}
			}
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}
