package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMongodbatlasDataSource_Project(t *testing.T) {
	projectName := "test"
	dataSourceName := "data.mongodbatlas_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasDataSourceProject(projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_count"),
					resource.TestCheckResourceAttr(dataSourceName, "name", projectName),
				),
			},
		},
	})
}

func testAccMongodbatlasDataSourceProject(projectName string) string {
	return fmt.Sprintf(`data "mongodbatlas_project" "test" {
  name = "%s"
}`, projectName)
}
