# Scripting Language

The scripting language is a tiny language meant to make it easier to wire up
inputs and outputs. It is by no means a fully fledged language and it's not its
goal to become one either. More complicated logic should be pushed into the
build scripts, where it can be properly tested. 

## Examples

### Follow Dependency Version

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

## Syntax

### Simple values

```
123
"string value"
```

### Map lookups

```
$this.inputs.input_variable
```

### Indexing and slicing

```
$this.inputs.list_input[0]
$this.inputs.list_input[0:2]
$this.inputs.list_input[1:]
$this.inputs.list_input[:-1]
```

### Function calls

```
$this.inputs.variable.split()
```


## Context

`$this` always refers to the current unit.

### Dependency context


## Enabled Fields

The language is only enabled in a few fields of the Escape plan so that common 
tasks can be handled. 

### Version 

### Metadata

### Variables

Most importantly the language can be used in variable definition blocks. 


## Built-in Functions

```

```
