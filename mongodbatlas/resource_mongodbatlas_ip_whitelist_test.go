package mongodbatlas

import (
	"errors"
	"fmt"
	"testing"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongodbatlasIPWhitelist_basic(t *testing.T) {
	var whitelist ma.Whitelist
	projectName := "test"
	cidrBlock := "179.154.224.127/32"
	comment := "running acceptance tests"

	resourceName := "mongodbatlas_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasWhitelistCidr(projectName, cidrBlock, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasWhitelistExists(resourceName, &whitelist),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongodbatlasWhitelistCidr(projectName, cidrBlock, "testing changing comments"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasWhitelistExists(resourceName, &whitelist),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", "testing changing comments"),
				),
			},
		},
	})
}

func TestAccAWSEcsWhitelist_importBasic(t *testing.T) {
	projectName := "test"
	projectID := "5ba8c5c396e8211ae8272486"
	cidrBlock := "179.154.224.127/32"
	comment := "running acceptance tests"
	importStateID := fmt.Sprintf("%s-%s", projectID, cidrBlock)

	resourceName := "mongodbatlas_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasWhitelistCidr(projectName, cidrBlock, comment),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     importStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongodbatlasWhitelistExists(n string, res *ma.Whitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Whitelist ID is set")
		}

		if rs.Primary.Attributes["group"] == "" {
			return errors.New("No Whitelist group ID is set")
		}

		if rs.Primary.Attributes["cidr_block"] == "" {
			return errors.New("No Whitelist CIDR Block is set")
		}

		client := testAccProvider.Meta().(*ma.Client)

		c, _, err := client.Whitelist.Get(rs.Primary.Attributes["group"], rs.Primary.Attributes["cidr_block"])
		if err != nil {
			return err
		}

		*res = *c
		return nil
	}
}

func testAccCheckMongodbatlasWhitelistDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ma.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_ip_whitelist" {
			continue
		}

		whitelists, _, err := client.Whitelist.List(rs.Primary.Attributes["group"])

		if err == nil {
			if len(whitelists) != 0 {
				return fmt.Errorf("Whitelist %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if err != nil {
			return fmt.Errorf("Error listing MongoDB Whitelists: %s", err)
		}
	}

	return nil
}

func testAccMongodbatlasWhitelistCidr(projectName, cidrBlock, comment string) string {
	return fmt.Sprintf(`resource "mongodbatlas_ip_whitelist" "test" {
  group = "${data.mongodbatlas_project.test.id}"
  cidr_block = "%s"
  comment = "%s"
}

data "mongodbatlas_project" "test" {
  name = "%s"
}`, cidrBlock, comment, projectName)
}
