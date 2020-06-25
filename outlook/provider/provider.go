package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token_cache": {
				Type:        schema.TypeString,
				Description: "Token cache path used to read the oauth2 token",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN_CACHE_PATH", ""),
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
		return nil, nil
	}
}
