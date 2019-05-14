package mongodbatlas

import (
	"errors"
	"fmt"
	"testing"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongodbatlasContainer_basic(t *testing.T) {
	var container ma.Container
	projectName := "test"
	cidrBlock := "10.1.0.0/21"

	resourceName := "mongodbatlas_container.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasContainerCidr(projectName, cidrBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasContainerExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttrSet(resourceName, "identifier"),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttr(resourceName, "atlas_cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "region", "US_EAST_1"),
				),
			},
			{
				Config: testAccMongodbatlasContainerCidr(projectName, "192.168.0.0/21"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasContainerExists(resourceName, &container),
					resource.TestCheckResourceAttr(resourceName, "atlas_cidr_block", "192.168.0.0/21"),
				),
			},
		},
	})
}

func testAccCheckMongodbatlasContainerExists(n string, res *ma.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Container ID is set")
		}

		if rs.Primary.Attributes["group"] == "" {
			return errors.New("No Container group ID is set")
		}

		client := testAccProvider.Meta().(*ma.Client)

		c, _, err := client.Containers.Get(rs.Primary.Attributes["group"], rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *c
		return nil
	}
}

func testAccCheckMongodbatlasContainerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ma.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_container" {
			continue
		}

		containers, _, err := client.Containers.List(rs.Primary.Attributes["group"], rs.Primary.Attributes["provider_name"])

		if err == nil {
			if len(containers) != 0 {
				return fmt.Errorf("Container %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if err != nil {
			return fmt.Errorf("Error listing MongoDB Containers: %s", err)
		}
	}

	return nil
}

func testAccMongodbatlasContainerCidr(projectName, cidrBlock string) string {
	return fmt.Sprintf(`resource "mongodbatlas_container" "test" {
  group = "${data.mongodbatlas_project.test.id}"
  atlas_cidr_block = "%s"
  provider_name = "AWS"
  region = "US_EAST_1"
}

data "mongodbatlas_project" "test" {
  name = "%s"
}`, cidrBlock, projectName)
}
