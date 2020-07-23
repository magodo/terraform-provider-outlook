package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMailFolderDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDsMailFolderConfig_basic(acctest.RandString(3)),
				Check: resource.ComposeTestCheckFunc(
					// testCheckMailFolderExists(t, "name"),
					resource.TestCheckResourceAttrSet("data.outlook_mail_folder.test", "name"),
				),
			},
		},
	})
}

func TestAccMailFolderDataSource_wellKnownName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDsMailFolderConfig_wellKnownName("inbox"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckMailFolderExists(t, "name"),
					resource.TestCheckResourceAttr("data.outlook_mail_folder.test", "well_known_name", "inbox"),
				),
			},
		},
	})
}

func testAccDsMailFolderConfig_basic(suffix string) string {
	return fmt.Sprintf(`
%s

data "outlook_mail_folder" "test" {
  name = outlook_mail_folder.test.name
}
`, testAccMailFolderConfig_basic(suffix))
}

func testAccDsMailFolderConfig_wellKnownName(name string) string {
	return fmt.Sprintf(`
data "outlook_mail_folder" "test" {
  well_known_name = "%s"
}
`, name)
}
