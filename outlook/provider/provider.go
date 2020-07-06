package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/msauth"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/services"
	"github.com/pkg/browser"
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
			"token_cache_path": {
				Type:        schema.TypeString,
				Description: "Token Cache Path",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN_CACHE_PATH", ".terraform-provider-outlook.json"),
			},
		},

		DataSourcesMap: SupportedDataSources(),
		ResourcesMap:   SupportedResources(),
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		tokenCachePath := d.Get("token_cache_path").(string)
		const (
			// Use msgraph tutorial client id as client id. As custom registered app
			// at tenant AAD level is not able to invoke outlook ms graph API.
			// (seems only "first-party" app at "common" auth endpoint can work)
			clientID = "6731de76-14a6-49ae-97bc-6eba6914391e"
			tenantID = "common"
		)
		app := msauth.NewApp()
		if err := app.ImportCache(tokenCachePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		ts, err := app.ObtainTokenSourceViaDeviceFlow(context.Background(), tenantID, clientID,
			func(auth msauth.DeviceAuthorization) {
				browser.OpenReader(strings.NewReader(fmt.Sprintf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate (with in %d sec).\n",
					auth.VerificationURI, auth.UserCode, auth.ExpiresIn)))
			},
			"mailboxsettings.readwrite",
			"mail.readwrite",
			"offline_access",
		)
		if err != nil {
			return nil, err
		}
		if err := app.ExportCache(tokenCachePath); err != nil {
			return nil, err
		}
		return clients.NewClient(msgraph.NewClient(oauth2.NewClient(context.Background(), ts)).BaseRequestBuilder), nil
	}
}
