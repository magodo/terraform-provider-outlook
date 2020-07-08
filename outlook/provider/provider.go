package provider

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/msauth"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/services"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

func SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"outlook_mail_folder":  services.ResourceMailFolder(),
		"outlook_message_rule": services.ResourceMessageRule(),
	}
}

func SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"outlook_mail_folder": services.DataSourceMailFolder(),
	}
}

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"browser_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the environment running terraform is able to open a browser",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_BROWSER_ENABLED", false),
			},
			"token_cache_path": {
				Type:        schema.TypeString,
				Description: "Token cache file path that the provider will export the token info into this file for reuse. Accordingly, the provider will try to load the token from this file if file exists.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN_CACHE_PATH", ".terraform-provider-outlook.json"),
			},
		},

		DataSourcesMap: SupportedDataSources(),
		ResourcesMap:   SupportedResources(),
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (meta interface{}, diags diag.Diagnostics) {
		const (
			// Use msgraph tutorial client id as client id. As custom registered app
			// at tenant AAD level is not able to invoke outlook ms graph API.
			// (seems only "first-party" app at "common" auth endpoint can work)
			clientID     = "6731de76-14a6-49ae-97bc-6eba6914391e"
			clientSecret = "JqQX2PNo9bpM0uEihUPzyrh"
			redirectURL  = "http://localhost:8888/myapp/"
			tenantID     = "common"
		)
		scopes := []string{
			"mailboxsettings.readwrite",
			"mail.readwrite",
			"offline_access",
		}
		app := msauth.NewApp()

		// Import token cache if specified, accordingly export the updated token cache at the end of configuring provider.
		tokenCachePath := d.Get("token_cache_path").(string)
		if tokenCachePath != "" {
			if err := app.ImportCache(tokenCachePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return nil, diag.FromErr(err)
			}
			defer func() {
				if !diags.HasError() {
					if err := app.ExportCache(tokenCachePath); err != nil {
						diags = diag.FromErr(err)
					}
				}
			}()
		}

		// For device with browser available, use authorization code flow.
		// Otherwise, use device flow.
		var (
			ts  oauth2.TokenSource
			err error
		)

		if d.Get("browser_enabled").(bool) {
			ts, err = app.ObtainTokenSourceViaAuthorizationCodeFlow(context.Background(), tenantID, clientID, clientSecret, redirectURL, scopes...)
		} else {
			ts, err = app.ObtainTokenSourceViaDeviceFlow(context.Background(), tenantID, clientID,
				func(auth msauth.DeviceAuthorizationAuth) error {
					log.Printf("[INFO] To sign in, use a web browser to open %s and enter the code %s to authenticate (with in %d sec)", auth.VerificationURI, auth.UserCode, auth.ExpiresIn)
					return nil
				},
				scopes...,
			)
		}
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return clients.NewClient(msgraph.NewClient(oauth2.NewClient(context.Background(), ts)).BaseRequestBuilder), nil
	}
}
