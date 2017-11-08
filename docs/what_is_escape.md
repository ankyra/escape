---
title: "What is Escape?"
slug: what-is-escape 
type: "docs"
toc: true
---

Escape is a cross-cloud release engineering kit, ....

It’s trying to solve these problems:
Promoting artifacts from check-in to live
Expensive ad-hoc pipelines and processes that are hard to change
Having different deployment processes for different layers of the infrastructure
Knowing what’s running where. What version is live, what version is in CI, etc.
Complexity of configuration management
Reusability of code and components
Close coupling between configuration and code

Escape makes it possible to ensure that
Source code
Cloud infrastructure
Configuration
Open Source Components
Data
Documentation

Is:
Identifiable
Reproducible
Consistent
Automated
Operatable

A good, cohesive CD pipeline can capture a big part of a software business’ processes.

How does Escape work?
Escape can build, test and deploy packages targeting any layer of the stack by wrapping around other tools like terraform, packer, kubectl, helm, etc.
Packages are centrally stored in the Escape Inventory
Packages can depend on and extend other packages; making different layers of the stack composable.
Packages can be deployed into environments

