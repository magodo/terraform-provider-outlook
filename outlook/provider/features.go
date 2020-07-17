package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
)

var featureSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mail_folder_delete_parallelism": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Specify the parallelism used at the that point when mail folder moves its containing messages back to inbox before deletion.",
			},
		},
	},
	Description: "Provider level features",
}

func expandFeature(input []interface{}) clients.UserFeature {
	if len(input) == 0 || input[0] == nil {
		return clients.UserFeature{
			MailFolderDeleteParallelism: 1,
		}
	}

	raw := input[0].(map[string]interface{})

	return clients.UserFeature{
		MailFolderDeleteParallelism: raw["mail_folder_delete_parallelism"].(int),
	}
}
