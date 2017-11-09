---
title: "What is Escape?"
slug: what-is-escape 
type: "docs"
toc: true
---

Escape is a cross-cloud release engineering kit that aims to reach a few goals:

## Goal 1: Managing Multiple Environments

Whether it's source code. cloud infrastructure, data or documentation; most
artefacts go through one or several phases before going live.  Often this means
that there will be at least one <i>environment</i> that is as close to the real
thing as possible so that the business can develop and test new features, but
on big projects there may be dozens; some only short-lived. 

Managing, promoting and changing configuration between these different
environments can often be quite difficult.  Escape's goal is to make
managing multiple environments easier and automatable.

## Goal 2: Different Layers, Different Tools, Same Process

Different artefacts are often built and configured using different tools. For
example: we may be building a Docker image using the `docker` command and
deploy the image using the `kubectl` command on a Kubernetes cluster that was
built on AWS using the `terraform` command.

Having different tools targetting different layers makes it hard to orchestrate
the deployment of a complete environment and to verify its consistency.
Escape's goal is to make it easy to compose packages, potentially containing
different kind of artefacts, into entire platforms.

## Goal 3: Enable Best Practices in Modern Release Engineering

A modern software application stacks consists of many different layers that all
need to be versioned, configured, deployed and operated differently. The
integration points between these layers are often bespoke, hard to change, and
full of technical debt.

Escape's goal is to bring and somewhat standardise the best practices in
release engineering, but without being overly opinionated. 

## Goal 4: To Make All Packages Identifiable and Environments Self-Documenting

It can be hard to find out what version is live, what version has passed
integration tests, what change was deployed, who authored it, etc. 

Escape should enable you to make everything:

* Identifiable
* Reproducible
* Consistent

[&gt; Next: Get Started](/docs/escape-installation/)
