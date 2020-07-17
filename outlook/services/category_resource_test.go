package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOutlookCategory_basic(t *testing.T) {
	suffix := acctest.RandString(3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccOutlookCategory_basic(suffix),
				//Check: resource.ComposeTestCheckFunc(),
			},
			importStep("outlook_category.test"),
		},
	})
}

func TestAccOutlookCategory_update(t *testing.T) {
	suffix := acctest.RandString(3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccOutlookCategory_basic(suffix),
				//Check: resource.ComposeTestCheckFunc(),
			},
			importStep("outlook_category.test"),
			{
				Config: testAccOutlookCategory_update(suffix),
				//Check: resource.ComposeTestCheckFunc(),
			},
			importStep("outlook_category.test"),
		},
	})
}

func testAccOutlookCategory_basic(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_category" "test" {
  name  = "category-%s"
  color = "Black"
}
`, suffix)
}

func testAccOutlookCategory_update(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_category" "test" {
  name  = "category-%s"
  color = "Brown"
}
`, suffix)
}
