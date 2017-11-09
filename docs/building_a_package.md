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
escape plan init --minify --name my-project/my-package
```

This should create an Escape plan in the default location `escape.yml`;
minified for education purposes, looking a little something like this:

```
name: my-project/my-package
version: 0.0.@
```

The `.@` at the end of our version signals to Escape that it should start
auto-versioning from there. This is great, and covered in depth
[here](/docs/escape_versioning), but it requires an Escape Inventory, which
we'll skip for now. Let's version our package explicitly instead, by changing
`escape.yml` to: 

```
name: my-project/my-package
version: 0.1
```

We can use the [plan preview](/docs/escape_plan_preview) command to make sure 
that our plan compiles and to have a look at what Escape makes of it. (NB. We don't 
have to run this command explicitly for any of our build steps as Escape will 
do it automatically, but it can be handy in some cases)

```
escape plan preview
```

That's looking tidy. 
