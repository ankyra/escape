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
	"fmt"
)

func (c *Compiler) compileMetadata(metadata map[string]string) error {
	result := map[string]string{}
	for key, val := range metadata {
		str, err := RunScriptForCompileStep(val, c.VariableCtx)
		if err != nil {
			return fmt.Errorf("%s in metadata field.", err.Error())
		}
		result[key] = str
	}
	c.metadata.Metadata = result
	return nil
}
