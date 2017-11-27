---
title: "Dependencies"
slug: quickstart-dependencies
type: "docs"
toc: true
wip: true

back: /docs/quickstart-deploying/
backLabel: Deploying Environments
next: /docs/quickstart-output-variables/
nextLabel: Output Variables
---

More often than not a package will depend on a another package to function
correctly. For example: an application could need a database, some files on 
the filesystem, a Docker image, a configured virtual machine, etc. 

Dependencies can be expressed in Escape using the
[depends](/docs/escape-plan/#depends) field of the Escape plan. Dependencies
are just regular old Escape packages, so this mechanism provides us one way to
componentize our infrastructure and software estate (we will look at more in
later sections)

There are many useful examples, but to keep things simple we're going to do
something very contrived and build upon our previous work.  So let's create a
new Escape plan to build a package that depends on our previous
`quickstart/hello-world` package. 

```yaml
name: quickstart/introduction
version: 0.0.@
description: 
logo: 

depends:
- quickstart/hello-world-latest

```

## Configuring Dependencies

Using literals:

```yaml
name: quickstart/introduction
version: 0.0.@
description: 
logo: 

depends:
- release_id: quickstart/hello-world-latest
  mapping:
    who: Everyone I Know

```

Passing variables from the parent: 

```yaml
name: quickstart/introduction
version: 0.0.@
description: 
logo: 

depends:
- release_id: quickstart/hello-world-latest
  mapping:
    who: $this.inputs.who

inputs:
- who
```

Which happens by default, so is equivalent to:

```yaml
name: quickstart/introduction
version: 0.0.@
description: 
logo: 

depends:
- quickstart/hello-world-latest

inputs:
- who
```

We might change this default behaviour, because it can get a bit confusing, so
we recommend to keep these mappings explicit.
