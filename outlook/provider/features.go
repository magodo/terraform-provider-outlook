package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/clients"
)

var featureSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	MinItems:    1,
	Elem:        &schema.Resource{},
	Description: "Provider level features",
}

func expandFeature(input []interface{}) clients.UserFeature {
	if len(input) == 0 || input[0] == nil {
		return clients.UserFeature{}
	}

	//raw := input[0].(map[string]interface{})

	return clients.UserFeature{}
}
