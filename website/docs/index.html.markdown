---
layout: "outlook"
page_title: "Provider: Outlook"
description: |-
  The Outlook provider is used to manage outlook related resources.
---

# Outlook Provider

The Azure Provider can be used to configure Outlook Mail Settings using the Microsoft Graph API's. Documentation regarding the [Data Sources](/docs/configuration/data-sources.html) and [Resources](/docs/configuration/resources.html) supported by the Outlook Provider can be found in the navigation to the left.

Use the navigation to the left to read about the available resources.

## Authenticating to MS Graph

Terraform supports a number of different methods authenticating to MS Graph:

* Authenticating to MS Graph using Authorization Code Flow
* Authenticating to MS Graph using Device Flow

---

~> **NOTE** We do not support non-interactively authentication method currently.

### Authenticating to MS Graph using Authorization Code Flow

The authorization code flow is used for devices which has browser installed.

Set provider configuration as below:

```hcl
provider "outlook" {
  auth_method         = "..." # e.g. auth_code_flow
  client_id           = "..." # e.g. 23bd8cd9-a50b-4839-b522-67b77d5db7da
  client_secret       = "..." # not necessary for public native app
  client_redirect_url = "..." # e.g. http://localhost:3000/
}
```

Then run terraform command, there will automatically launch a web browser to allow user to do the authentication.

### Authenticating to MS Graph using Device Flow

The device flow is used for devices which has no browser installed or has limited input capability.

Set provider configuration as below:

```hcl
provider "outlook" {
  auth_method = "device_flow"
  client_id   = "..." # e.g. 23bd8cd9-a50b-4839-b522-67b77d5db7da
}
```

Currently, terraform doesn't allow provider to print any message directly, we have to output the device login link via terraform log. User needs to enable [terraform debug level](https://www.terraform.io/docs/internals/debugging.html) via setting `TF_LOG` to `DEBUG` or `INFO`.

Then when user runs terraform command and should see following line in logs:

```
...
[INFO] To sign in, use a web browser to open https://microsoft.com/devicelogin and enter the code *** to authenticate (with in 900 sec)
...

```

In this point, user should follow the instruction shown above to use another device to finish the login flow.

### Token Cache File

Once the user finishes the authentication, the provider will write the token (including **refresh token**) into a local file (as defined in `token_cache_path` provider configuration or `OUTLOOK_TOKEN_CACHE_PATH` environment variable), in plain text for now. So user needs to make sure to keep this cache file in secure.

## Example Usage

```hcl
# Configure the Outlook Provider
provider "outlook" {
  # auth_method = "auth_code_flow"
  # client_id = "23bd8cd9-a50b-4839-b522-67b77d5db7da"
  # client_secret = ""
  # client_redirect_url = "http://localhost:3000/"
  # token_cache_path = ".terraform-provider-outlook.json"
}

# Create a mail folder
resource "outlook_mail_folder" "example" {
  name = "Foo"
}

# Create a message rule to move message meet certain condition into the folder created above
resource "outlook_message_rule" "example" {
  name     = "move message from foo@bar.com to Foo"
  sequence = "1"
  enabled  = true
  condition {
    from_addresses = ["foo@bar.com"]
  }
  action {
    copy_to_folder = outlook_mail_folder.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `auth_method` - (Optional) The oauth2 authentication method to use. Possible values are `auth_code_flow` and `device_flow`. This can also be sourced from the `OUTLOOK_AUTH_METHOD` Environment Variable. Defaults to `auth_code_flow`.

* `client_id` - (Optional) The AzureAD registered application's Object ID (i.e. oauth2 client_id). This can also be sourced from the `OUTLOOK_CLIENT_ID` Environment Variable. Defaults to `23bd8cd9-a50b-4839-b522-67b77d5db7da`.

* `client_secret` - (Optional) The AzureAD registered application's secret (i.e. oauth2 client_secret). For native public application, you can leave it unset. This can also be sourced from the `OUTLOOK_CLIENT_SECRET` Environment Variable.

* `client_redirect_url` - (Optional) The AzureAD registered application's redirect URL. This can also be sourced from the `OUTLOOK_CLIENT_REDIRECT_URL` Environment Variable. Defaults to `http://localhost:3000/`.

* `token_cache_path` - (Optional) Token cache file path that the provider will export the token info into this file for reuse. Accordingly, the provider will try to load the token from this file if file exists. This can also be sourced from the `OUTLOOK_TOKEN_CACHE_PATH` Environment Variable. Defaults to `.terraform-provider-outlook.json`.
