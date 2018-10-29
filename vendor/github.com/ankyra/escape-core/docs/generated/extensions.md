---
date: 2017-11-11 00:00:00
title: "Extensions"
slug: extensions
type: "reference"
toc: true
wip: false
contributeLink: https://github.com/ankyra/escape-core/blob/master/extension_config.go
---

A package can extend another package and inherit its build scripts and input
and output variables. This makes it possible to reuse build and deployment
patterns. To extend another package the `extends` field in the Escape plan can
be used:

```
extends:
- my-project/my-application-latest
```

This effectively copies the following fields from the extension to the package:

* `depends`
* `consumes`
* `provides`
* `inputs`
* `outputs`
* `errands`
* `metadata`
* `templates`
* `deploy`, `build`, `test`, `pre_deploy`, etc.

If both the parent release and the extension release define a particular input,
output, or build/deployment script then the definition in the parent always
wins.  This allows you to override parts of the extension that need
specialisation or changing.

## How do scripts get executed?

When executing a script defined in the extension, Escape will run the script in the
parent's root directory. For example: if we have an extension:

```
name: my-extension
version: 1.0

deploy: deploy.sh
```

And a package that uses that extension:

```
name: my-extender
version: 1.0
extends:
- my-extension-latest
```

Then `escape run deploy` for that package will execute `deploy.sh`, but in
my-extender's root directory.

## Example extensions

At Ankyra we use extensions heavily in our delivery pipeline and we've open
sourced a few of them. These extensions might not do exactly what you want, but
hopefully give a flavour of what can be done:

* [extension-golang-binary](https://github.com/ankyra/extension-golang-binary): Builds Go binaries in a Docker image.
* [extension-docker](https://github.com/ankyra/extension-docker): Does Docker image builds and pushes.
* [extension-kubespec](https://github.com/ankyra/extension-kubespec): Does Kubernetes deployments.

## Shortcomings

Extensions are a great way to capture build and deployment patterns, but there are
some things to note:

* Because you can override fields in a parent release, you can also
	_accidentally_ override extension fields in a parent release, and Escape
	won't warn you. This can make debugging harder.
* It can be tricky to release the extension itself, because as soon as you
	define a testing step that step will get inherited by packages that extend
	it; which is usually not what you want.

## Escape Plan

Extensions are configured in the [`extends`](/docs/reference/escape-plan/#extends)
field of the Escape Plan.


Field | Type | Description
------|------|-------------
|release_id|`string`|The release id is required and is resolved at *build* time and then persisted in the release metadata ensuring that deployments always use the same versions. 

