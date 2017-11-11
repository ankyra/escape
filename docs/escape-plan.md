---
date: 2017-11-11 00:00:00
title: "The Escape Plan"
slug: escape-plan
type: "docs"
toc: true
---

The Escape Plan describes a package. 

Field | Type | Description
------|------|-------------
|<a name='name'></a>name|`string`|The build name. Format: /[a-zA-Z]+[a-zA-Z0-9-]*/ 
|||
|<a name='version'></a>version|`string`|The version. Either specify the full version or use the '@' symbol to let Escape pick the next version at build time. Format: /[0-9]+(\.[0-9]+)*(\.@)?/ 
|||Examples: 
|||Build version 1.5: version: 1.5 
|||Build the next minor release in the 1.* series: version: 1.@ 
|||Build the next path release in the 1.1.* series: version: 1.1.@ 
|||
|<a name='description'></a>description|`string`|An (optional) description for this release. 
|||
|<a name='logo'></a>logo|`string`|
|||
|<a name='extends'></a>extends|`[string]` [Extensions](/docs/extensions/) |
|||
|<a name='depends'></a>depends|`[string]` [Dependencies](/docs/dependencies/) |Dependencies. Reference dependencies by their full release ID or use the '@' symbol to resolve versions at build time. 
|||Examples: 
|||Reference the full release ID to pin to a particular version: depends: [archive-example-v0.1] 
|||To always get the latest version of a particular release: depends: [archive-example-latest] 
|||Or: depends: [archive-example-@] 
|||To resolve the latest minor release: depends: [archive-example-v0.@] 
|||To resolve the latest path release: depends: [archive-example-v0.1.@] 
|||
|<a name='consumes'></a>consumes|`[string]` [Consumers](/docs/providers-and-consumers/) |The release can consume zero or more providers from the environment its deployed in. 
|||
|<a name='build_consumes'></a>build_consumes|`[string]` [Consumers](/docs/providers-and-consumers/) |The release can consume zero or more providers at build time. 
|||
|<a name='deploy_consumes'></a>deploy_consumes|`[string]` [Consumers](/docs/providers-and-consumers/) |The release can consume zero or more providers at deploy time. 
|||
|<a name='provides'></a>provides|`[string]` [Consumers](/docs/providers-and-consumers/) |The release can provide zero or more providers for other releases to consume at deployment time. 
|||
|<a name='inputs'></a>inputs|`[string]` [Variables](/docs/input-and-output-variables/) |Input variables. 
|||Examples: 
|||inputs: - string_input 
|||- id: full_string description: "A nice description" friendly: "Friendly variable display name" type: string 
|||- id: int type: int 
|||- id: choice_string type: string default: first items: - first - second 
|||Escape script can be used in the default field to reference values from other dependencies: 
|||inputs: - id: example default: $dependency.outputs.output_variable 
|||Supported types are "string" (default), "int", "list" and the special types: "version", "project", "environment" and "deployment" which are automatically populated by Escape. 
|||
|<a name='build_inputs'></a>build_inputs|`[string]` [Variables](/docs/input-and-output-variables/) |
|||
|<a name='deploy_inputs'></a>deploy_inputs|`[string]` [Variables](/docs/input-and-output-variables/) |
|||
|<a name='outputs'></a>outputs|`[string]` [Variables](/docs/input-and-output-variables/) |Output variables (see input variables for documentation) 
|||
|<a name='metadata'></a>metadata|`{}` |Metadata key value pairs. 
|||Escape script can be used as values, but note that the metadata is compiled at build time, so dependency inputs and outputs can't be referenced. 
|||Example: 
|||metadata: author: Fictional Character co_author: $dependency.metadata.author 
|||
|<a name='includes'></a>includes|`[]string` |The files to includes in this release. The files don't have to exist at build time. Globbing patterns are supported. 
|||
|<a name='errands'></a>errands|[Errands](/docs/errands/) |Errands. 
|||Errands are scripts that can be run against the deployment of this release. The scripts receive the deployment's inputs and outputs as environment variables. 
|||Examples: 
|||errands: my-errand: description: "Run this errand to do something special" script: bin/my_errand.sh inputs: - extra_input 
|||For information on the syntax of the input variables see the "inputs" field. 
|||
|<a name='downloads'></a>downloads|[Downloads](/docs/downloads/) |
|||
|<a name='templates'></a>templates|[Templates](/docs/templates/) |
|||
|<a name='build_templates'></a>build_templates|[Templates](/docs/templates/) |
|||
|<a name='deploy_templates'></a>deploy_templates|[Templates](/docs/templates/) |
|||
|<a name='path'></a>path|`string`|
|||
|<a name='pre_build'></a>pre_build|`string`|A script to run before the build. 
|||The script has access to the input variables (prepended with INPUT_) in the environment. 
|||Examples: 
|||Given the escape plan: 
|||inputs: - test_input pre_build: pre_build.sh 
|||We can get the value of the input variable in pre_build.sh: 
|||echo $INPUT_test_input 
|||
|<a name='build'></a>build|`string`|
|||
|<a name='post_build'></a>post_build|`string`|A script to run after the build. 
|||The script has access to the input variables (prepended with INPUT_) and output variables (prepended with OUTPUT_) in the environment. 
|||Examples: 
|||Given the escape plan: 
|||inputs: - test_input outputs: - test_output post_build: post_build.sh 
|||We can get the value of the variables in post_build.sh: 
|||echo $INPUT_test_input $OUTPUT_test_output 
|||
|<a name='test'></a>test|`string`|
|||
|<a name='pre_deploy'></a>pre_deploy|`string`|
|||
|<a name='deploy'></a>deploy|`string`|
|||
|<a name='post_deploy'></a>post_deploy|`string`|
|||
|<a name='smoke'></a>smoke|`string`|
|||
|<a name='pre_destroy'></a>pre_destroy|`string`|
|||
|<a name='destroy'></a>destroy|`string`|
|||
|<a name='post_destroy'></a>post_destroy|`string`|
|||
