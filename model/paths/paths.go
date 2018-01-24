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

package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"os/user"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape/util"
)

type Path struct {
	os            string
	homeDirectory string
	baseDir       string
}

func NewPath() *Path {
	currentUser, _ := user.Current()
	currentDir, _ := os.Getwd()

	return &Path{
		os:            runtime.GOOS,
		homeDirectory: currentUser.HomeDir,
		baseDir:       currentDir,
	}
}

func NewPathWithBaseDir(baseDir string) *Path {
	p := NewPath()
	p.baseDir = baseDir
	return p
}

func (p *Path) NewPathForDependency(metadata *core.ReleaseMetadata) *Path {
	return NewPathWithBaseDir(filepath.Join(p.baseDir, "deps", metadata.Project, metadata.Name))
}

func (p *Path) GetBaseDir() string {
	return p.baseDir
}

func (p *Path) GetAppConfigDir() string {
	return util.GetAppConfigDir(p.os, p.homeDirectory)
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

func (p *Path) EscapeBinaryPath() string {
	return filepath.Join(p.GetAppConfigDir(), ".bin")
}

func (p *Path) DependencyCacheDirectory(project string) string {
	return filepath.Join(p.GetAppConfigDir(), project)
}

func (p *Path) GetDefaultStateLocation() string {
	return filepath.Join(p.GetAppConfigDir(), "escape_state.json")
}

func (p *Path) DependencyReleaseArchive(dependency *core.Dependency) string {
	return filepath.Join(p.baseDir, dependency.GetReleaseId()+".tgz")
}

func (p *Path) DependencyDownloadTarget(dependency *core.Dependency) string {
	return filepath.Join(p.DependencyCacheDirectory(dependency.Project), dependency.GetReleaseId()+".tgz")
}

func (p *Path) LocalReleaseMetadata(metadata *core.ReleaseMetadata) string {
	return filepath.Join(p.DependencyCacheDirectory(metadata.Project), metadata.GetReleaseId()+".json")
}

func (p *Path) DepTypeDirectory(dependency *core.Dependency) string {
	return filepath.Join(p.baseDir, "deps", dependency.Project)
}
func (p *Path) UnpackedDepDirectory(dependency *core.Dependency) string {
	return filepath.Join(p.DepTypeDirectory(dependency), dependency.Name)
}
func (p *Path) UnpackedDepCfgDirectory(dependency *core.DependencyConfig) string {
	return filepath.Join(p.baseDir, "deps", dependency.Project, dependency.Name)
}
func (p *Path) UnpackedDepDirectoryByReleaseMetadata(metadata *core.ReleaseMetadata) string {
	return filepath.Join(p.baseDir, "deps", metadata.Project, metadata.Name)
}
func (p *Path) UnpackedDepDirectoryReleaseMetadata(dependency *core.Dependency) string {
	return filepath.Join(p.UnpackedDepDirectory(dependency), "release.json")
}

func (p *Path) ExtensionPath(extension *core.ReleaseMetadata, path string) string {
	return filepath.Join("deps", extension.Project, extension.Name, path)
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

func (p *Path) EnsureEscapeConfigDirectoryExists() error {
	return util.MkdirRecursively(p.GetAppConfigDir())
}

func (p *Path) EnsureDependencyCacheDirectoryExists(project string) error {
	return util.MkdirRecursively(p.DependencyCacheDirectory(project))
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

func (p *Path) EnsureEscapePathDirectoryExists() error {
	return util.MkdirRecursively(p.EscapeBinaryPath())
}
