🚧 This is very much WIP, do not use in production. 🚧

# Important Disclaimer

I am a Microsoft employee, but this is not an official Microsoft product nor an endorsed product. Purely a project for fun.

# Terraform Provider for Outlook

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

## Using the Provider 

### Terraform 0.13 and above

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/magodo/outlook).

### Terraform 0.12 or manual installation

You can download a pre-built binary from the [releases](https://github.com/magodo/terraform-provider-outlook/releases) page, these are built using [goreleaser](https://goreleaser.com/) (the [configuration](.goreleaser.yml) is in the repo). You can verify the signature and my [key ownership via Keybase](https://keybase.io/magodo).

If you want to build from source, you can simply use `go build` in the root of the repository.

To use an external provider binary with Terraform ([until the provider registry is GA](https://www.hashicorp.com/blog/announcing-providers-in-the-new-terraform-registry/)), you need to place the binary in a [plugin location for Terraform](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) to find it.

## Provider Documents

The document of this provider is available on [Terraform Provider Registry](https://registry.terraform.io/providers/magodo/outlook/latest/docs).
