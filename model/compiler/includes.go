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
	"path/filepath"
)

func (c *Compiler) compileIncludes(includes []string) {
	for _, globPattern := range includes {
		paths, err := filepath.Glob(globPattern)
		if err != nil {
			fmt.Println("Warning: ignoring pattern error: " + err.Error())
			continue
		}
		if paths == nil {
			continue
		}
		for _, path := range paths {
			err = c.addFileDigest(path)
			if err != nil {
				fmt.Println("Ignoring problem with path " + path + ": " + err.Error())
			}
		}
	}
}
