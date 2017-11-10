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

```bash
mkdir workspace
escape plan init --minify --name my-project/my-package
```

This should create an Escape plan in the default location `escape.yml`;
minified for education purposes, looking a little something like this:

```yaml
name: my-project/my-package
version: 0.0.@
```

The `.@` at the end of our version signals to Escape that it should start
auto-versioning from there. This is great, and covered in depth
[here](/docs/versioning/), but it requires an Escape Inventory, which
we'll skip for now. Let's version our package explicitly instead, by changing
`escape.yml` to: 

```yaml
name: my-project/my-package
version: 0.1
```

We can use the [plan preview](/docs/escape_plan_preview) command to make sure 
that our plan compiles and to have a look at what Escape makes of it. (NB. We don't 
have to run this command explicitly for any of our build steps as Escape will 
do it automatically, but it can be handy in some cases)

```bash
escape plan preview
```

That's looking tidy. 

We now have enough to create an empty package, but usually we do want to
actually put something inside it. Let's create a file:

```bash
echo "Hello world" > hello.txt
```

Next we'll need to tell Escape to put these files in our package:

```yaml
name: my-project/my-package
version: 0.1

includes:
- hello.txt
```

We can also use [globbing patterns and add whole
directories](/docs/escape-plan/#includes) using this `includes` field, but for
now we can keep it simple. 

We are ready to create our first package!

```bash
escape run release
```

Which outputs:

```
Release: Releasing my-project/my-package-v0.1
  Build: ✔️ Completed build
  Test: ✔️ Tests passed.
  Destroy: ✔️ Destruction complete
  Deploy: ✔️ Successfully deployed my-project/my-package-v0.1 with deployment name my-project/my-package in the dev environment.
  Smoke tests: ✔️ Smoke tests passed.
  Destroy: ✔️ Destruction complete
  Package: ✔️ Packaged my-project/my-package-v0.1 at /home/user/workspace/.escape/target/hello-v0.1.tgz
  Push: ✔️ Push successful.
Release: ✔️ Successfully released my-project/my-package-v0.1%          
```

Voila, we've built our package and made it available in the Inventory.  There
is a lot going on here as Escape runs through all the different phases, but
hopefully all becomes clear in the next section:

[&lt; Back: Hello World](/docs/quickstart-hello-world/)
[&gt; Next: The Release Cycle](/docs/quickstart-the-release-cycle/)
