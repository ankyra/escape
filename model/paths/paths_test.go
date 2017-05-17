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
	"testing"

	"os"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type pathsSuite struct{}

var unit Path
var workingDirectory string

var _ = Suite(&pathsSuite{})

func (*pathsSuite) SetUpSuite(c *C) {
	workingDirectory, _ = os.Getwd()
}

func (s *pathsSuite) TearDownTest(c *C) {
	os.RemoveAll(workingDirectory + "/.escape")
}

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

func (s *pathsSuite) Test_Path_EscapePlan(c *C) {
	unit := Path{}

	escapeDirectory := unit.EscapeDirectory()

	c.Assert(escapeDirectory, Equals, ".escape")
}

//func (s *pathsSuite) Test_Path_ScratchSpaceDirectory(c *C) {
//	unit := Path{}
//
//	scratchSpaceDirectory := unit.ScratchSpaceDirectory(&ReleaseMetadata{
//		Type:    "type",
//		Name:    "name",
//		Version: "1",
//	})
//
//	c.Assert(scratchSpaceDirectory, Equals, ".escape/type-name-v1")
//}
//
//func (s *pathsSuite) Test_Path_ScratchSpaceReleaseMetadata(c *C) {
//	unit := Path{}
//
//	scratchSpaceReleaseMetadata := unit.ScratchSpaceReleaseMetadata(&ReleaseMetadata{
//		Type:    "type",
//		Name:    "name",
//		Version: "1",
//	})
//
//	c.Assert(scratchSpaceReleaseMetadata, Equals, ".escape/type-name-v1/release.json")
//}

func (s *pathsSuite) Test_Path_ReleaseTargetDirectory(c *C) {
	unit := Path{}

	releaseTargetDirectory := unit.ReleaseTargetDirectory()

	c.Assert(releaseTargetDirectory, Equals, ".escape/target")
}

//func (s *pathsSuite) Test_Path_ReleaseLocation(c *C) {
//	unit := Path{}
//
//	releaseLocation := unit.ReleaseLocation(&ReleaseMetadata{
//		Type:    "type",
//		Name:    "name",
//		Version: "1",
//	})
//
//	c.Assert(releaseLocation, Equals, ".escape/target/type-name-v1.tgz")
//}

func (s *pathsSuite) Test_Path_DependencyCacheDirectory(c *C) {
	unit := Path{
		os:            "linux",
		homeDirectory: "/home/test-user",
	}

	dependencyCacheDirectory := unit.DependencyCacheDirectory()

	c.Assert(dependencyCacheDirectory, Equals, "/home/test-user/.config/escape/deps")
}

//func (s *pathsSuite) Test_Path_DependencyReleaseArchive(c *C) {
//	unit := Path{}
//
//	dependencyReleaseArchive := unit.DependencyReleaseArchive(&Dependency{
//		Type:    "type",
//		Build:   "build",
//		Version: "1",
//	})
//
//	c.Assert(dependencyDownloadTarget, Equals, ".config/escape/deps/type-build-v1.tgz")
//}

//func (s *pathsSuite) Test_Path_UnpackedDepDirectory(c *C) {
//	unit := Path{}
//
//	unpackedDepDirectory := unit.UnpackedDepDirectory(&Dependency{
//		Type:    "type",
//		Build:   "build",
//		Version: "1",
//	})
//
//	c.Assert(unpackedDepDirectory, Equals, "type/build")
//}
//
//func (s *pathsSuite) Test_Path_UnpackedDepDirectoryReleaseMetadata(c *C) {
//	unit := Path{}
//
//	unpackedDepDirectoryReleaseMetadata := unit.UnpackedDepDirectoryReleaseMetadata(&Dependency{
//		Type:    "type",
//		Build:   "build",
//		Version: "1",
//	})
//
//	c.Assert(unpackedDepDirectoryReleaseMetadata, Equals, "type/build/release.json")
//}

//func (s *pathsSuite) Test_Path_EnsureEscapeDirectoryExists_CreatesDirectory(c *C) {
//	unit.EnsureEscapeDirectoryExists()
//
//	_, err := os.Stat(workingDirectory + ".escape")
//
//	c.Assert(err, IsNil)
//}
//
//func (s *pathsSuite) Test_Path_EnsureEscapeDirectoryExists_WithExistDir(c *C) {
//	os.Mkdir(workingDirectory + ".escape", 0777)
//
//	err := unit.EnsureEscapeDirectoryExists()
//	c.Assert(err, IsNil)
//}
