---
subcategory: ""
layout: "outlook"
page_title: "Outlook Resource: outlook_mail_folder"
description: |-
  Manages a Mail Folder.
---

# outlook_mail_folder

Manages a Mail Folder.

~> **NOTE**: Deleting a Mail Folder will not deleting the containing messages, instead those messages will be moved back to inbox. 

## Example Usage

```hcl
resource "outlook_mail_folder" "example" {
  name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Mail Folder.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mail Folder.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Mail Folder.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mail Folder.
* `update` - (Defaults to 30 minutes) Used when updating the Mail Folder.
* `delete` - (Defaults to 30 minutes) Used when deleting the Mail Folder.

## Import

Mail Folders can be imported using the `resource id`, e.g.

```shell
terraform import outlook_mail_folder.example <id>
```
