# Escape 

[![Build status](https://circleci.com/gh/ankyra/escape.svg?style=shield&circle-token=2f7f6d97a01eefe7b3d1967ce11bb183034c963d)](https://circleci.com/gh/ankyra/escape) 

![Escape Logo](/header.png)

Escape is a tool that can help with building, testing, versioning, deploying,
composing and operating software platforms. Its goal is to provide best
practices in release engineering to make it easier to perform these tasks
across environments and layers. 

Some things you can do with Escape:

* Create repeatable builds with predictable versioning, packaging and distribution 
* Use the same delivery process for different tools and layers: deploy your
  infrastructure and container code like you deploy your application code.
* Manage multiple environments: promote from Dev to Prod.
* Composition: break your platform up into logical components and compose them
  into a cohesive platform.
* Simplify configuration management
* Operate running deployments
* Create self-documenting environments and releases

This repository contains the official Escape client. The Escape Inventory can
be found [here](https://github.com/ankyra/escape-inventory).

## Downloads

Cross-platform binaries can be downloaded from [the
webite](https://escape.ankyra.io/downloads/) where you'll also find the
[installation
instructions](https://escape.ankyra.io/docs/escape-installation/).

## Docker

You can also use the Escape docker image which is published in the [central
Docker hub](https://hub.docker.com/r/ankyra/escape/).

`docker run -it ankyra/escape:latest`

## Usage and Documentation

See the [Escape Docs](https://escape.ankyra.io/docs/) for the full documentation.

## Support and Contact

Issues and general enquiries can be raised on the Github issue tracker. 
You can also join our [Community Slack
channel](https://join.slack.com/t/ankyra-escape/shared_invite/enQtMzI4NDU4NDUwMDk2LTYwNjQ5Nzc1ZThlYTEyMjJkMTYzMDMxNzkxYzg0ZTE3ZjNlNWM2MmExNWFlYzU1NTQ2MTM2NjVlMGI0NjhhMmY)
for realtime interaction.

## License

```
Copyright 2017, 2018 Ankyra

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
