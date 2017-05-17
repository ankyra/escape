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

import (
	core "github.com/ankyra/escape-core"
)

type Paths interface {
	GetBaseDir() string
	GetAppConfigDir() string
	EscapeDirectory() string
	ReleaseJson() string
	ScratchSpaceDirectory(*core.ReleaseMetadata) string
	ScratchSpaceReleaseMetadata(*core.ReleaseMetadata) string
	ReleaseTargetDirectory() string
	ReleaseLocation(*core.ReleaseMetadata) string
	DependencyCacheDirectory() string
	DependencyReleaseArchive(*core.Dependency) string
	DependencyDownloadTarget(*core.Dependency) string
	OutputsFile() string
	DepTypeDirectory(*core.Dependency) string
	UnpackedDepDirectory(*core.Dependency) string
	UnpackedDepDirectoryReleaseMetadata(*core.Dependency) string
	LocalReleaseMetadata(*core.ReleaseMetadata) string
	EnsureEscapeDirectoryExists() error
	ExtensionPath(*core.ReleaseMetadata, string) string

	EnsureDependencyTypeDirectoryExists(*core.Dependency) error
	EnsureScratchSpaceDirectoryExists(*core.ReleaseMetadata) error
	EnsureReleaseTargetDirectoryExists() error
	EnsureDependencyCacheDirectoryExists() error

	Script(script string) string

	NewPathForDependency(*core.ReleaseMetadata) Paths
}
