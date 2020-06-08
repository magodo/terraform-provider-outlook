package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/provider/client"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Description: `
Outlook data source can be used to retrieve the ID for a user by name.
`,

		Read: dataUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the user to look up.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"principal_name": {
				Description: "The user principal name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(client.Client)

	name := d.Get("name").(string)

	users, err := c.ListUser(context.TODO())
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.DisplayName == name {
			d.SetId(u.ID)
			d.Set("principal_name", u.PrincipalName)
			return nil
		}
	}

	return fmt.Errorf("user not found with name %s", name)
}
