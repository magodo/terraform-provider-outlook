package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOutlookCategoryDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOutlookCategory_basic(acctest.RandString(3)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outlook_category.test", "name"),
					resource.TestCheckResourceAttr("data.outlook_category.test", "color", "Red"),
				),
			},
		},
	})
}

func testAccDsOutlookCategory_basic(suffix string) string {
	return fmt.Sprintf(`
%s

data "outlook_category" "test" {
  name = outlook_category.test.name
}
`, testAccOutlookCategory_basic(suffix))
}
