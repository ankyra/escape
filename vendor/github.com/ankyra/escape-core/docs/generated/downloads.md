---
date: 2017-11-11 00:00:00
title: "Downloads"
slug: downloads
type: "reference"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/download_config.go
---

Downloading files at build or deployment time is one of those common tasks
that Escape tries to cover.

## Escape Plan

Downloads are configured in the Escape Plan under the
[`downloads`](/docs/reference/escape-plan/#downloads) field.


Field | Type | Description
------|------|-------------
|url|`string`|The URL to download from. This field is required. 
|||Example: `https://www.google.com/` 
|dest|`string`|The destination path. 
|overwrite|`bool`|Overwrite the destination path if it already exists. 
|if_not_exists|`[string]`|Only perform this download if none of the paths in this list exist. Supports glob patterns (for example: `"*.zip"`) 
|unpack|`bool`|Should Escape try and unpack the destination path after download? Supported extensions: `.zip`, `.tgz`, `.tar.gz`, `.tar`. 
|platform|`string`|Only perform this download if the platform matches this value. Can be used to do platform dependent builds. 
|arch|`string`|Only perform this download if the architecture matches this string. Can be used to do architecture dependent builds. 
|scopes|`[string]`|A list of scopes (`build`, `deploy`) that defines during which stage(s) this download should be performed. 

