---
subcategory: ""
layout: "outlook"
page_title: "Outlook Resource: outlook_category"
description: |-
  Manages a Category.
---

# outlook_category

Manages a Category.

## Example Usage

```hcl
resource "outlook_category" "example" {
  name  = "Foo"
  color = "Red"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Category. Changing this forces a new Category to be created.

---

* `color` - (Optional) The color of this Category, possible values are `None`, `Red`, `Orange`, `Brown`, `Yellow`, `Green`, `Teal`, `Olive`, `Blue`, `Purple`, `Cranberry`, `Steel`, `DarkSteel`, `Gray`, `DarkGray`, `Black`, `DarkRed`, `DarkOrange`, `DarkBrown`, `DarkYellow`, `DarkGreen`, `DarkTeal`, `DarkOlive`, `DarkBlue`, `DarkPurple`, `DarkCranberry`. Defaults to `None`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Category.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Category.
* `read` - (Defaults to 5 minutes) Used when retrieving the Category.
* `update` - (Defaults to 30 minutes) Used when updating the Category.
* `delete` - (Defaults to 30 minutes) Used when deleting the Category.

## Import

Categories can be imported using the `resource id`, e.g.

```shell
terraform import outlook_category.example <id>
```
