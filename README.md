# Escape 

Escape was written to solve common deployment woes. It was originially built to
support "Infrastructure as Code" platform delivery from start to finish, but
nowadays it can also be used to version, package and deploy documentation,
data, applications, or whatever you want really.

Its goal is to replace ad-hoc release engineering, deployment, orchestration
and operational systems with a cohesive, but otherwise unopinionated model.

## Features

* Compose a platform out of versioned _units of deployment_
* Separate configuration from code
* Multi-environment
* Easy to integrate into CI and CD processes
* Common release engineering tasks: optional automatic versioning, version
  control linking, packaging, uploading, downloading, templating
* Operations as code
* Extensions to work with Docker, Packer, Kubernetes, Terraform, ...

## Installation

The easiest way to install Escape is to download the binaries from the website. 
See https://escape.ankyra.io/downloads/

### Build

```
git clone https://github.com/ankyra/escape-client.git "$GOPATH/src/github.com/ankyra/escape-client"
go install
mv "$GOPATH/bin/escape-client" "$GOPATH/bin/escape"
```

## License

```
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
