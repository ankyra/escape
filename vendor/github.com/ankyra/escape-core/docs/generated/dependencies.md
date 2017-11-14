---
date: 2017-11-11 00:00:00
title: "Dependencies"
slug: dependencies
type: "docs"
toc: true
---

Dependencies are configured in the (`depends`)[/docs/escape-plan/#depends]
field of the Escape plan.


Field | Type | Description
------|------|-------------
|release_id|`string`|The release id is required and is resolved at *build* time and then persisted in the release metadata ensuring that deployments always use the same versions. 
|||Examples: - To always use the latest version: `my-organisation/my-dependency-latest` - To always use version 0.1.1: `my-organisation/my-dependency-v0.1.1` - To always use the latest version in the 0.1 series: `my-organisation/my-dependency-v0.1.@` 
|build_mapping|`{string:any}`|
|deploy_mapping|`{string:any}`|
|consumes|`{string:string}`|
|scopes|`[string]`|

