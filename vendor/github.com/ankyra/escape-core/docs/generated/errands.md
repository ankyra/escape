---
date: 2017-11-11 00:00:00
title: "Errands"
slug: errands
type: "docs"
toc: true
wip: true
contributeLink: https://github.com/ankyra/escape-core/blob/master/errand.go
---

Errands are an Escape mechanism that make it easy to run operational and
publication tasks against deployed packages. They can be used to implement
backup procedures, user management, scalability controls, binary
publications, etc. Errands are a good idea whenever a task needs to be aware
of Environments.

You can inspect and run Errands using the [`escape
errands`](/docs/escape_errands/) command.

## Escape Plan

Errands are configured in the Escape Plan under the
[`errands`](/docs/escape-plan/#errands) field.


Field | Type | Description
------|------|-------------
|name|`string`|The name of the errand. This field is required. 
|description|`string`|An optional description of the errand. 
|script|`string`|The location of the script performing the actual work. 
|||The script has access to the deployment inputs and outputs as enviroment variables. For example: an input with `"id": "input_variable"` will be accessible as `INPUT_input_variable`; and an output with `"id": "output_variable"` as `OUTPUT_output_variable`. 
|inputs|`[variables.Variable]`|A list of [Variables](/docs/input-and-output-variables/_. The values will be made available to the `script` (along with the regular deployment inputs and outputs) as environment variables. For example: a variable with `"id": "input_variable"` will be accessible as environment variable `INPUT_input_variable` 

