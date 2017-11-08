---
title: "Hello World! 👋"
slug: quickstart-hello-world
type: "docs"
toc: true
---

To make sure our local [installation](/docs/escape-installation/) works we
can try and run the [version](/docs/escape_version/) command to output the
version that Escape is built with.

```bash
escape version
```

If that works then the time has come to have Escape greet the world by 
deploying the latest `hello-world` package, which will be fetched from the 
[Public Inventory](https://escape.ankyra.io/app/registry/_/hello-world/latest/).

```bash
mkdir workspace
cd workspace
escape run deploy hello-world-latest
```

The output may surprise you:

```
$ escape run deploy hello-world-latest  
Deploy: Running deploy step /home/user/workspace/deps/_/hello-world/hello.sh.
Deploy: hello.sh: Hello world! 👋
  Smoke tests: ✔️ Smoke tests passed.
Deploy: ✔️ Successfully deployed.
```

One of Escape's defining features is the ability to codify configuration inputs
and outputs, so that code can be properly separated from configuration.

Let's try and deploy hello-world again, but this time using an input variable:

```bash
escape run deploy hello-world-latest -v who=you
```

We will learn more about this as we try to [build our own packages](/docs/quickstart-building-a-package/)

[&lt; Back: Installation](/docs/escape-installation/)
[&gt; Next: Hello World](/docs/quickstart-building-a-package/)
