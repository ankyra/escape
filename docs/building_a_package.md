---
title: "Building a Package"
slug: quickstart-building-a-package
type: "docs"
toc: true
---

At its core Escape provides abstractions to work with _packages_. A _package_
is a collection of files, plus a bit of metadata to tell Escape what it can do.
Based on the metadata Escape knows how to build, test, deploy, destroy and
operate the unit.

Metadata doesn't have to be written by hand, but can be compiled from an
[Escape plan](/docs/escape-plan/).  Let's create a new workspace and initialise
a new Escape plan using the [plan init](/docs/escape_plan_init) command.

```
mkdir workspace
escape plan init -n my-project/my-package
```

This should create an Escape plan in the default location `escape.yml`.
