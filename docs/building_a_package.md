---
title: "Building a Package"
slug: quickstart-building-a-package
type: "docs"
toc: true

back: /docs/quickstart-configure-inventory/
backLabel: Configuring an Inventory
next: /docs/quickstart-release-cycle/
nextLabel: The Release Cycle
contributeLink: https://example.com/
---

At its core Escape provides abstractions to work with _packages_. A _package_
is a collection of files, plus a bit of metadata to tell Escape what it can do.
Based on the metadata Escape knows how to build, test, deploy, destroy and
operate the unit.

Metadata doesn't have to be written by hand, but can be compiled from an
[Escape plan](/docs/escape-plan/).  Let's create a new workspace and initialise
a new Escape plan using the [plan init](/docs/escape_plan_init/) command.

```bash
mkdir workspace
cd workspace
escape plan init --name quickstart/hello-world
```

This should create an Escape plan in the default location `escape.yml`; looking
a little something like this:

```yaml
name: quickstart/hello-world
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: 
deploy:
```

<div class='docling'>
The `.@` at the end of our version signals to Escape that it should
auto-version from there. Versioning is covered in depth
[here](/docs/versioning/), 
</div>

We can use the [plan preview](/docs/escape_plan_preview/) command to make sure 
that our plan compiles and to have a look at what Escape makes of it:


```bash
escape plan preview
```

That's looking tidy. We don't have to run the `escape plan preview` command
explicitly for any of our build steps as Escape will do it automatically, but
it can be a handy validation step). NB. If you're getting an "Unauthorized"
error here, you need to configure your
[Inventory](/docs/quickstart-configure-inventory/).

We now have enough to create an empty package, but usually we do want to put
something inside it. Our Escape plan includes a reference to `README.md`, but
the file doesnae exist! Let's create it:

```bash
echo "Thanks for reading" > README.md
```

We've only told Escape to include `README.md`, but we can also add [globbing
patterns and whole directories](/docs/escape-plan/#includes).  In anyway, we
are ready to create our first package!

```bash
escape run release
```

Which outputs:

```
Release: Releasing quickstart/hello-world-v0.0.0
  Build: ✔️ Completed build
  Test: ✔️ Tests passed.
  Destroy: ✔️ Destruction complete
  Deploy: ✔️ Successfully deployed hello-world-v0.0.0 with deployment name quickstart/hello-world in the dev environment.
  Smoke tests: ✔️ Smoke tests passed.
  Destroy: ✔️ Destruction complete
  Package: ✔️ Packaged quickstart/hello-world-v0.0.0 at /home/user/workspace/.escape/target/hello-world-v0.0.0.tgz
  Push: ✔️ Push successful.
Release: ✔️ Successfully released quickstart/hello-world-v0.0.0
```

We can keep running `escape run release` and see the version number increase
for each successful push.

We've built our package and made it available in the Inventory.  There
is a lot going on here as Escape runs through all the different phases, but
hopefully all becomes clear in the next section:

