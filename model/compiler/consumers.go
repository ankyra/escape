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

	core "github.com/ankyra/escape-core"
)

func compileConsumers(ctx *CompilerContext) error {
	for _, consumer := range ctx.Plan.Consumes {
		c, err := core.NewConsumerConfigFromInterface(consumer)
		if err != nil {
			return fmt.Errorf("%s in 'consumes' field", err)
		}
		ctx.Metadata.AddConsumes(c)
	}
	for _, consumer := range ctx.Plan.BuildConsumes {
		c, err := core.NewConsumerConfigFromInterface(consumer)
		if err != nil {
			return fmt.Errorf("%s in 'build_consumes' field", err)
		}
		c.Scopes = []string{"build"}
		ctx.Metadata.AddConsumes(c)
	}
	for _, consumer := range ctx.Plan.DeployConsumes {
		c, err := core.NewConsumerConfigFromInterface(consumer)
		if err != nil {
			return fmt.Errorf("%s in 'deploy_consumes' field", err)
		}
		c.Scopes = []string{"deploy"}
		ctx.Metadata.AddConsumes(c)
	}
	return nil
}
