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
| **Block List** | `field { ... }` (repeated per entry) | `field.0.subfield` |
| **Block Set** | `field { ... }` (repeated per entry) | Not indexable, see [Referencing Sets](#referencing-sets) |
| **Nested Attribute** | `field = { ... }` | `field.subfield` |
| **Nested Attribute List** | `field = [{ ... }]` | `field[0].subfield` |
| **Nested Attribute Set** | `field = [{ ... }]` | Not indexable, see [Referencing Sets](#referencing-sets) |
| **Read-Only Object** | Not user-declarable | `field.subfield` |
| **Read-Only Object List** | Not user-declarable | `field.0.subfield` |
| **Read-Only Object Set** | Not user-declarable | Not indexable, see [Referencing Sets](#referencing-sets) |

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

List nested attributes are referenced by index, like any other list:

```hcl
linode_interface.example.vpc.ipv4.addresses[0].address
```

### Referencing Sets

Sets are unordered, so their elements have no addressable keys. Any attempt to index into a
set fails at plan time:

```hcl
# Invalid: `domain_grant` is a Block Set
output "grant_id" {
  value = linode_user.example.domain_grant.0.id
}
```

```
Error: Cannot index a set value

Block type "domain_grant" is represented by a set of objects, and set elements do not have
addressable keys. To find elements matching specific criteria, use a "for" expression with
an "if" clause.
```

Instead, convert the set to a list, or project the values you need:

```hcl
# Reference every element
output "grant_ids" {
  value = [for grant in linode_user.example.domain_grant : grant.id]
}

# Select a specific element
output "example_grant_id" {
  value = one([for grant in linode_user.example.domain_grant : grant.id if grant.label == "example.com"])
}

# Convert the set to a list when an index is unavoidable
output "first_grant_id" {
  value = tolist(linode_user.example.domain_grant)[0].id
}
```

-> `tolist(...)` gives no ordering guarantee, so only use it when any element will do.
