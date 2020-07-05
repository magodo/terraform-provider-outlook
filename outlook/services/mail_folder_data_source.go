package services

import (
	"fmt"
	"time"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
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
	req := client.Request()
	req.Filter(fmt.Sprintf(`displayName eq '%s'`, name))
	objs, err := req.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(objs) != 1 {
		return diag.Errorf("expect one mail folder but got %d", len(objs))
	}
	obj := objs[0]
	if obj.ID == nil || *obj.ID == "" {
		return diag.Errorf("empty or nil ID returned for Mail Folder %q ID", name)
	}

	d.SetId(*obj.ID)

	return nil
}
