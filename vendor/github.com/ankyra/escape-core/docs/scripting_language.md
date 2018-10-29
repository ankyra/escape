---
title: "Scripting Language"
slug: scripting-language 
type: "reference"
toc: true
---

The scripting language is a tiny language that is meant to help with wiring up
inputs and outputs. It is by no means a fully fledged language and it's not its
goal to become one either. More complicated logic should be pushed into the
build scripts, where it can be properly tested. 


# Examples

## Using provider outputs as defaults

If we have this provider:

```
id: aws
version: 0.1.@
provides:
- aws

outputs:
- id: default_zone
  default: us-east
```

We can use its `default_zone` output variable as an input for our unit:

```
id: my-unit
version: 0.1.@

consumes:
- aws

inputs:
- id: zone
  default: $aws.outputs.default_zone
```

## Propagating Dependency Metadata

```
id: my-unit
version: 0.1.@

depends:
- my-dependency as dep

metadata:
  author: $dep.metadata.author
  some_other_key: $dep.metadata.some_other_key
```

## Configuring dependencies

```
id: my-unit
version: 0.1.@

inputs:
- input_variable

depends:
- id: my-dependency
  mapping:
    dependency_input_variable: $this.inputs.input_variable
```

Note that this is only necessary if you want to rename or modify the inputs to
the dependency, as the parent's inputs are mapped to the dependency by default.

## Exposing inputs as outputs

```
id: my-unit
version: 0.1.@

inputs:
- input_variable

outputs:
- id: output_variable
  default: $this.inputs.input_variable
```

This can be a handy pattern for providers, where the user selects or enters
something and the unit provides it to the environment:

```
id: my-provider
version: 0.1.@

provides:
- zone

inputs:
- id: zone
  description: Please select a zone.
  items:
  - europe-west1-b
  - europe-west1-c
  - europe-west1-d

outputs:
- id: zone
  default: $this.inputs.zone
```

## Tracking a Dependency's Version

If we have a dependency:

```
id: my-dependency
version: 0.1.@
```
We can use the same versioning scheme in the parent:

```
id: my-project/my-application
depends:
- my-dependency as dep
version: $dep.version.track_minor_version()
```

Which is short for this beauty:

```
version: $dep.version.split(".")[:2].join(".").concat(".@")
```

Also see: `track_major_version`, `track_minor_version`, `track_patch_version`
in the [Standard Library Reference](../scripting-language-stdlib/)


## Configuring Templates

```
id: my-project/my-application
version: 0.1.@
inputs:
- input_variable

templates:
- file: template.txt.tpl
  mapping:
    title: $this.inputs.input_variable
```

# Syntax

## Literals

At the moment only integer and string literals are supported:

```
123
-123
"string value"
```

## Dictionary lookups

A dictionary lookup is performed using the `.` operator. Lookups in `$`, the
global context, are a special case where the `.` should be skipped:

```
$this.inputs.input_variable
```

## Indexing and slicing

Indexing a list uses familiar syntax:

```
$this.inputs.list_input[0]
```

And slicing is also supported:

```
$this.inputs.list_input[0:2]
$this.inputs.list_input[1:]
$this.inputs.list_input[:-1]
```

## Function calls

A function is usually called on an object:

```
$this.inputs.variable.upper()
$this.inputs.variable.split(",")
```

Functions can also be called directly (although the manner in which this is
done is likely to change), which can be handy when you don't have a reference
to an object:

```
$__upper("my string")
$__split("my:string", ":")
$__timestamp()
```

For a full overview of supported functions see the [Standard Library
Reference](../scripting-language-stdlib/).


# Context

`$this` always refers to the current unit and whether we're compiling the
Escape plan, building or deploying, the following fields will always be
available:

```
$this.id                    # e.g. "project/unit-v1.0.0"
$this.name                  # e.g. "unit"
$this.release               # e.g. "unit-v1.0.0"
$this.versionless_release   # e.g. "project/unit"
$this.version               # e.g. "1.0.0"

$this.branch                # e.g. "master"
$this.revision              # e.g. "abcdefff12312312312123"
$this.repository            # e.g. "github.com/ankyra/escape-core"
```

These variables are generally used in templates; we could tag virtual machines,
docker images or Kubernetes deployments with this information for example.

Besides these variables, we also have access to:

```
$this.description           # e.g. "Escape plan description"
$this.metadata.key          # e.g. "Escape plan metadata value"
```

This is also true for all of our dependencies. Given:

```
id: unit
version: 0.1.@
depends:
- my-dependency-latest as dep
```

We can:

```
$dep.id
$dep.name
$dep.release
```

etc.

## Stateful Context

When we're building or deploying we also have access to the state.

```
$this.inputs        # e.g. {"input_variable": "value"}
$this.outputs       # e.g. {"output_variable": "value"}
$this.project       # e.g. "project"
$this.environment   # e.g. "ci"
$this.deployment    # e.g. "my-deployment"
```

Note that the project variable does not refer to the unit's project, but to the
state's project. This will probably change, because it's confusing.

Again, these variables are also available for our dependencies (unless we're
evaluating variables and `eval_before_dependencies` is set to true).

```
$deps.inputs
$deps.outputs
$deps.project       # same as $this.project
$deps.environment   # same as $this.environment
$deps.deployment    # e.g. "my-project/dependency"
```

# Enabled Fields

The language is only enabled in a few fields of the Escape plan so that common 
tasks can be handled. 

## Version 

The version field can be scripted to enabled custom versioning schemes and 
dependency version tracking. No state information will be available in the 
context when this field is compiled, so (dependency) inputs and outputs can't 
be referenced.

Examples:

Follow a dependency one-to-one:

```
version: $dep.version
```

Follow a minor version:

```
version: $dep.version.track_minor_version()
```

## Metadata

The metadata field can be scripted to propagate metadata from dependencies.  No
state information will be available in the context when this field is compiled,
so (dependency) inputs and outputs can't be referenced.

```
metadata:
  website: $dep.metadata.website
```

## Variables

Most importantly the language can be used in variable definition blocks. Both
the `default` and the `items` field can be scripted. 

```
inputs:
- id: input_variable
  default: $this.version

outputs:
- id: output_variable
  default: $this.inputs.input_variable
```

The `inputs` field is a special case, because it generally won't have access to
the full state, because most input variables are evaluated before its
dependencies are run. This can however be configured using the
`eval_before_dependencies` setting, which gives it access to dependency
outputs:

```
inputs:
- id: input_variable
  default: $dep.outputs.image
  eval_before_dependencies: false
```

## Template Mapping

The values in the `mapping` field of a template are evaluated in a stateful
context, which gives full access to inputs and outputs:

```
templates:
- file: kubespec.yml.tpl
  mapping:
    image: $dep.outputs.docker_image
```

## Dependency Mapping

The values in the `mapping` field of a dependency can be scripted:

```
inputs:
- input_variable

depends:
- id: my-dependency-latest
  mapping:
    image: $this.inputs.input_variable
```


# See also

* [Standard Library Reference](../scripting-language-stdlib/)
