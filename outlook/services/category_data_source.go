package services

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
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
	// we do not use filter here since the filter in category list API does not work
	objs, err := req.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var category *msgraph.OutlookCategory
	for _, c := range objs {
		if c.DisplayName == nil {
			continue
		}
		if name == *c.DisplayName {
			if category != nil {
				return diag.Errorf("more than one categories are called %s", name)
			}
			category = &c
		}
	}

	if category.ID == nil || *category.ID == "" {
		return diag.Errorf("empty of nil ID returned for Outlook Category %q", name)
	}

	d.SetId(*category.ID)
	if err := d.Set("color", flattenCategoryColor(colorMap, category.Color)); err != nil {
		return diag.Errorf("setting `color`: %+v", err)
	}

	return nil
}
