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

func TestAccMessageRuleResource_multipleRuleSequence(t *testing.T) {
	suffix := acctest.RandString(3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccMessageRuleConfig_multipleRuleSeq1(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test1"),
			importStep("outlook_message_rule.test2"),
			importStep("outlook_message_rule.test3"),
			{
				Config: testAccMessageRuleConfig_multipleRuleSeq2(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test1"),
			importStep("outlook_message_rule.test2"),
			importStep("outlook_message_rule.test3"),
			{
				Config: testAccMessageRuleConfig_multipleRuleSeq1(suffix),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckMailFolderExists(t, "name"),
				// ),
			},
			importStep("outlook_message_rule.test1"),
			importStep("outlook_message_rule.test2"),
			importStep("outlook_message_rule.test3"),
		},
	})
}

func testAccMessageRuleConfig_basic(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_message_rule" "test" {
  name     = "msgrule-%[1]s"
  sequence = "1"
  enabled  = false
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
  name     = "msgrule-%[1]s"
  sequence = "2"
  enabled  = true
  condition {
    from_addresses = [
      "foo@bar.com"
    ]
    importance         = "high"
    is_meeting_request = true
  }
  action {
    mark_as_read   = true
    copy_to_folder = outlook_mail_folder.test.id
  }
}
`, suffix)
}

func testAccMessageRuleConfig_multipleRuleSeq1(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_message_rule" "test1" {
  name     = "msgrule-%[1]s1"
  sequence = "1"
  enabled  = false
  action {
    mark_as_read = true
  }
}
resource "outlook_message_rule" "test2" {
  name     = "msgrule-%[1]s2"
  sequence = "2"
  enabled  = false
  action {
    mark_as_read = true
  }
  depends_on = [outlook_message_rule.test1]
}
resource "outlook_message_rule" "test3" {
  name     = "msgrule-%[1]s3"
  sequence = "3"
  enabled  = false
  action {
    mark_as_read = true
  }
  depends_on = [outlook_message_rule.test2]
}
`, suffix)
}

func testAccMessageRuleConfig_multipleRuleSeq2(suffix string) string {
	return fmt.Sprintf(`
resource "outlook_message_rule" "test1" {
  name     = "msgrule-%[1]s1"
  sequence = "3"
  enabled  = false
  action {
    mark_as_read = true
  }
  depends_on = [outlook_message_rule.test3]
}
resource "outlook_message_rule" "test2" {
  name     = "msgrule-%[1]s2"
  sequence = "1"
  enabled  = false
  action {
    mark_as_read = true
  }
}
resource "outlook_message_rule" "test3" {
  name     = "msgrule-%[1]s3"
  sequence = "2"
  enabled  = false
  action {
    mark_as_read = true
  }
  depends_on = [outlook_message_rule.test2]
}
`, suffix)
}
