---
title: "Input Variables"
slug: quickstart-input-variables
type: "docs"
toc: true
wip: true

back: /docs/quickstart-building-a-package/
backLabel: Building a Package
next: /docs/quickstart-deploying/
nextLabel: Deploying Environments
---

Our Hello World package from the previous section is a beautiful piece of
engineering:


```yaml
name: quickstart/hello-world
version: 0.0.@
description: 
logo: 

includes:
- README.md

build: hello_world.sh
deploy: hello_world.sh
```

When we deploy our package it greets the world, but our Product Manager 
wants our package to greet the Universe when we deploy to production.

This kind of thing happens a lot, even in less contrived examples; we often
have configuration that changes depending on what environment it's being
deployed in.

In our Escape plan we can make these configuration variables explicit and define 
them as inputs:

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
- who
```

When we try and run a build or deployment step Escape will now start complaining:

```bash
$ escape run build   
Build: Starting build.
Build: Error: Missing value for variable 'who'
```

Defining input variables makes Escape expect a value for both builds and
deployments. We can scope variables to specific phases, add type checks, a
description, defaults, [and all sorts of
stuff](/docs/input-and-output-variables/), but we'll keep ours simple for now:

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

Because we've set a default value Escape will use this when we don't specify
one ourselves:

```bash
$ escape run build
Build: Running build step /home/bspaans/src/workspace/escape/hello_world.sh.
Build: hello_world.sh: Hello World!
Build: ✔️ Completed build
```

The next step is to update our `hello_world.sh` script to use the new variable.
Variables are passed into build scripts by their `id` (in this case "who" -- we
know it should be "whom", don't email in) and receive the `INPUT_` prefix. This
is a complicated way to say that we should update our script to:

```bash
#!/bin/bash -e

echo "Hello ${INPUT_who}!"
```

And just like that we've made our build and deployment configurable:

```bash
$ escape run build -v who=You
Build: Running build step /home/bspaans/src/workspace/escape/hello_world.sh.
Build: hello_world.sh: Hello You!
Build: ✔️ Completed build
```

Inputs are stored in the Escape state, which means we don't have to set the
input again on the next run:

```bash
$ escape run build           
Build: Running build step /home/bspaans/src/workspace/escape/hello_world.sh.
Build: hello_world.sh: Hello You!
Build: ✔️ Completed build
```

We will find out more about the state in the next session as we try to deploy
our package to multiple environments, but let's release our progress first:

```
escape run release
```

