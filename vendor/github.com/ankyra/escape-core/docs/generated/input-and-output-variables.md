---
date: 2017-11-11 00:00:00
title: "Input and Output Variables"
slug: input-and-output-variables
type: "docs"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/variables/variable.go
---

Variables can be used to defined inputs and outputs for the build and
deployment stages. They can also be used to make [Errands](/docs/reference/errands/)
configurable.

Variables are strongly typed, which is checked at both build and deploy
time.  A task can't succeed if the required variables have not been
configured correctly.

## Escape Plan

Variables can be configured in the Escape Plan under the
[`inputs`](/docs/reference/escape-plan/#inputs),
[`build_inputs`](/docs/reference/escape-plan/#build_inputs),
[`deploy_inputs`](/docs/reference/escape-plan/#deploy_inputs) and
[`outputs`](/docs/reference/escape-plan/#outputs) fields.


Field | Type | Description
------|------|-------------
|id|`string`|A unique name for this variable. Required field. 
|type|`string`|The variable type. Before executing any steps Escape will make sure that all the values match the types that are set on the variables. 
|||One of: `string`, `list`, `integer`, `bool`. 
|||Default: `string` 
|default|`any`|A default value for this variable. This value will be used if no value has been specified by the user. 
|description|`string`|A description of the variable. 
|friendly|`string`|A friendly name for this variable for presentational purposes only. 
|visible|`bool`|Control whether or not this variable should be visible when deploying interactively. In other words: should the user be asked to input this value?  It only really makes sense to set this to `true` if there a `default` is set. 
|options|`{string:any}`|Options that put more constraints on the type. 
|sensitive|`bool`|Is this sensitive data? 
|items|`any`|If set, this should contain all the valid values for this variable. 
|eval_before_dependencies|`bool`|Should the variables be evaluated before the dependencies are deployed? 
|scopes|`[string]`|A list of scopes (`build`, `deploy`) that defines during which stage(s) this variable should be active. You wouldn't usually use this field directly, but use something like [`build_inputs`](/docs/escape-plan/#build_inputs) or [`deploy_inputs`](/docs/escape-plan/#deploy_inputs), which usually express intent better. 

