---
page_title: "Linode: Blocks vs. Nested Attributes"
description: |-
  Explains the difference between block and nested attribute syntax in the Linode provider.
---

# Blocks vs. Nested Attributes

Fields with a nested structure are implemented in this provider as either a **block** or a
**nested attribute**. The two implementations are declared and referenced using different
syntax, so the documentation for every nested field is tagged with the implementation it uses:

| Tag | Declaration | Reference |
| --- | --- | --- |
| **Block** | `field { ... }` | `field.0.subfield` |
| **Block List** / **Block Set** | `field { ... }` (repeated per entry) | `field.0.subfield` |
| **Nested Attribute** | `field = { ... }` | `field.subfield` |
| **Nested Attribute List** / **Nested Attribute Set** | `field = [{ ... }]` | `field[0].subfield` |
| **Read-Only Object** | Not user-declarable | `field.subfield` |
| **Read-Only Object List** | Not user-declarable | `field.0.subfield` |

Read-only tags are used for computed fields that can only be referenced, never declared.
A **Read-Only Object List** is often a single logical object that is exposed as a list
because of the same legacy limitation that affects blocks, so it must be referenced through
the `0` index.

All new nested fields are implemented as nested attributes. Blocks are only used by
existing fields, and are not used for new functionality because they are effectively
unsupported by HashiCorp and carry a number of bugs and limitations.

## Declaration Syntax

A **block** is declared without an equals sign, and each entry is a separate block:

```hcl
resource "linode_instance" "example" {
  # ...

  metadata {
    user_data = base64encode("...")
  }
}
```

A **nested attribute** is declared with an equals sign, using an object (or a list of
objects for list and set nested attributes):

```hcl
resource "linode_interface" "example" {
  # ...

  public = {
    ipv4 = {
      addresses = [
        {
          address = "auto"
          primary = true
        }
      ]
    }
  }
}
```

## Reference Syntax

Blocks that hold a single logical object still have to be implemented as a list of blocks,
so they must be referenced through the `0` index:

```hcl
# `metadata` is a Block
linode_instance.example.metadata.0.user_data
```

Nested attributes do not have this limitation, so a single nested attribute is referenced
directly:

```hcl
# `public` is a Nested Attribute
linode_interface.example.public.ipv4
```

List and set nested attributes are referenced by index, like any other list or set:

```hcl
linode_interface.example.public.ipv4.addresses[0].address
```
