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

package interfaces

type Paths interface {
	GetBaseDir() string
	GetAppConfigDir() string
	EscapeDirectory() string
	ReleaseJson() string
	ScratchSpaceDirectory(ReleaseMetadata) string
	ScratchSpaceReleaseMetadata(ReleaseMetadata) string
	ReleaseTargetDirectory() string
	ReleaseLocation(ReleaseMetadata) string
	DependencyCacheDirectory() string
	DependencyReleaseArchive(Dependency) string
	DependencyDownloadTarget(Dependency) string
	OutputsFile() string
	DepTypeDirectory(Dependency) string
	UnpackedDepDirectory(Dependency) string
	UnpackedDepDirectoryReleaseMetadata(Dependency) string
	LocalReleaseMetadata(ReleaseMetadata) string
	EnsureEscapeDirectoryExists() error
	ExtensionPath(ReleaseMetadata, string) string

	EnsureDependencyTypeDirectoryExists(Dependency) error
	EnsureScratchSpaceDirectoryExists(ReleaseMetadata) error
	EnsureReleaseTargetDirectoryExists() error
	EnsureDependencyCacheDirectoryExists() error

	Script(script string) string

	NewPathForDependency(ReleaseMetadata) Paths
}
