---
title: "Configuring Escape"
slug: quickstart-configure-escape
type: "docs"
toc: true

back: /docs/inventory-installation/
backLabel: Inventory Installation
next: /docs/quickstart-building-a-package/
nextLabel: Building a Package
contributeLink: https://example.com/
---

We have our own instance of the Inventory running now, but we need to tell
Escape to use it, because our default configuration profile is set up to go the
central Ankyra repository:

```
escape config profile
```

We can "login" to our local instance and create and activate a new profile 
using the [escape login](/docs/escape_login/) command:

```
escape login --url http://localhost:7770/ --target-profile local
```

If everything is working correctly we should now be able to query our empty repository 
from the command line using:

```
escape inventory query --json
```

Beautiful, let's build something!
