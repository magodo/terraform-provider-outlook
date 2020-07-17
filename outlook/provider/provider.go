package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/magodo/terraform-provider-outlook/msauth"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/services"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

const (
	AUTH_METHOD_AUTH_CODE_FLOW = "auth_code_flow"
	AUTH_METHOD_DEVICE_FLOW    = "device_flow"
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
			"auth_method": {
				Type:        schema.TypeString,
				Description: "The oauth2 authentication method to use.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_AUTH_METHOD", AUTH_METHOD_AUTH_CODE_FLOW),
				ValidateFunc: validation.StringInSlice([]string{
					AUTH_METHOD_AUTH_CODE_FLOW,
					AUTH_METHOD_DEVICE_FLOW,
				}, false),
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "The AzureAD registered application's Object ID (i.e. oauth2 client_id)",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_CLIENT_ID", "23bd8cd9-a50b-4839-b522-67b77d5db7da"),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "The AzureAD registered application's secret (i.e. oauth2 client_secret). For native public application, you can leave it unset.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_CLIENT_SECRET", ""),
			},
			"client_redirect_url": {
				Type:        schema.TypeString,
				Description: "The AzureAD registered application's redirect URL",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_CLIENT_REDIRECT_URL", "http://localhost:3000/"),
			},
			"token_cache_path": {
				Type:        schema.TypeString,
				Description: "Token cache file path that the provider will export the token info into this file for reuse. Accordingly, the provider will try to load the token from this file if file exists.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN_CACHE_PATH", ".terraform-provider-outlook.json"),
			},
			"feature": featureSchema,
		},

		DataSourcesMap: SupportedDataSources(),
		ResourcesMap:   SupportedResources(),
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (meta interface{}, diags diag.Diagnostics) {
		var (
			clientID     = d.Get("client_id").(string)
			clientSecret = d.Get("client_secret").(string)
			redirectURL  = d.Get("client_redirect_url").(string)
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
		err := app.ImportCache(tokenCachePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				defer func() {
					if !diags.HasError() {
						if err := app.ExportCache(tokenCachePath); err != nil {
							diags = diag.FromErr(err)
						}
					}
				}()
			} else {
				return nil, diag.FromErr(err)
			}
		}

		var ts oauth2.TokenSource

		switch d.Get("auth_method").(string) {

		case AUTH_METHOD_AUTH_CODE_FLOW:
			ts, err = app.ObtainTokenSourceViaAuthorizationCodeFlow(context.Background(), tenantID, clientID, clientSecret, redirectURL, scopes...)

		case AUTH_METHOD_DEVICE_FLOW:
			ts, err = app.ObtainTokenSourceViaDeviceFlow(context.Background(), tenantID, clientID,
				func(auth msauth.DeviceAuthorizationAuth) error {
					// Currently there is no way for a provider to print messsage to user's console unless using debug message.
					log.Printf("[INFO] To sign in, use a web browser to open %s and enter the code %s to authenticate (with in %d sec)", auth.VerificationURI, auth.UserCode, auth.ExpiresIn)
					return nil
				},
				scopes...,
			)
		default:
			return nil, diag.FromErr(fmt.Errorf("Unknown auth method: %s", d.Get("auth_method").(string)))
		}

		if err != nil {
			return nil, diag.FromErr(err)
		}

		feature := expandFeature(d.Get("feature").([]interface{}))
		return clients.NewClient(msgraph.NewClient(oauth2.NewClient(context.Background(), ts)).BaseRequestBuilder, feature), nil
	}
}
