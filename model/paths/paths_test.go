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
	core "github.com/ankyra/escape-core"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type pathsSuite struct{}

var _ = Suite(&pathsSuite{})

func (s *pathsSuite) Test_Path_GetAppConfigDir_GetsCorrectDarwinDirectory(c *C) {
	unit := Path{
		os:            "darwin",
		homeDirectory: "/Users/test-user",
	}

	appConfigDir := unit.GetAppConfigDir()

	c.Assert(appConfigDir, Equals, "/Users/test-user/Library/Application Support/escape")
}

func (s *pathsSuite) Test_Path_GetAppConfigDir_GetsCorrectLinuxDirectory(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}

	appConfigDir := unit.GetAppConfigDir()

	c.Assert(appConfigDir, Equals, "/home/test-user/.config/escape")
}

func (s *pathsSuite) Test_Path_GetAppConfigDir_GetsCorrectWindowsDirectory(c *C) {
	unit := Path{
		os:            "windows",
		homeDirectory: "/Users/test-user",
	}

	appConfigDir := unit.GetAppConfigDir()

	c.Assert(appConfigDir, Equals, "/Users/test-user/escape")
}

func (s *pathsSuite) Test_Path_GetBaseDir(c *C) {
	unit := NewPathWithBaseDir("/base/")
	c.Assert(unit.GetBaseDir(), Equals, "/base/")
}

func (s *pathsSuite) Test_Path_EscapePlan(c *C) {
	unit := NewPathWithBaseDir("/base/")
	c.Assert(unit.EscapeDirectory(), Equals, "/base/.escape")
}

func (s *pathsSuite) Test_Path_ReleaseJson(c *C) {
	unit := Path{baseDir: "/base/dep/"}
	c.Assert(unit.ReleaseJson(), Equals, "/base/dep/release.json")
}

func (s *pathsSuite) Test_Path_NewPathForDependency(c *C) {
	unit := NewPathWithBaseDir("/base/")
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.NewPathForDependency(metadata).GetBaseDir(), Equals, "/base/deps/_/build")
}

func (s *pathsSuite) Test_Path_ScratchSpaceDirectory(c *C) {
	unit := NewPathWithBaseDir("/base/")
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.ScratchSpaceDirectory(metadata), Equals, "/base/.escape/build-v1")
}

func (s *pathsSuite) Test_Path_ScratchSpaceReleaseMetadata(c *C) {
	unit := NewPathWithBaseDir("/base/")
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.ScratchSpaceReleaseMetadata(metadata), Equals, "/base/.escape/build-v1/release.json")
}

func (s *pathsSuite) Test_Path_ReleaseTargetDirectory(c *C) {
	unit := NewPathWithBaseDir("/base/")
	c.Assert(unit.ReleaseTargetDirectory(), Equals, "/base/.escape/target")
}

func (s *pathsSuite) Test_Path_ReleaseLocation(c *C) {
	unit := NewPathWithBaseDir("/base/")
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.ReleaseLocation(metadata), Equals, "/base/.escape/target/build-v1.tgz")
}

func (s *pathsSuite) Test_Path_DependencyCacheDirectory(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}
	dependencyCacheDirectory := unit.DependencyCacheDirectory("_")
	c.Assert(dependencyCacheDirectory, Equals, "/home/test-user/.config/escape/_")
}

func (s *pathsSuite) Test_Path_DependecyReleaseArchive(c *C) {
	unit := NewPathWithBaseDir("/base/")
	dep, err := core.NewDependencyFromString("dep-v1.0")
	c.Assert(err, IsNil)
	c.Assert(unit.DependencyReleaseArchive(dep), Equals, "/base/dep-v1.0.tgz")
}

func (s *pathsSuite) Test_Path_Escape_Binary_Path(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}
	c.Assert(unit.EscapeBinaryPath(), Equals, "/home/test-user/.config/escape/.bin")
}

func (s *pathsSuite) Test_Path_DependecyDownloadTarget(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}
	dep, err := core.NewDependencyFromString("dep-v1.0")
	c.Assert(err, IsNil)
	c.Assert(unit.DependencyDownloadTarget(dep), Equals, "/home/test-user/.config/escape/_/dep-v1.0.tgz")
}

func (s *pathsSuite) Test_Path_LocalReleaseMetadata(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.LocalReleaseMetadata(metadata), Equals, "/home/test-user/.config/escape/_/build-v1.json")
}

func (s *pathsSuite) Test_Path_DepTypeDirectory(c *C) {
	unit := NewPathWithBaseDir("/base/")
	dep, err := core.NewDependencyFromString("dep-v1.0")
	c.Assert(err, IsNil)
	c.Assert(unit.DepTypeDirectory(dep), Equals, "/base/deps/_")
}

func (s *pathsSuite) Test_Path_UnpackedDepDirectory(c *C) {
	unit := NewPathWithBaseDir("/base/")
	dep, err := core.NewDependencyFromString("dep-v1.0")
	c.Assert(err, IsNil)
	c.Assert(unit.UnpackedDepDirectory(dep), Equals, "/base/deps/_/dep")
}

func (s *pathsSuite) Test_Path_UnpackedDepDirectoryReleaseMetadata(c *C) {
	unit := NewPathWithBaseDir("/base/")
	dep, err := core.NewDependencyFromString("dep-v1.0")
	c.Assert(err, IsNil)
	c.Assert(unit.UnpackedDepDirectoryReleaseMetadata(dep), Equals, "/base/deps/_/dep/release.json")
}

func (s *pathsSuite) Test_Path_ExtensionPath(c *C) {
	unit := NewPathWithBaseDir("/base/")
	metadata := core.NewReleaseMetadata("build", "1")
	c.Assert(unit.ExtensionPath(metadata, "test.sh"), Equals, "deps/_/build/test.sh")
}

func (s *pathsSuite) Test_Path_OutputsFile(c *C) {
	unit := NewPathWithBaseDir("/base/")
	c.Assert(unit.OutputsFile(), Equals, "/base/.escape/outputs.json")
}

func (s *pathsSuite) Test_Path_Script(c *C) {
	unit := NewPathWithBaseDir("/base/")
	c.Assert(unit.Script("test.sh"), Equals, "/base/test.sh")
}
