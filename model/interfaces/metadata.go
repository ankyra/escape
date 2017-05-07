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

type ReleaseMetadata interface {
	ToJson() string
	ToDict() (map[string]interface{}, error)
	WriteJsonFile(string) error
	GetDirectories() []string
	GetVersionlessReleaseId() string
	GetReleaseId() string
	ToDependency() Dependency
	AddInputVariable(Variable)
	AddOutputVariable(Variable)
	SetVariableInContext(string, string)
	AddFileWithDigest(string, string)
	SetConsumes([]string)

	GetApiVersion() string
	GetBranch() string
	GetConsumes() []string
	GetDependencies() []string
	GetDescription() string
	GetErrands() map[string]Errand
	GetFiles() map[string]string
	GetRevision() string
	GetInputs() []Variable
	GetLogo() string
	GetMetadata() map[string]string
	GetName() string
	GetOutputs() []Variable
	GetPath() string
	GetProvides() []string
	GetType() string
	GetVersion() string
	GetVariableContext() map[string]string

	SetStage(stage, script string)
	GetScript(stage string) string
}
