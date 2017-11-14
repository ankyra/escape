---
title: "Installation"
slug: escape-installation 
type: "docs"
toc: true
---

There are a few ways to install the Escape command line tool onto your own
machine. The preferred way for now is to download one of our pre-built binaries
or use the official Docker image, but you can also build it from source. 

# Pre-built binaries

<div class='docling'>
Note: The following instructions assume 64 bit architectures, which is probably what
you have. You can find 32 bit builds on the <a href='/downloads/'>Downloads</a> page.
</div>

## Linux

```bash
curl -O https://storage.googleapis.com/escape-releases-eu/escape-client/{{version}}/escape-v{{ version }}-linux-amd64.tgz
tar -xvzf escape-v{{version}}-linux-amd64.tgz
sudo mv escape /usr/bin/escape
```

## MacOS

```bash
curl -O https://storage.googleapis.com/escape-releases-eu/escape-client/{{version}}/escape-v{{ version }}-linux-darwin.tgz
tar -xvzf escape-v{{version}}-linux-amd64.tgz
sudo mv escape /usr/bin/escape
```

# Docker images

Ankyra publishes images for Escape into the central Docker hub. 

```bash
docker run -it ankyra/escape:v{{version}} 
```

# From Source

Escape is written in Go and its code is hosted on Github. 

## Using the Go toolchain

To build Escape from source you'll need a functioning Go toolchain, which is
outside the scope of this document. 

```bash
go get -u github.com/ankyra/escape
```

## From Source Using Escape

If you already have an Escape binary (and you love recursion) then you can also
use Escape to build Escape.

```bash
escape run build
escape run test
./escape
```

[&lt; Back: Installation](/docs/what-is-escape/)
[&gt; Next: Hello World](/docs/quickstart-hello-world/)
