---
subcategory: ""
layout: "outlook"
page_title: "Outlook Resource: Data Source: outlook_category"
description: |-
  Gets information about an existing Category.
---

# Data Source: outlook_category

Use this data source to access information about an existing Category.

## Example Usage

```hcl
data "outlook_category" "example" {
  name = "existing"
}

output "id" {
  value = data.outlook_category.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Category.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Category.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Category.
