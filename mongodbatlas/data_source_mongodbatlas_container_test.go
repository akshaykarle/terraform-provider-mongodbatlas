package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMongodbatlasDataSource_Container(t *testing.T) {
	projectName := "test"
	cidrBlock := "10.0.0.0/21"

	dataSourceName := "data.mongodbatlas_container.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasDataSourceContainerCidr(projectName, cidrBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "identifier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),
					resource.TestCheckResourceAttr(dataSourceName, "atlas_cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(dataSourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "US_EAST_1"),
				),
			},
		},
	})
}

func testAccMongodbatlasDataSourceContainerCidr(projectName, cidrBlock string) string {
	return fmt.Sprintf(`resource "mongodbatlas_container" "test" {
  group = "${data.mongodbatlas_project.test.id}"
  atlas_cidr_block = "%s"
  provider_name = "AWS"
  region = "US_EAST_1"
}

data "mongodbatlas_container" "test" {
  group = "${data.mongodbatlas_project.test.id}"
  identifier = "${mongodbatlas_container.test.identifier}"
}

data "mongodbatlas_project" "test" {
  name = "%s"
}`, cidrBlock, projectName)
}
