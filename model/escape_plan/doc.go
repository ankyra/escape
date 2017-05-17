/*
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
*/

package escape_plan

var docMap = map[string]string{
	"name":        nameDoc,
	"version":     versionDoc,
	"description": descriptionDoc,
	"depends":     dependsDoc,
	"includes":    includesDoc,
	"consumes":    consumesDoc,
	"provides":    providesDoc,
	"metadata":    metadataDoc,
	"inputs":      inputsDoc,
	"outputs":     outputsDoc,
	"errands":     errandsDoc,
	"pre_build":   prebuildDoc,
	"post_build":  postbuildDoc,
}

func GetDoc(key string) []byte {
	doc, ok := docMap[key]
	if !ok {
		return []byte{}
	}
	return []byte(doc)
}

const nameDoc = `# The build name. Format: /[a-zA-Z]+[a-zA-Z0-9-]*/
#`

const versionDoc = `# The version. Either specify the full version or use the '@' symbol to 
# let Escape pick the next version at build time. Format: /[0-9]+(\.[0-9]+)*(\.@)?/
# 
# Examples:
#
# Build version 1.5:
#   version: 1.5
#
# Build the next minor release in the 1.* series:
#   version: 1.@
#
# Build the next path release in the 1.1.* series:
#   version: 1.1.@
#`

const descriptionDoc = `# An (optional) description for this release.
#`

const dependsDoc = `# Dependencies. Reference dependencies by their full release ID or use
# the '@' symbol to resolve versions at build time.
#
# Examples:
#
# Reference the full release ID to pin to a particular version:
#   depends: [archive-example-v0.1]
#
# To always get the latest version of a particular release:
#   depends: [archive-example-latest]
#
# Or:
#   depends: [archive-example-@]
#
# To resolve the latest minor release:
#   depends: [archive-example-v0.@]
#
# To resolve the latest path release:
#   depends: [archive-example-v0.1.@]
#`

const includesDoc = `# The files to includes in this release. The files don't have to exist at 
# build time. Globbing patterns are supported.
#`

const consumesDoc = `# The release can consume zero or more providers from the environment its
# deployed in.
#`

const providesDoc = `# The release can provide zero or more providers for other releases to 
# consume at deployment time.
#`

const metadataDoc = `# Metadata key value pairs.
#
# Escape script can be used as values, but note that the metadata is compiled
# at build time, so dependency inputs and outputs can't be referenced.
# 
# Example:
#
#   metadata:
#     author: Fictional Character
#     co_author: $dependency.metadata.author
#`

const inputsDoc = `# Input variables. 
#
# Examples:
#
#   inputs:
#   - string_input
#
#   - id: full_string
#     description: "A nice description"
#     friendly: "Friendly variable display name"
#     type: string
#
#   - id: int
#     type: int
#
#   - id: choice_string
#     type: string
#     default: first
#     items:
#     - first
#     - second
#
# Escape script can be used in the default field to reference values from other 
# dependencies:
#
#   inputs:
#   - id: example
#     default: $dependency.outputs.output_variable
#
# Supported types are "string" (default), "int", "list" and the special types:
# "version", "project", "environment" and "deployment" which are automatically 
# populated by Escape.
#`

const outputsDoc = `# Output variables (see input variables for documentation)
#`

const errandsDoc = `# Errands.
#
# Errands are scripts that can be run against the deployment of this release.
# The scripts receive the deployment's inputs and outputs as environment
# variables.
#
# Examples:
#
#   errands:
#     my-errand:
#       description: "Run this errand to do something special"
#       script: bin/my_errand.sh
#       inputs:
#       - extra_input
#
# For information on the syntax of the input variables see the "inputs" field.
#`

const prebuildDoc = `# A script to run before the build.
#
# The script has access to the input variables (prepended with INPUT_)
# in the environment.
#
# Examples:
#
# Given the escape plan:
#
#   inputs:
#    - test_input
#   pre_build: pre_build.sh
#
# We can get the value of the input variable in pre_build.sh:
#
#   echo $INPUT_test_input
#`

const postbuildDoc = `# A script to run after the build.
#
# The script has access to the input variables (prepended with INPUT_) 
# and output variables (prepended with OUTPUT_) in the environment.
#
# Examples:
#
# Given the escape plan:
#
#   inputs:
#    - test_input
#   outputs:
#   - test_output
#   post_build: post_build.sh
#
# We can get the value of the variables in post_build.sh:
#
#   echo $INPUT_test_input $OUTPUT_test_output
#`
