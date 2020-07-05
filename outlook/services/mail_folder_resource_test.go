package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccMailFolderResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMailFolderConfig_basic(acctest.RandString(3)),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_mail_folder.test"),
		},
	})
}

func testAccMailFolderConfig_basic(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_mail_folder" "test" {
	name = "foo%s"
}
`, suffix)
}
