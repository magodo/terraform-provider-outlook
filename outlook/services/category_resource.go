package services

import (
	"context"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/validation"
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
				ForceNew: true,
			},

			"color": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.StringInSlice(func(m map[string]msgraph.CategoryColor) []string {
					keys := make([]string, len(m))
					for k := range m {
						keys = append(keys, k)
					}
					return keys
				}(colorMap), false),
			},
		},
	}
}

var (
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

	return resourceArmCategoryRead(ctx, d, meta)
}

func resourceArmCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	return nil
}

func resourceArmCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	return resourceArmCategoryRead(ctx, d, meta)
}

func resourceArmCategoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Categories

	return nil
}
