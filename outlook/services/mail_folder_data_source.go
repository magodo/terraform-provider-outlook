package services

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	msgraph "github.com/yaegashi/msgraph.go/beta"
)

func DataSourceMailFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMailRead,

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

func dataSourceMailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MailFolders
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	name := d.Get("name").(string)
	resp, err := client.ID(name).Request().Get(ctx)
	if err != nil {
		return err
	}
	if err != nil {
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return fmt.Errorf("Mail Folder %q not found", name)
			}
		}
		return fmt.Errorf("retrieving Mail Folder %q: %w", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Mail Folder %q ID", name)
	}

	d.SetId(*resp.ID)

	return nil
}
