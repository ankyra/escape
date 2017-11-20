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

Suppose we have three environments: Testing, Staging and Production. Each
Environment is more volatile than the next. Testing will have constant changes
and once approved things might make it to Staging and then Production. It is
likely, or at least possible, that we'll have different versions running in
each of our environments.

```bash
escape run deploy quickstart/hello-world-latest -v who=You -e testing
escape run deploy quickstart/hello-world-latest -v who=World -e staging
escape run deploy quickstart/hello-world-latest -v who=Universe -e production
```
