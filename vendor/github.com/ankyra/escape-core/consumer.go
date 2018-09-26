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

package core

import (
	"fmt"

	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/scopes"
)

/*

Unlike Dependencies, which are resolved at build time and provide tight
coupling, we can use Consumers and Providers to resolve and loosely couple
packages at deployment time. Providers make their output variables available to
each consumer, making it possible to share credentials and host details for
example. Providers and Consumers are often used to model the different layers
in an architecture; where the layer below is consumed by the layer on top (e.g.
AWS -> Kubernetes -> Helm -> Service).

To signal that a package implements a certain interface, e.g. "my-interface", we can
define it as a provider in the Escape plan:

```yaml
provides:
- my-interface
```

Packages that require a "my-interface" define this joyful fact in their Escape
Plan as well:

```yaml
consumes:
- my-interface
```

When building or deploying the consumer Escape now makes sure that it also has
access to a provider's output variables. You can only link consumers to
providers in the same environment. Escape will link up consumers with providers
automatically if there's only a single provider of a particular interface; other
times providers need to be specified with the `-p` flag. For example:

```
escape run deploy my-project/my-consumer-v1.0.0 -p my-interface=provider-deployment
```

To list providers in an environment you can use the [`escape state
show-providers`](/docs/reference/escape_state_show-providers/) command.

## Wrapper Packages

Providers and consumers provide a loose coupling, but sometimes we know exactly
what provider implementation we want to use. In this case we can create a wrapper
release that uses one dependency as the provider for the next:

```yaml
depends:
- release_id: my-project/postgres-provider-latest as postgres
- release_id: my-project/my-application-latest
  consumes:
	  postgres: $postgres.deployment
```

To read more about wrapper releases see the [blog post](https://www.ankyra.io/blog/combining-packages-into-platforms/).

## Provider Activation and Deactivation

When a package consumes another package as a provider, the provider has the
ability to run activation and deactivation scripts. The scripts can be defined by
adding the following fields to the Escape plan:

```yaml
activate_provider: activate.sh
deactivate_provider: deactivate.sh
```

The scripts gets full access to the provider's deployment state and is in that
way similar to running a smoke test. These steps are often used to activate
credentials, install packages, or otherwise manage state on the deployment
machine or container.

To disable activation and deactivation see the `skip_activate` and `skip_deactivate` options
on the consumer configuration below.

## Escape Plan

Consumers are configured in the [`consumes`](/docs/reference/escape-plan/#consumes)
field of the Escape Plan.

Providers are configured in the [`provides`](/docs/reference/escape-plan/#provides)
field of the Escape Plan.

*/
type ConsumerConfig struct {
	// The name of the interface. Can be renamed using the `as` syntax.
	// For example: `kubernetes as k8s`, `postgres`, `postgres as db`
	Name string `json:"name" yaml:"name"`

	// A list of scopes (`build`, `deploy`) that defines during which stage(s)
	// this dependency should be fetched and deployed. Also see
	// [`build_consumes`](/docs/reference/escape-plan/#build_consumes] and
	// [`deploy_consumes`](/docs/reference/escape-plan/#deploy_consumes].
	Scopes scopes.Scopes `json:"scopes" yaml:"scopes"`

	// The variable used to reference this consumer. Overwriting this field in
	// the Escape plan has no effect.
	VariableName string `json:"variable" yaml:"variable"`

	// Skips the provider's activation step.
	SkipActivate bool `json:"skip_activate" yaml:"skip_activate"`

	// Skips the provider's deactivation step. Only relevant when
	// `skip_activate` is false.
	SkipDeactivate bool `json:"skip_deactivate" yaml:"skip_deactivate"`
}

// Only used for testing purposes.
//
func NewConsumerConfig(name string) *ConsumerConfig {
	return &ConsumerConfig{
		Name:           name,
		Scopes:         scopes.AllScopes,
		VariableName:   name,
		SkipActivate:   false,
		SkipDeactivate: false,
	}
}

func NewConsumerConfigFromString(str string) (*ConsumerConfig, error) {
	id, err := parsers.ParseConsumer(str)
	if err != nil {
		return nil, err
	}
	cfg := NewConsumerConfig(id.Interface)
	if id.VariableName != "" {
		cfg.VariableName = id.VariableName
	}
	return cfg, nil
}

func NewConsumerConfigFromInterface(v interface{}) (*ConsumerConfig, error) {
	switch v.(type) {
	case string:
		return NewConsumerConfigFromString(v.(string))
	case map[interface{}]interface{}:
		return NewConsumerConfigFromMap(v.(map[interface{}]interface{}))
	}
	return nil, fmt.Errorf("Expecting dict or string type")
}

func NewConsumerConfigFromMap(dep map[interface{}]interface{}) (*ConsumerConfig, error) {
	var name string
	var skipActivate bool
	var skipDeactivate bool
	consumesScopes := scopes.Scopes{}
	for key, val := range dep {
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string key in consumer")
		}
		if keyStr == "name" {
			valString, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("Expecting string for consumer 'name' got '%T'", val)
			}
			name = valString
		} else if key == "scopes" {
			s, err := scopes.NewScopesFromInterface(val)
			if err != nil {
				return nil, err
			}
			consumesScopes = s
		} else if key == "skip_activate" {
			valBool, ok := val.(bool)
			if !ok {
				return nil, fmt.Errorf("Expecting bool for consumer field 'skip_activate' got '%T'", val)
			}
			skipActivate = valBool
		} else if key == "skip_deactivate" {
			valBool, ok := val.(bool)
			if !ok {
				return nil, fmt.Errorf("Expecting bool for consumer field 'skip_deactivate' got '%T'", val)
			}
			skipDeactivate = valBool
		}
	}
	if name == "" {
		return nil, fmt.Errorf("Missing 'name' in consumer")
	}
	cfg, err := NewConsumerConfigFromString(name)
	if err != nil {
		return nil, err
	}
	cfg.Scopes = consumesScopes
	cfg.SkipActivate = skipActivate
	cfg.SkipDeactivate = skipDeactivate
	return cfg, cfg.ValidateAndFix()
}

func (c *ConsumerConfig) ValidateAndFix() error {
	if c.Scopes == nil || len(c.Scopes) == 0 {
		c.Scopes = []string{"build", "deploy"}
	}
	if c.Name == "" {
		return fmt.Errorf("Missing name for Consumer")
	}
	if c.VariableName == "" {
		c.VariableName = c.Name
	}
	return nil
}
func (c *ConsumerConfig) Copy() *ConsumerConfig {
	result := NewConsumerConfig(c.Name)
	result.Scopes = c.Scopes.Copy()
	result.VariableName = c.VariableName
	result.SkipActivate = c.SkipActivate
	result.SkipDeactivate = c.SkipDeactivate
	return result
}

func (c *ConsumerConfig) InScope(scope string) bool {
	return c.Scopes.InScope(scope)
}
