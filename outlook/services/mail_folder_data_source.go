package services

import (
	"net/http"
	"time"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	msgraph "github.com/yaegashi/msgraph.go/beta"
)

func DataSourceMailFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMailRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceMailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders

	name := d.Get("name").(string)
	resp, err := client.ID(name).Request().Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return diag.Errorf("Mail Folder %q not found", name)
			}
		}
		return diag.Errorf("retrieving Mail Folder %q: %w", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return diag.Errorf("empty or nil ID returned for Mail Folder %q ID", name)
	}

	d.SetId(*resp.ID)

	return nil
}
