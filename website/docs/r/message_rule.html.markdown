---
subcategory: ""
layout: "outlook"
page_title: "Outlook Resource: outlook_message_rule"
description: |-
  Manages a Message Rule.
---

# outlook_message_rule

Manages a Message Rule.

## Example Usage (Single Rule)

```hcl
resource "outlook_mail_folder" "example" {
  name = "Foo"
}

resource "outlook_message_rule" "example" {
  name    = "move message from foo@bar.com to Foo"
  enabled = true
  condition {
    from_addresses = ["foo@bar.com"]
  }
  action {
    move_to_folder = outlook_mail_folder.example.id
  }
}
```

## Example Usage (Multiple Rules)

```hcl
resource "outlook_mail_folder" "example" {
  name = "Foo"
}

resource "outlook_message_rule" "move" {
  name     = "move message from foo@bar.com to Foo"
  sequence = 1
  enabled  = true
  condition {
    from_addresses = ["foo@bar.com"]
  }
  action {
    move_to_folder = outlook_mail_folder.example.id
  }
}

resource "outlook_message_rule" "mark" {
  name     = "flag"
  sequence = outlook_message_rule.move.sequence + 1 # Explicitly refer to the "move" rule so as to ensure the creation order.
  enabled  = true
  condition {
    from_addresses = ["foo@bar.com"]
  }
  action {
    mark_importance = "low"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Message Rule. Changing this forces a new Message Rule to be created.

* `action` - (Required) A `action` block as defined below.

---

* `condition` - (Optional) A `condition` block as defined below. The messages meet the condition will be processed.

* `enabled` - (Optional) Should the Message Rule be enabled?

* `sequence` - (Optional) Indicates the order in which the rule is executed, among other rules (the lower number executes first). User should specify a **unique** and **sequential number from starts 1** to each rule. By default, it equals to the amount of existing rules plus one (i.e. append to the end).

~> **NOTE**: Even if `sequence` is specified, it has no effect on creation. Outlook API will reset the sequence number based on the creation order in FIFO. User can rerun `terraform apply` until `terraform plan` doesn't give any differences. A best practice is to explicitly control the creation order via using the reference between resources, as illustrated in example above.

* `exception` - (Optional) Same as `condition`, except the messages meet the condition will not be processed.

---

A `action` block supports the following:

* `copy_to_folder` - (Optional) The ID of a folder that a message is to be copied to.

* `delete` - (Optional) Indicates whether a message should be moved to the Deleted Items folder.

* `forward_as_attachment_to` - (Optional) Specifies a list of the email addresses of the recipients to which a message should be forwarded as an attachment.

* `forward_to` - (Optional) Specifies a list of the email addresses of the recipients to which a message should be forwarded.

* `mark_as_read` - (Optional) Indicates whether a message should be marked as read.

* `mark_importance` - (Optional) Sets the importance of the message, possible values are `low`, `normal`, `high`.

* `move_to_folder` - (Optional) The ID of the folder that a message will be moved to.

* `permanent_delete` - (Optional) Indicates whether a message should be permanently deleted (rather than saved to the Deleted Items folder).

* `redirect_to` - (Optional) Specifies a list of the email addresses to which a message should be redirected.

* `stop_processing_rules` - (Optional) Indicates whether subsequent rules should be evaluated.

---

A `condition` block supports the following:

* `body_contains` - (Optional) Specifies the strings that should appear in the body of an incoming message in order for the condition or exception to apply.

* `body_or_subject_contains` - (Optional) Specifies a list of the strings that should appear in the body or subject of an incoming message in order for the condition or exception to apply.

* `from_addresses` - (Optional) Specifies a list of the specific sender email addresses of an incoming message in order for the condition or exception to apply.

* `has_attachments` - (Optional) Whether an incoming message must have attachments in order for the condition or exception to apply.

* `header_contains` - (Optional) Specifies a list of the strings that appear in the headers of an incoming message in order for the condition or exception to apply.

* `importance` - (Optional) Whether the importance that is stamped on an incoming message in order for the condition or exception to apply? Possible values are `low`, `normal`, `high`.

* `is_approval_request` - (Optional) Whether an incoming message must be an approval request in order for the condition or exception to apply?

* `is_automatic_forward` - (Optional) Whether an incoming message must be automatically forwarded in order for the condition or exception to apply.

* `is_automatic_reply` - (Optional) Whether an incoming message must be an auto reply in order for the condition or exception to apply.

* `is_encrypted` - (Optional) Whether an incoming message must be encrypted in order for the condition or exception to apply.

* `is_meeting_request` - (Optional) Whether an incoming message must be a meeting request in order for the condition or exception to apply.

* `is_meeting_response` - (Optional) Whether an incoming message must be a meeting response in order for the condition or exception to apply.

* `is_non_delivery_report` - (Optional) Whether an incoming message must be a non-delivery report in order for the condition or exception to apply.

* `is_permission_controlled` - (Optional) Whether an incoming message must be permission controlled (RMS-protected) in order for the condition or exception to apply.

* `is_read_receipt` - (Optional) Whether an incoming message must be a read receipt in order for the condition or exception to apply.

* `is_signed` - (Optional) Whether an incoming message must be S/MIME-signed in order for the condition or exception to apply.

* `is_voicemail` - (Optional) Whether an incoming message must be a voice mail in order for the condition or exception to apply.

* `message_action_flag` - (Optional) Specifies the flag-for-action value that appears on an incoming message in order for the condition or exception to apply. Possible values are `any`, `call`,`doNotForward`, `followUp`, `fyi`, `forward`, `noResponseNecessary`, `read`, `reply`, `replyToAll`, `review`.

* `not_sent_to_me` - (Optional) Whether the owner of the mailbox must not be a recipient of an incoming message in order for the condition or exception to apply.

* `recipient_contains` - (Optional) Specifies a list of the strings that appear in either the **toRecipients** or **ccRecipients** properties of an incoming message in order for the condition or exception to apply.

* `sender_contains` - (Optional) Specifies a list of the strings that appear in the **from** property of an incoming message in order for the condition or exception to apply.

* `sensitivity` - (Optional) Specifies the sensitivity level that must be stamped on an incoming message in order for the condition or exception to apply. Possible values are `normal`, `personal`, `private`, `confidential`.

* `sent_cc_me` - (Optional) Whether the owner of the mailbox must be in the **ccRecipients** property of an incoming message in order for the condition or exception to apply.

* `sent_only_to_me` - (Optional) Whether the owner of the mailbox must be the only recipient in an incoming message in order for the condition or exception to apply.

* `sent_to_addresses` - (Optional) Specifies a list of the email addresses that an incoming message must have been sent to in order for the condition or exception to apply.

* `sent_to_me` - (Optional) Whether the owner of the mailbox must be in the **toRecipients** property of an incoming message in order for the condition or exception to apply.

* `sent_to_or_cc_me` - (Optional) Whether the owner of the mailbox must be in either a **toRecipients** or **ccRecipients** property of an incoming message in order for the condition or exception to apply.

* `subject_contains` - (Optional) Specifies a list of the strings that appear in the subject of an incoming message in order for the condition or exception to apply.

* `within_size_range` - (Optional) A `within_size_range` block as defined below.

---

A `within_size_range` block supports the following:

* `max_size` - (Required) Specifies the maximum size (in kilobytes) that an incoming message must have in order for a condition or exception to apply.

* `min_size` - (Required) Specifies the minimum size (in kilobytes) that an incoming message must have in order for a condition or exception to apply.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Message Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Message Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Message Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Message Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Message Rule.

## Import

Message Rules can be imported using the `resource id`, e.g.

```shell
terraform import outlook_message_rule.example <id>
```
