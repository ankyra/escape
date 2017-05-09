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

type EscapePlan interface {
	Init(typ, buildId string) EscapePlan
	LoadConfig(string) error
	ToYaml() []byte
	GetReleaseId() string
	GetVersionlessReleaseId() string

	GetBuild() string
	GetConsumes() []string
	GetDepends() []string
	GetDescription() string
	GetErrands() map[string]interface{}
	GetIncludes() []string
	GetInputs() []interface{}
	GetLogo() string
	GetMetadata() map[string]string
	GetOutputs() []interface{}
	GetPath() string
	GetPostBuild() string
	GetPostDeploy() string
	GetPostDestroy() string
	GetPreBuild() string
	GetPreDeploy() string
	GetPreDestroy() string
	GetProvides() []string
	GetTest() string
	GetSmoke() string
	GetType() string
	GetVersion() string
	GetTemplates() []interface{}

	SetBuild(string)
	SetConsumes([]string)
	SetDepends([]string)
	SetDescription(string)
	SetErrands(map[string]interface{})
	SetIncludes([]string)
	SetInputs([]interface{})
	SetLogo(string)
	SetMetadata(map[string]string)
	SetOutputs([]interface{})
	SetPath(string)
	SetPostBuild(string)
	SetPostDeploy(string)
	SetPostDestroy(string)
	SetPreBuild(string)
	SetPreDeploy(string)
	SetPreDestroy(string)
	SetProvides([]string)
	SetTest(string)
	SetType(string)
	SetVersion(string)
}
