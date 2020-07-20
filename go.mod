module github.com/magodo/terraform-provider-outlook

go 1.14

require (
	github.com/bflad/tfproviderlint v0.14.0
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.0-rc.2
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/sergi/go-diff v1.0.0
	github.com/yaegashi/msgraph.go v0.1.2
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/tools v0.0.0-20200529172331-a64b76657301
)

replace github.com/yaegashi/msgraph.go => ../msgraph.go
