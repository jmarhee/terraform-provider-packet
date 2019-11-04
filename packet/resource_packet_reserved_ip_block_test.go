package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckPacketReservedIPBlockConfig_Global(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "tfacc-reserved_ip_block-%s"
}

resource "packet_reserved_ip_block" "test" {
    project_id = "${packet_project.foobar.id}"
    type     = "global_ipv4"
	description = "testdesc"
	quantity = 1
}`, name)
}

func testAccCheckPacketReservedIPBlockConfig_Public(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "tfacc-reserved_ip_block-%s"
}

resource "packet_reserved_ip_block" "test" {
    project_id  = "${packet_project.foobar.id}"
    facility    = "ewr1"
    type        = "public_ipv4"
	quantity    = 2
}`, name)
}

func TestAccPacketReservedIPBlock_Global(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketReservedIPBlockConfig_Global(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "quantity", "1"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "description", "testdesc"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "type", "global_ipv4"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "netmask", "255.255.255.255"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "management", "false"),
				),
			},
		},
	})
}

func TestAccPacketReservedIPBlock_Public(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketReservedIPBlockConfig_Public(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "facility", "ewr1"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "type", "public_ipv4"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "quantity", "2"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "netmask", "255.255.255.254"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "management", "false"),
				),
			},
		},
	})
}

func TestAccPacketReservedIPBlock_importBasic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketReservedIPBlockConfig_Public(rs),
			},
			{
				ResourceName:      "packet_reserved_ip_block.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketReservedIPBlockDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_reserved_ip_block" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Reserved IP block still exists")
		}
	}

	return nil
}
