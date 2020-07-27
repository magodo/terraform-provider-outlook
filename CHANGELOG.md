## 0.0.4

FEATURES:

* Data Source `outlook_mail_folder` supports `well_known_name`.

## 0.0.3 

FEATURES:

* New Resource/Data Source: `outlook_category`

BUG FIXES:

* Add announcement about MS Graph Outlook API of up to 4 concurrent requests allowed.
* `outlook_mail_folder`: remove the ability to move mail out of folder in parallel

## 0.0.2

FEATURES:

* `outlook_mail_folder`: deleting a mail folder will move the containing messages to inbox, then delete the folder.
* Provider: allow user to specify oauth2 client and auth method underused.

## 0.0.1

FEATURES:

* New Resource/Data Source: `outlook_mail_folder`
* New Resource: `outlook_message_rule`
