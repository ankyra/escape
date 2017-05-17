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

package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"os/user"

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
)

type Path struct {
	os            string
	homeDirectory string
	baseDir       string
}

func NewPath() Paths {
	currentUser, _ := user.Current()
	currentDir, _ := os.Getwd()

	return &Path{
		os:            runtime.GOOS,
		homeDirectory: currentUser.HomeDir,
		baseDir:       currentDir,
	}
}

func NewPathWithBaseDir(baseDir string) Paths {
	p := NewPath()
	p.(*Path).baseDir = baseDir
	return p
}

func (p *Path) GetBaseDir() string {
	return p.baseDir
}

func (p *Path) NewPathForDependency(metadata *core.ReleaseMetadata) Paths {
	return NewPathWithBaseDir(filepath.Join(p.baseDir, "deps", metadata.GetName()))
}

func (p *Path) GetAppConfigDir() string {
	if p.os == "windows" {
		folder := os.Getenv("APPDATA")
		if folder == "" {
			folder = p.homeDirectory
		}
		return filepath.Join(folder, "escape")
	} else if p.os == "darwin" {
		return filepath.Join(p.homeDirectory, "Library", "Application Support", "escape")
	}
	folder := os.Getenv("XDG_CONFIG_HOME")
	if folder == "" {
		folder = filepath.Join(p.homeDirectory, ".config")
	}
	return filepath.Join(folder, "escape")
}

func (p *Path) EscapeDirectory() string {
	return filepath.Join(p.baseDir, ".escape")
}

func (p *Path) ReleaseJson() string {
	return filepath.Join(p.baseDir, "release.json")
}

func (p *Path) ScratchSpaceDirectory(metadata *core.ReleaseMetadata) string {
	return filepath.Join(p.EscapeDirectory(), metadata.GetReleaseId())
}

func (p *Path) ScratchSpaceReleaseMetadata(metadata *core.ReleaseMetadata) string {
	return filepath.Join(p.ScratchSpaceDirectory(metadata), "release.json")
}

func (p *Path) ReleaseTargetDirectory() string {
	return filepath.Join(p.EscapeDirectory(), "target")
}

func (p *Path) ReleaseLocation(metadata *core.ReleaseMetadata) string {
	pkgName := metadata.GetReleaseId()
	return filepath.Join(p.ReleaseTargetDirectory(), pkgName+".tgz")
}

func (p *Path) DependencyCacheDirectory() string {
	return filepath.Join(p.GetAppConfigDir(), "deps")
}

func (p *Path) DependencyReleaseArchive(dependency *core.Dependency) string {
	return filepath.Join(p.baseDir, dependency.GetReleaseId()+".tgz")
}

func (p *Path) DependencyDownloadTarget(dependency *core.Dependency) string {
	return filepath.Join(p.DependencyCacheDirectory(), dependency.GetReleaseId()+".tgz")
}

func (p *Path) LocalReleaseMetadata(metadata *core.ReleaseMetadata) string {
	return filepath.Join(p.DependencyCacheDirectory(), metadata.GetReleaseId()+".json")
}

func (p *Path) ExtensionPath(extension *core.ReleaseMetadata, path string) string {
	return filepath.Join("deps", extension.GetName(), path)
}

func (p *Path) DepTypeDirectory(dependency *core.Dependency) string {
	return filepath.Join(p.baseDir, "deps")
}
func (p *Path) UnpackedDepDirectory(dependency *core.Dependency) string {
	return filepath.Join(p.DepTypeDirectory(dependency), dependency.GetBuild())
}

func (p *Path) UnpackedDepDirectoryReleaseMetadata(dependency *core.Dependency) string {
	return filepath.Join(p.UnpackedDepDirectory(dependency), "release.json")
}

func (p *Path) OutputsFile() string {
	return filepath.Join(p.EscapeDirectory(), "outputs.json")
}

func (p *Path) Script(script string) string {
	return filepath.Join(p.baseDir, script)
}

func (p *Path) EnsureEscapeDirectoryExists() error {
	return util.MkdirRecursively(p.EscapeDirectory())
}

func (p *Path) EnsureDependencyCacheDirectoryExists() error {
	return util.MkdirRecursively(p.DependencyCacheDirectory())
}

func (p *Path) EnsureDependencyTypeDirectoryExists(dep *core.Dependency) error {
	return util.MkdirRecursively(p.DepTypeDirectory(dep))
}

func (p *Path) EnsureScratchSpaceDirectoryExists(metadata *core.ReleaseMetadata) error {
	return util.MkdirRecursively(p.ScratchSpaceDirectory(metadata))
}

func (p *Path) EnsureReleaseTargetDirectoryExists() error {
	return util.MkdirRecursively(p.ReleaseTargetDirectory())
}
