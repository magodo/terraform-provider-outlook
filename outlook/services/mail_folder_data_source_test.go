package services_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccMailFolderDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDsMailFolderConfig_basic(acctest.RandString(3)),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckMailFolderExists(t, "name"),
					resource.TestCheckResourceAttrSet("data.outlook_mail_folder.test", "name"),

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
`,  testAccMailFolderConfig_basic(suffix))
}
