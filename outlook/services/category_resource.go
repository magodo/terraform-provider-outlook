package services

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/utils"
	"github.com/magodo/terraform-provider-outlook/outlook/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func ResourceCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArmCategoryCreate,
		ReadContext:   resourceArmCategoryRead,
		UpdateContext: resourceArmCategoryUpdate,
		DeleteContext: resourceArmCategoryDelete,

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

			"color": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				ValidateDiagFunc: validation.StringInSlice(keySlice(colorMap), false),
			},
		},
	}
}

var (
	// colorMap fetched from https://docs.microsoft.com/en-us/graph/api/resources/outlookcategory?view=graph-rest-1.0
	colorMap = map[string]msgraph.CategoryColor{
		"None":          msgraph.CategoryColorVNone,
		"Red":           msgraph.CategoryColorVPreset0,
		"Orange":        msgraph.CategoryColorVPreset1,
		"Brown":         msgraph.CategoryColorVPreset2,
		"Yellow":        msgraph.CategoryColorVPreset3,
		"Green":         msgraph.CategoryColorVPreset4,
		"Teal":          msgraph.CategoryColorVPreset5,
		"Olive":         msgraph.CategoryColorVPreset6,
		"Blue":          msgraph.CategoryColorVPreset7,
		"Purple":        msgraph.CategoryColorVPreset8,
		"Cranberry":     msgraph.CategoryColorVPreset9,
		"Steel":         msgraph.CategoryColorVPreset10,
		"DarkSteel":     msgraph.CategoryColorVPreset11,
		"Gray":          msgraph.CategoryColorVPreset12,
		"DarkGray":      msgraph.CategoryColorVPreset13,
		"Black":         msgraph.CategoryColorVPreset14,
		"DarkRed":       msgraph.CategoryColorVPreset15,
		"DarkOrange":    msgraph.CategoryColorVPreset16,
		"DarkBrown":     msgraph.CategoryColorVPreset17,
		"DarkYellow":    msgraph.CategoryColorVPreset18,
		"DarkGreen":     msgraph.CategoryColorVPreset19,
		"DarkTeal":      msgraph.CategoryColorVPreset20,
		"DarkOlive":     msgraph.CategoryColorVPreset21,
		"DarkBlue":      msgraph.CategoryColorVPreset22,
		"DarkPurple":    msgraph.CategoryColorVPreset23,
		"DarkCranberry": msgraph.CategoryColorVPreset24,
	}
)

func resourceArmCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	name := d.Get("name").(string)

	if d.IsNewResource() {
		req := client.Request()
		// we do not use filter here since the filter in category list API does not work
		objs, err := req.Get(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		existing := getCategoryByName(objs, name)
		if existing != nil {
			return utils.ImportAsExistsError("outlook_category", *existing.ID)
		}
	}

	param := &msgraph.OutlookCategory{
		DisplayName: utils.String(name),
		Color:       expandCategoryColor(colorMap, d.Get("color").(string)),
	}

	resp, err := client.Request().Add(ctx, param)
	if err != nil {
		return diag.Errorf("creating Outlook Category %q: %+v", name, err)
	}

	if resp.ID == nil {
		return diag.Errorf("nil ID for Outlook Category %q", name)
	}
	d.SetId(*resp.ID)

	return resourceArmCategoryRead(ctx, d, meta)
}

func resourceArmCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	resp, err := client.ID(d.Id()).Request().Get(ctx)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			log.Printf("[WARN] Outlook Category %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", resp.DisplayName)
	if err := d.Set("color", flattenCategoryColor(colorMap, resp.Color)); err != nil {
		return diag.Errorf("setting `color`: %+v", err)
	}

	return nil
}

func resourceArmCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	var param msgraph.OutlookCategory

	if d.HasChange("color") {
		param.Color = expandCategoryColor(colorMap, d.Get("color").(string))
	}

	if err := client.ID(d.Id()).Request().Update(ctx, &param); err != nil {
		return diag.Errorf("updating Outlook Category %q: %+v", d.Get("name").(string), err)
	}

	return resourceArmCategoryRead(ctx, d, meta)
}

func resourceArmCategoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	if err := client.ID(d.Id()).Request().Delete(ctx); err != nil {
		return diag.Errorf("deleting Outlook Category %q: %+v", d.Get("name").(string), err)
	}

	return nil
}

func keySlice(m map[string]msgraph.CategoryColor) []string {
	keys := make([]string, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func expandCategoryColor(m map[string]msgraph.CategoryColor, color string) *msgraph.CategoryColor {
	return utils.ToPtr(m[color]).(*msgraph.CategoryColor)
}

func flattenCategoryColor(m map[string]msgraph.CategoryColor, color *msgraph.CategoryColor) string {
	if color == nil {
		return "None"
	}
	for k, v := range m {
		if v == *color {
			return k
		}
	}
	return "None"
}

func getCategoryByName(categories []msgraph.OutlookCategory, displayName string) *msgraph.OutlookCategory {
	for _, c := range categories {
		if c.DisplayName == nil {
			continue
		}
		if displayName == *c.DisplayName {
			return &c
		}
	}
	return nil
}
