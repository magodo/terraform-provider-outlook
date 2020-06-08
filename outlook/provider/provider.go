package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/provider/client"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Description: "MS Graph access token with proper scope defined. Can be specified with the `OUTLOOK_TOKEN`",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"outlook_user": dataUser(),
		},
		ResourcesMap: map[string]*schema.Resource{},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		token := d.Get("token").(string)
		c := client.BuildClient(token)
		return c, nil
	}
}
