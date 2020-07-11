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
		CreateContext: resourceMailFolderCreate,
		ReadContext:   resourceMailFolderRead,
		UpdateContext: resourceMailFolderUpdate,
		DeleteContext: resourceMailFolderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"parent_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMailFolderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders
	name := d.Get("name").(string)
	parent := d.Get("parent_folder_id").(string)

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

	var (
		resp *msgraph.MailFolder
		err  error
	)
	if parent == "" {
		resp, err = client.Request().Add(ctx, param)
	} else {
		resp, err = client.ID(parent).ChildFolders().Request().Add(ctx, param)
	}

	if err != nil {
		return diag.Errorf("creating Mail Folder %q: %+v", name, err)
	}

	if resp.ID == nil {
		return diag.Errorf("nil ID for Mail Folder %q", name)
	}
	d.SetId(*resp.ID)

	return resourceMailFolderRead(ctx, d, meta)
}

func resourceMailFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("parent_folder_id", resp.ParentFolderID)
	return nil
}

func resourceMailFolderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders.ID(d.Id())

	var param msgraph.MailFolder
	if d.HasChange("name") {
		param.DisplayName = utils.String(d.Get("name").(string))
	}
	if err := client.Request().Update(ctx, &param); err != nil {
		return diag.FromErr(err)
	}

	return resourceMailFolderRead(ctx, d, meta)
}

func resourceMailFolderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MailFolders

	// Avoid delete the folder when it has child folder
	children, err := client.ID(d.Id()).ChildFolders().Request().Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(children) != 0 {
		return diag.Errorf("deleting a folder with child folder is not allowed")
	}

	const dustbinName = "Deleted Items"
	req := client.Request()
	req.Filter(fmt.Sprintf(`displayName eq '%s'`, dustbinName))
	objs, err := req.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(objs) != 1 {
		return diag.Errorf("expect one mail folder with name '%s' but got %d", dustbinName, len(objs))
	}

	dustbin := objs[0]
	if dustbin.ID == nil || *dustbin.ID == "" {
		return diag.Errorf("empty or nil ID returned for Mail Folder %q ID", dustbinName)
	}

	if _, err := client.ID(d.Id()).Move(&msgraph.MailFolderMoveRequestParameter{DestinationID: dustbin.ID}).Request().Post(ctx); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
