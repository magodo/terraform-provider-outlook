package services

import (
	"fmt"
	"time"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func DataSourceMailFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMailRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "well_known_name"},
			},
			"parent_folder_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"well_known_name"},
			},
			"well_known_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validation.StringInSlice([]string{
					"archive",
					"clutter",
					"conflicts",
					"conversationhistory",
					"deleteditems",
					"drafts",
					"inbox",
					"junkemail",
					"localfailures",
					"msgfolderroot",
					"outbox",
					"recoverableitemsdeletions",
					"scheduled",
					"searchfolders",
					"sentitems",
					"severfailures",
					"syncissues",
				},
					false),
				ExactlyOneOf:  []string{"name", "well_known_name"},
				ConflictsWith: []string{"parent_folder_id"},
			},
		},
	}
}

func dataSourceMailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders

	name := d.Get("name").(string)
	wellKnownName := d.Get("well_known_name").(string)
	parent := d.Get("parent_folder_id").(string)

	var obj *msgraph.MailFolder
	if wellKnownName != "" {
		var err error
		obj, err = client.ID(wellKnownName).Request().Get(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		var (
			objs []msgraph.MailFolder
			err  error
		)
		if parent == "" {
			req := client.Request()
			req.Filter(fmt.Sprintf(`displayName eq '%s'`, name))
			objs, err = req.Get(ctx)
		} else {
			req := client.ID(parent).ChildFolders().Request()
			req.Filter(fmt.Sprintf(`displayName eq '%s'`, name))
			objs, err = req.Get(ctx)
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if len(objs) != 1 {
			return diag.Errorf("expect one mail folder but got %d", len(objs))
		}
		obj = &objs[0]
	}
	if obj.ID == nil || *obj.ID == "" {
		return diag.Errorf("empty or nil ID returned for Mail Folder %q ID", name)
	}
	d.SetId(*obj.ID)

	return nil
}
