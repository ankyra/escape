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

package compiler

import (
	"github.com/ankyra/escape-core"
)

func compileErrands(ctx *CompilerContext) error {
	for name, errandDict := range ctx.Plan.Errands {
		newErrand, err := core.NewErrandFromDict(name, errandDict)
		if err != nil {
			return err
		}
		if newErrand.Run.RelativeScript != "" {
			if err := ctx.AddFileDigest(newErrand.Run.RelativeScript); err != nil {
				return err
			}
		}
		ctx.Metadata.Errands[name] = newErrand
	}
	return nil
}
