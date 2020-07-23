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

* `name` - (Optional) The name of this Mail Folder.

~> **NOTE** Either `name` or `well_known_name` should be specified.

* `parent_folder_id` - (Optional) The ID of the parent folder of the Mail Folder which is specified by `name`.

* `well_known_name` - (Optional) The [well-known Mail Folder name](https://docs.microsoft.com/en-us/graph/api/resources/mailfolder?view=graph-rest-1.0).

~> **NOTE** Either `name` or `well_known_name` should be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mail Folder.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mail Folder.
