---
title: "The Release Cycle"
slug: quickstart-release-cycle
type: "docs"
toc: true
---

## Build and test

The life cycle of almost every unit starts with a build step. This step is
generally only run on developer's workstations and CI servers. 

```
$ escape build
```

Our almost empty Escape plan doesn't define a build step however, 
but we can add one:

```
name: my-project/my-deployment-unit
version: 1.0
build: build_my_thing.sh
```

Next we might wanna test our build:

```
$ escape test
```

And again we can define a field:

```
name: my-project/my-deployment-unit
version: 1.0
build: build_my_thing.sh
test: test_my_thing.sh
```

## Deploy and smoke

Once a unit has been built we want to be able to deploy it into an
_environment_.

```
$ escape deploy
```

Same thing again. And let's also add a field for `escape smoke` to run our
smoke tests.


```
name: my-project/my-deployment-unit
version: 1.0
build: build_my_thing.sh
test: test_my_thing.sh
deploy: deploy_my_thing.sh
smoke: smoke_test_my_thing.sh
```


## Package

So far, so good. We're just running some tasks. Nothing fancy, but let's
package it all up into a distributable unit:

```
$ escape package
```

This takes all the files referenced in the Escape plan and adds them into an
archive, combined with the compiled metadata. We can include more files by glob
patterns:

```
name: my-project/my-deployment-unit
version: 1.0
build: build_my_thing.sh
test: test_my_thing.sh
deploy: deploy_my_thing.sh
smoke: smoke_test_my_thing.sh

includes:
- src/*.src
- assets/
- README

```

And then package again (we need to use the `-f` flag, because we already
created an archive for this version in the previous command)

```
$ escape package -f
```

## Push

Once happy we can push it to a server (note: the client points to the Ankyra
registry by default, which doesn't allow public writes, but you can run a
stand-alone registry).

```
$ escape push
```

This will make `my-project/my-deployment-unit-v1.0` available in the registry
and it will effectively freeze our version, because we can't upload the same
version twice. Our fellow engineers can now deploy this release by running:

```
$ escape deploy my-project/my-deployment-unit-v1.0
```

## Release

Instead of running all the steps one by one we can also use the `escape
release` command, which is generally preferable in CI settings.

```
$ escape release
```

This will run the build, test, deploy, smoke, destroy, package and push steps
in succession to make sure the unit is working end-to-end (although parts can
be skipped see `escape release --help`)

