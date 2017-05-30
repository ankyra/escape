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

package compiler

import (
	"github.com/ankyra/escape-core"
)

func (c *Compiler) compileDependencies(depends []string) error {

	result := []string{}
	for _, depend := range depends {
		dep, err := core.NewDependencyFromString(depend)
		if err != nil {
			return err
		}
		if err := c.ResolveVersion(dep, c.context); err != nil {
			return err
		}
		resolvedDep := dep.GetReleaseId()
		versionlessDep := dep.GetVersionlessReleaseId()
		metadata, err := c.context.GetDependencyMetadata(resolvedDep)
		if err != nil {
			return err
		}
		for _, consume := range metadata.GetConsumes() {
			c.metadata.AddConsumes(consume)
		}
		for _, input := range metadata.GetInputs() {
			if !input.HasDefault() {
				c.metadata.AddInputVariable(input)
			}
		}
		c.VariableCtx[versionlessDep] = metadata
		c.metadata.SetVariableInContext(versionlessDep, metadata.GetReleaseId())
		if dep.GetVariableName() != "" {
			c.VariableCtx[dep.GetVariableName()] = metadata
			c.metadata.SetVariableInContext(dep.GetVariableName(), metadata.GetReleaseId())
		}
		result = append(result, resolvedDep)
	}
	c.metadata.SetDependencies(result)
	return nil
}
