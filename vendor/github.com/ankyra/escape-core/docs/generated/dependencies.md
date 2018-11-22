---
date: 2017-11-11 00:00:00
title: "Dependencies"
slug: dependencies
type: "reference"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/dependency_config.go
---

## Escape Plan

Dependencies are configured in the [`depends`](/docs/reference/escape-plan/#depends)
field of the Escape plan.


Field | Type | Description
------|------|-------------
|release_id|`string`|The release id is required and is resolved at *build* time and then persisted in the release metadata ensuring that deployments always use the same versions. 
|||Examples: - To always use the latest version: `my-organisation/my-dependency-latest` - To always use version 0.1.1: `my-organisation/my-dependency-v0.1.1` - To always use the latest version in the 0.1 series: `my-organisation/my-dependency-v0.1.@` - To make it possible to reference a dependency using a different name: `my-organisation/my-dependency-latest as my-name` 
|mapping|`{string:any}`|Define the values of dependency inputs using Escape Script. 
|build_mapping|`{string:any}`|Define the values of dependency inputs using Escape Script when running stages in the build scope. 
|deploy_mapping|`{string:any}`|Define the values of dependency inputs using Escape Script when running stages in the deploy scope. 
|consumes|`{string:string}`|Map providers from the parent to dependencies. 
|||Example: ``` consumes: - my-provider depends: - release_id: my-org/my-dep-latest consumes: provider: $my-provider.deployment ``` 
|deployment_name|`string`|The name of the (sub)-deployment. This defaults to the versionless release id; e.g. if the release_id is `my-org/my-dep-v1.0` then the DeploymentName will be `my-org/my-dep` by default. 
|variable|`string`|The variable used to reference this dependency. By default the variable name is the versionless release id of the dependency, but this can be overruled by renaming the dependency (e.g. `my-org/my-release-latest as my-variable`. This field will be set automatically at build time. Overwriting this field in the Escape plan has no effect. 
|scopes|`scopes.Scopes`|A list of scopes (`build`, `deploy`) that defines during which stage(s) this dependency should be fetched and deployed. *Currently not implemented!* 
|-|`string`|Parsed out of the release ID. For example: when release id is `"my-org/my-name-v1.0"` this value is `"my-org"`. 
|-|`string`|Parsed out of the release ID. For example: when release id is `"my-org/my-name-v1.0"` this value is `"my-name"`. 
|-|`string`|Parsed out of the release ID. For example: when release id is `"my-org/my-name-v1.0"` this value is `"1.0"`. 
|-|`string`|Parsed out of the release ID. For example: when release id is `"my-org/my-name:tag"` this value is `"tag"`. 

