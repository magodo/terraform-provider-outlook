package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMailFolderResource_basic(t *testing.T) {
	suffix := acctest.RandString(3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMailFolderConfig_basic(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_mail_folder.test"),
		},
	})
}

func TestAccMailFolderResource_parent(t *testing.T) {
	suffix := acctest.RandString(3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMailFolderConfig_parent_foo(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_mail_folder.test_child"),
			{
				Config: testAccMailFolderConfig_parent_bar(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_mail_folder.test_child"),
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

func testAccMailFolderConfig_parent_foo(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_mail_folder" "test_foo" {
  name = "foo%[1]s"
}
resource "outlook_mail_folder" "test_bar" {
  name = "bar%[1]s"
}

resource "outlook_mail_folder" "test_child" {
  name = "child%[1]s"
  parent_folder_id = outlook_mail_folder.test_foo.id
}
`, suffix)
}

func testAccMailFolderConfig_parent_bar(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_mail_folder" "test_foo" {
  name = "foo%[1]s"
}
resource "outlook_mail_folder" "test_bar" {
  name = "bar%[1]s"
}

resource "outlook_mail_folder" "test_child" {
  name = "child%[1]s"
  parent_folder_id = outlook_mail_folder.test_bar.id
}
`, suffix)
}
