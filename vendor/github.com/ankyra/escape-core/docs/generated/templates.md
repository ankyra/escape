---
date: 2017-11-11 00:00:00
title: "Templates"
slug: templates
type: "docs"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/templates/templates.go
---

Escape provides the Mustache templating language and integrates it with the
package's [Variables](/docs/reference/input-and-output-variables/), making for a quick
and easy way to render files at either build or deploy time.

## Escape Plan

Templates are configured in the Escape Plan under the
[`templates`](/docs/reference/escape-plan/#templates) field.


Field | Type | Description
------|------|-------------
|file|`string`|The file containing the template. This field is required. 
|target|`string`|The target location for the rendered template. If the source location specified in `file` has the `.tpl` extension this `target` will default to source location minus that extension. 
|||For example: if `file` is `"hello.txt.tpl"` then the default value for target will be `"hello.txt"` 
|scopes|`scopes.Scopes`|A list of scopes (`build`, `deploy`) that defines during which stage(s) the template should be rendered. 
|mapping|`{string:any}`|This mapping can be used to relate template variables to Escape variables. 

