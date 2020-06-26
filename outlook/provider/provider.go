package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/msauth"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/services"
	"github.com/pkg/browser"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"outlook_mail_folder": services.DataSourceMailFolder(),
		},
		ResourcesMap: map[string]*schema.Resource{},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		const (
			// Use msgraph tutorial client id as client id. As custom registered app
			// at tenant AAD level is not able to invoke outlook ms graph API.
			// (seems only "first-party" app at "common" auth endpoint can work)
			clientID = "6731de76-14a6-49ae-97bc-6eba6914391e"
			tenantID = "common"
		)
		ts, err := msauth.NewClientViaDeviceFlow(tenantID, clientID,
			func(auth msauth.DeviceAuthorization) {
				browser.OpenReader(strings.NewReader(fmt.Sprintf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate (with in %d sec).\n",
					auth.VerificationURI, auth.UserCode, auth.ExpiresIn)))
			},
			"mailboxsettings.readwrite",
			"mail.readwrite",
			"offline_access",
		).ObtainTokenSource(context.Background())
		if err != nil {
			return nil, err
		}
		return clients.NewClient(msgraph.NewClient(oauth2.NewClient(context.Background(), ts)).BaseRequestBuilder), nil
	}
}
