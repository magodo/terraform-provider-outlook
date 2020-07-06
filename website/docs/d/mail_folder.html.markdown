---
subcategory: ""
layout: "outlook"
page_title: "Outlook Resource: Data Source: outlook_mail_folder"
description: |-
  Gets information about an existing Mail Folder.
---

# Data Source: outlook_mail_folder

Use this data source to access information about an existing Mail Folder.

## Example Usage

```hcl
data "outlook_mail_folder" "example" {
  name = "existing"
}

output "id" {
  value = data.outlook_mail_folder.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Mail Folder.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mail Folder.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mail Folder.
