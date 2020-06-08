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
$ GO111MODULE=off go get -u github.com/magodo/terraform-provider-outlook
```

Then follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

## Provider Documents

Currently the documents for this provider is not hosted by the official site [Terraform Providers](https://www.terraform.io/docs/providers/index.html). Please enter the provider directory and build the website locally.

```sh
$ make website
```

The commands above will start a docker-based web server powered by [Middleman](https://middlemanapp.com/), which hosts the documents in `website` directory. Simply open `http://localhost:4567/docs/providers/outlook` and enjoy them.
