---
date: 2017-11-11 00:00:00
title: "Providers and Consumers"
slug: providers-and-consumers
type: "docs"
toc: true
wip: true
contributeLink: https://github.com/ankyra/escape-core/blob/master/consumer.go
---

Unlike Dependencies, which are resolved at build time and provide tight
coupling, we can use Consumers and Providers to resolve and loosely couple
dependencies at deployment time.

## Escape Plan

Consumers are configured in the [`consumes`](/docs/escape-plan/#consumes)
field of the Escape Plan.

Providers are configured in the [`provides`](/docs/escape-plan/#provides)
field of the Escape Plan.


Field | Type | Description
------|------|-------------
|name|`string`|
|scopes|`[string]`|

