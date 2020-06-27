package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/utils"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func ResourceMailFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMailCreate,
		ReadContext:   resourceMailRead,
		UpdateContext: resourceMailUpdate,
		DeleteContext: resourceMailDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMailCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders
	name := d.Get("name").(string)

	if d.IsNewResource() {
		req := client.Request()
		req.Filter(fmt.Sprintf(`displayName eq '%s'`, name))
		objs, err := req.Get(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(objs) != 0 {
			return utils.ImportAsExistsError("outlook_mail_folder", *(objs[0].ID))
		}
	}

	param := &msgraph.MailFolder{
		DisplayName: utils.String(name),
	}
	if pid, ok := d.GetOk("parent_folder_id"); ok {
		param.ParentFolderID = utils.String(pid.(string))
	}

	resp, err := client.Request().Add(ctx, param)
	if err != nil {
		return diag.Errorf("creating Mail Folder %q: %w", name, err)
	}

	if resp.ID == nil {
		return diag.Errorf("nil ID for Mail Folder %q", name)
	}
	d.SetId(*resp.ID)

	return resourceMailRead(ctx, d, meta)
}

func resourceMailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders

	resp, err := client.ID(d.Id()).Request().Get(ctx)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			log.Printf("[WARN] Mail Folder %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", resp.DisplayName)
	return nil
}

func resourceMailUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders.ID(d.Id())

	var param msgraph.MailFolder
	if d.HasChange("name") {
		param.DisplayName = utils.String(d.Get("name").(string))
	}
	if err := client.Request().Update(ctx, &param); err != nil {
		return diag.FromErr(err)
	}

	return resourceMailRead(ctx, d, meta)
}

func resourceMailDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders
	if err := client.ID(d.Id()).Request().Delete(ctx); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
