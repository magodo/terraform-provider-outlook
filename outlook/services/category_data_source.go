package services

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
)

func DataSourceOutlookCategory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOutlookCategoryRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"color": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutlookCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	name := d.Get("name").(string)

	req := client.Request()
	req.Filter(fmt.Sprintf("displayName eq '%s'", name))
	objs, err := req.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(objs) != 1 {
		return diag.Errorf("expect one category with display name %q but got %q", name, len(objs))
	}

	obj := objs[0]
	if obj.ID == nil || *obj.ID == "" {
		return diag.Errorf("empty of nil ID returned for Outlook Category %q", name)
	}

	d.SetId(*obj.ID)
	if err := d.Set("color", flattenCategoryColor(colorMap, obj.Color)); err != nil {
		return diag.Errorf("setting `color`: %+v", err)
	}

	return nil
}
