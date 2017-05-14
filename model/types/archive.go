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

package types

import (
	"fmt"
	plan "github.com/ankyra/escape-client/model/escape_plan"
	. "github.com/ankyra/escape-client/model/interfaces"
	"io/ioutil"
	"path/filepath"
)

type ArchiveReleaseType struct{}

func (a *ArchiveReleaseType) GetType() string {
	return "archive"
}

func (a *ArchiveReleaseType) InitEscapePlan(p *plan.EscapePlan) {
	p.SetIncludes([]string{})
	queue := []string{"."}
	for len(queue) != 0 {
		dirPath := queue[0]
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, file := range files {
			joinedPath := filepath.Join(dirPath, file.Name())
			if file.IsDir() {
				queue = append(queue, joinedPath)
			} else if file.Name() == "release.json" && dirPath == "." {
				continue
			} else if file.Name()[0] == '.' || joinedPath[0] == '.' {
				continue
			} else {
				newIncludes := append(p.GetIncludes(), joinedPath)
				p.SetIncludes(newIncludes)
			}
		}
		queue = queue[1:]
	}
}

func (a *ArchiveReleaseType) CompileMetadata(p *plan.EscapePlan, release ReleaseMetadata) error {
	return nil
}

func (a *ArchiveReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	return nil, nil
}

func (a *ArchiveReleaseType) Destroy(ctx RunnerContext) error {
	return nil
}
