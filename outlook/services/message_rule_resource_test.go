package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMessageRuleResource_basic(t *testing.T) {
	suffix := acctest.RandString(3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMessageRuleConfig_basic(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test"),
		},
	})
}

func TestAccMessageRuleResource_upgrade(t *testing.T) {
	suffix := acctest.RandString(3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMessageRuleConfig_basic(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test"),
			{
				Config: testAccMessageRuleConfig_complete(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test"),
			{
				Config: testAccMessageRuleConfig_basic(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test"),
		},
	})
}

func testAccMessageRuleConfig_basic(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_message_rule" "test" {
  name = "msgrule-%[1]s"
  sequence = "1"
  enabled = false
  action {
    mark_as_read = true
  }
}
`, suffix)
}

func testAccMessageRuleConfig_complete(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_mail_folder" "test" {
	name = "msgrule-%[1]s"
}

resource "outlook_message_rule" "test" {
  name = "msgrule-%[1]s"
  sequence = "2"
  enabled = true
  condition {
    from_addresses = [
      "foo@bar.com"
    ]
    importance = "high"
    is_meeting_request = true
  }
  action {
    mark_as_read = true
    copy_to_folder = outlook_mail_folder.test.id
  }
}
`, suffix)
}
