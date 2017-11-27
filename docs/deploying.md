---
title: "Deploying Environments"
slug: quickstart-deploying
type: "docs"
toc: true
wip: true

back: /docs/quickstart-input-variables/
backLabel: Input Variables
next: /docs/quickstart-errands/
nextLabel: Errands
---

In the previous section we've built a package that can greet a configurable
subject:

```yaml
name: quickstart/hello-world
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: hello_world.sh
deploy: hello_world.sh

inputs:
- id: who
  default: World
  type: string
  description: Who should we be greeting?
```

With the following script:

```bash
#!/bin/bash -e

echo "Hello ${INPUT_who}!"
```

We have released this package into our Inventory using `escape run release`,
but now we'd like to deploy this package into another environment. 

Suppose we have three environments: Test, Stage and Live. Each environment is
more volatile than the next. "Test" will have constant changes applied to it,
and once approved things might make it to "Stage" and then "Live". It is
likely, or at least possible, that we'll have different versions running in
each of our environments.

<img src='envs.png'>

There are a few ways we can control these deployments.

## Deploying by version

Suppose our Test environment is looking good and we've received word that
`quickstart/hello-world-v0.0.1` is good to go to Stage.

If we know exactly what version we want to deploy into our environment we can
ask Escape to deploy it using the [escape run deploy](/docs/escape_run_deploy)
command:

```
escape run deploy quickstart/hello-world-v0.0.1 -v who=Stage -e stage
```

Which should output something like this:

```
Deploy: Running deploy step /home/user/workspace/deps/quickstart/hello-world/hello_world.sh.
Deploy: hello_world.sh: Hello Stage!
Deploy: ✔️ Successfully deployed hello-world-v0.0.1 with deployment name quickstart/hello-world in the stage environment.
```

Escape keeps track of what's deployed where and what inputs were used in a
_state file_. The state file contains information for every deployment grouped by
environment. We can deploy a package multiple times in the same environment, but 
this does mean we have to give deployments themselves a unique name. 

In the usual case we only want to deploy things once so Escape generates a
default deployment name based on the package's project and name (e.g.
`quickstart/hello-world` for our `quickstart/hello-world-v0.0.1` package), but
we can override this behaviour by using the `-d / --deployment` flag to pass in 
another deployment name. 

## Promoting from another environment

People are happy with the deployment we've put on Stage and they'd like to see
this same version deployed into the Live environment. We could use the same
approach as before, but why should we when Escape knows everything about what
version is where by reading its state file?

Let's use the [escape run promote](/docs/escape_run_promote) command instead
and have Escape deploy whatever is in deployment name `quickstart/hello_world`
from Stage to Live:

```bash
escape run promote --deployment "quickstart/hello-world" --environment stage --to live
```

Which should output something like:

```
Promote: Deployment quickstart/hello-world in environment stage has quickstart/hello-world-v0.0.1.
Promote: Deployment quickstart/hello-world in environment live is not present.
Promote quickstart/hello-world-v0.0.1 from stage (quickstart/hello-world) to live (quickstart/hello-world)? [Yn]:
Promote: Promoting quickstart/hello-world-v0.0.1 from stage to live.
  Deploy: Running deploy step /home/user/workspace/deps/quickstart/hello-world/hello_world.sh.
  Deploy: hello_world.sh: Hello World!
  Deploy: ✔️ Successfully deployed hello-world-v0.0.1 with deployment name quickstart/hello-world in the live environment.
```

<div class='docling'>
To figure out what deployments are where we can examine the state using the
<a href='/docs/escape_state'>escape state</a> command.
</div>

## Deploying using dynamic versions

In some cases we always want to deploy the latest version:

```
escape run deploy quickstart/hello-world-latest -v who=World -e stage
```

In some cases we want to deploy the latest patch version:

```
escape run deploy quickstart/hello-world-v0.0.@ -v who=World -e stage
```


## Converge

TODO, sorry

```bash
escape state create
escape run converge
```
