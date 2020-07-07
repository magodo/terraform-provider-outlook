package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Description: "Token cache file path. If specified, the provider will export the token info into this file for reuse. Accordingly, the provider will try to load the token from this file if file exists.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTLOOK_TOKEN_CACHE_PATH", ""),
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
			clientID = "6731de76-14a6-49ae-97bc-6eba6914391e"
			tenantID = "common"
		)
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
		ts, err := app.ObtainTokenSourceViaDeviceFlow(context.Background(), tenantID, clientID,
			func(auth msauth.DeviceAuthorizationAuth) error {
				return browser.OpenReader(
					strings.NewReader(
						buildDeviceflowMessage(auth.VerificationURI, auth.UserCode, auth.ExpiresIn)),
				)
			},
			"mailboxsettings.readwrite",
			"mail.readwrite",
			"offline_access",
		)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return clients.NewClient(msgraph.NewClient(oauth2.NewClient(context.Background(), ts)).BaseRequestBuilder), nil
	}
}

func buildDeviceflowMessage(uri, code string, timeout int) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body>

<h1>Terraform Outlook Provider</h1>

<p>To sign in, use a web browser to open the <a href="%s">Microsoft device login page</a> and enter the code <p>%s</p> to authenticate (with in %d sec).</p>

</body>
</html>
`, uri, code, timeout)
}
