🚧 This is very much WIP, do not use in production. 🚧

# Important Disclaimer

I am a Microsoft employee, but this is not an official Microsoft product nor an endorsed product. Purely a project for fun.

# Terraform Provider for Outlook

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

## Using the Provider 

Run following command to install the provider into your `$GOPATH/bin`.

```
$ go get -u github.com/magodo/terraform-provider-outlook
```

Then follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

## Provider Documents

Currently the documents for this provider is not hosted by the official site [Terraform Providers](https://www.terraform.io/docs/providers/index.html). Hence you have to follow the instructions [here](https://github.com/hashicorp/terraform-website#new-provider-repositories) to manually setup a bit so that you can run `make website` to see the document.

In the near future, we will host the provider on terraform registry, and the document will be hosted there also.
