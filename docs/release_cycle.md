---
title: "The Release Cycle"
slug: quickstart-release-cycle
type: "docs"
toc: true

back: /docs/quickstart-building-a-package/
backLabel: Building a Package
next: /docs/quickstart-inputs-and-outputs/
nextLabel: Inputs and Outputs
---

Our Hello World package from the previous section contains a single file, 
but other than that it doesn't do much:

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

## Build and test

The life cycle of almost every unit starts with a build step. This step is
generally only run on developer's workstations and CI servers to make sure all
the artefacts that are necessary at deployment time are part of the package.

We can trigger the build step of an Escape plan using the [escape run
build](/docs/escape_run_build) command:

```bash
escape run build
```

Our almost empty Escape plan doesn't define a build step however, so this step
completes very quickly, but we can add one:

```yaml
name: quickstart/hello-world
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: hello_world.sh
deploy:
```

Now Escape expects this file to exist when we try and run the build again:

```bash
$ escape run build
Compile: Error: File 'hello_world.sh' was referenced in the escape plan, but it doesn't exist
```

Let's do Escape a favour and create the `hello_world.sh` file:

```bash
#!/bin/bash -e

echo "Hello World!"
```

And try that again:

```bash
$ escape run build
Build: Running build step /home/bspaans/src/workspace/escape/hello_world.sh.
Build: hello_world.sh: Hello World!
Build: ✔️ Completed build
```

After a build step we'd usually want to run a test step to make sure that what we've 
built the right thing:

```
escape run test
```

Again, we haven't defined this field in our Escape plan, so this step completes
quickly, but we can add a `test` field to make Escape run another script.
However, our Hello World package currently doesn't contain anything worth
testing so we'll skip this.

<div class='docling'>
You can use the `pre_build` and `post_build` fields to run scripts after and
before the main build script. This can be handy later when dealing with
Extensions. 
</div>


## Deploy and smoke

We have a build phase that can be used to generate files, binaries, docs, etc.
and then test them. Next, we want to be able to _deploy_ these files into
separate _environments_. 

We can trigger the deploy step of an Escape plan using the [escape run
deploy](/docs/escape_run_deploy) command:

```
$ escape run deploy
```

We haven't specified a deployment script however. We can re-use our "Hello
World" script from before:

```
name: quickstart/hello-world
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: hello_world.sh
deploy: hello_world.sh
```

```bash
$ escape run deploy
Deploy: Running deploy step /home/bspaans/src/workspace/escape/hello_world.sh.
Deploy: hello_world.sh: Hello World!
Deploy: ✔️ Successfully deployed hello-world-v0.0.1 with deployment name quickstart/hello-world in the dev environment.
```

## Release


Running the release step will run the build, test, deploy, smoke, destroy,
package and push steps in succession to make sure the unit is working
end-to-end (although parts can be skipped see [`escape run
release`](/docs/escape_run_release))

```bash
$ escape run release
Release: Releasing quickstart/hello-world-v0.0.1
  Build: Running build step /home/bspaans/src/workspace/escape/hello_world.sh.
  Build: hello_world.sh: Hello World!
  Build: ✔️ Completed build
  Test: ✔️ Tests passed.
  Destroy: ✔️ Destruction complete
  Deploy: Running deploy step /home/bspaans/src/workspace/escape/hello_world.sh.
  Deploy: hello_world.sh: Hello World!
  Deploy: ✔️ Successfully deployed hello-world-v0.0.1 with deployment name quickstart/hello-world in the dev environment.
  Smoke tests: ✔️ Smoke tests passed.
  Destroy: ✔️ Destruction complete
  Package: ✔️ Packaged quickstart/hello-world-v0.0.1 at /home/bspaans/src/workspace/escape/.escape/target/hello-world-v0.0.1.tgz
  Push: ✔️ Push successful.
Release: ✔️ Successfully released quickstart/hello-world-v0.0.1
```

As we can see our Hello World script gets run twice. So why does the `release`
command run both the build and deploy step? Our general philosophy is that once
a package is in the inventory all its steps should have been tested. The main
thing that a package will be used for after publishing is deployments so it's
good to test this early.
