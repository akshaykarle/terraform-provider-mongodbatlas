package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMongodbatlasProject_basic(t *testing.T) {
	projectName := "testAcc"
	resourceName := "mongodbatlas_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasProject(projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_count"),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
				),
			},
		},
	})
}

func testAccMongodbatlasProject(projectName string) string {
	return fmt.Sprintf(`resource "mongodbatlas_project" "test" {
  org_id = "5b71ff2f96e82120d0aaec14"
  name = "%s"
}`, projectName)
}
