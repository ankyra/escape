---
title: "Templates"
slug: templates 
type: "docs"
toc: true
---

```
name: my-project/my-deployment-unit
version: 0.0.@
inputs:
- my_variable

templates:
- file: output.txt.tpl
```

```
name: my-project/my-deployment-unit
version: 0.0.@
inputs:
- my_variable

templates:
- file: output.txt.tpl
  mapping:
    input: $this.inputs.my_variable
    input2: "yo"
```

