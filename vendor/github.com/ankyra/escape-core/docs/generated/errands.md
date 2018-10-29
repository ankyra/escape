---
date: 2017-11-11 00:00:00
title: "Errands"
slug: errands
type: "reference"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/errand.go
---

Errands are an Escape mechanism that make it easy to run operational and
publication tasks against deployed packages. They can be used to implement
backup procedures, user management, scalability controls, binary
publications, etc. Errands are a good idea whenever a task needs to be aware
of Environments.

You can inspect and run Errands using the [`escape
errands`](/docs/reference/escape_errands/) command.

## Escape Plan

Errands are configured in the Escape Plan under the
[`errands`](/docs/reference/escape-plan/#errands) field.


Field | Type | Description
------|------|-------------
|name|`string`|The name of the errand. This field is required. 
|description|`string`|An optional description of the errand. 
|script|`string`|The script or command performing the errand (deprecated, use 'run' instead). 
|||The script has access to the deployment inputs and outputs as enviroment variables. For example: an input with `"id": "input_variable"` will be accessible as `INPUT_input_variable`; and an output with `"id": "output_variable"` as `OUTPUT_output_variable`. 
|exec_stage|`ExecStage`|The script or command performing the errand. 
|||The command has access to the deployment inputs and outputs as enviroment variables. For example: an input with `"id": "input_variable"` will be accessible as `INPUT_input_variable`; and an output with `"id": "output_variable"` as `OUTPUT_output_variable`. 
|inputs|`[variables.Variable]`|A list of [Variables](/docs/reference/input-and-output-variables/). The values will be made available to the `script` (along with the regular deployment inputs and outputs) as environment variables. For example: a variable with `"id": "input_variable"` will be accessible as environment variable `INPUT_input_variable` 

