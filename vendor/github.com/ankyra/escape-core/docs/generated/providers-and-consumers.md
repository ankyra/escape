---
date: 2017-11-11 00:00:00
title: "Providers and Consumers"
slug: providers-and-consumers
type: "docs"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/consumer.go
---

Unlike Dependencies, which are resolved at build time and provide tight
coupling, we can use Consumers and Providers to resolve and loosely couple
dependencies at deployment time.

To signal that a package implements a certain interface, e.g. "my-interface", we can
define it as a provider in the Escape plan:

```yaml
provides:
- my-interface
```

Packages that require a "my-interface" can define that in their Escape Plan as well:

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

## Escape Plan

Consumers are configured in the [`consumes`](/docs/reference/escape-plan/#consumes)
field of the Escape Plan.

Providers are configured in the [`provides`](/docs/reference/escape-plan/#provides)
field of the Escape Plan.


Field | Type | Description
------|------|-------------
|name|`string`|The name of the interface. Can be renamed using the `as` syntax. For example: `kubernetes as k8s`, `postgres`, `postgres as db` 
|scopes|`[string]`|A list of scopes (`build`, `deploy`) that defines during which stage(s) this dependency should be fetched and deployed. Also see [`build_consumes`](/docs/reference/escape-plan/#build_consumes] and [`deploy_consumes`](/docs/reference/escape-plan/#deploy_consumes]. 
|variable|`string`|The variable used to reference this consumer. Overwriting this field in the Escape plan has no effect. 

