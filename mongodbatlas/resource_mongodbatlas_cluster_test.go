package mongodbatlas

import (
	"errors"
	"fmt"
	"testing"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongodbatlasCluster_basic(t *testing.T) {
	var cluster ma.Cluster
	projectName := "test"
	clusterName := fmt.Sprintf("test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	size := "M10"
	diskSize := "10"

	resourceName := "mongodbatlas_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasCluster(projectName, clusterName, size, diskSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttrSet(resourceName, "disk_size_gb"),
					resource.TestCheckResourceAttrSet(resourceName, "identifier"),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttrSet(resourceName, "mongodb_version"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri_updated"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri_with_options"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "mongodb_major_version", "3.6"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "region", "US_EAST_1"),
					resource.TestCheckResourceAttr(resourceName, "size", size),
					resource.TestCheckResourceAttr(resourceName, "backup", "false"),
					resource.TestCheckResourceAttr(resourceName, "disk_gb_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "provider_backup", "false"),
					resource.TestCheckResourceAttr(resourceName, "num_shards", "1"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_factor", "3"),
				),
			},
			{
				Config: testAccMongodbatlasCluster(projectName, clusterName, "M20", "20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "size", "M20"),
				),
			},
			{
				Config: testAccMongodbatlasCluster(projectName, clusterName, "M20", "35"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "35"),
				),
			},
		},
	})
}

func TestMongodbatlasCluster_importBasic(t *testing.T) {
	projectName := "test"
	projectID := "5ba8c5c396e8211ae8272486"
	clusterName := fmt.Sprintf("test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	size := "M10"
	diskSize := "10"
	importStateID := fmt.Sprintf("%s-%s", projectID, clusterName)

	resourceName := "mongodbatlas_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasCluster(projectName, clusterName, size, diskSize),
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

func testAccCheckMongodbatlasClusterExists(n string, res *ma.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cluster ID is set")
		}

		if rs.Primary.Attributes["group"] == "" {
			return errors.New("No Cluster group ID is set")
		}

		if rs.Primary.Attributes["name"] == "" {
			return errors.New("No Cluster name is set")
		}

		client := testAccProvider.Meta().(*ma.Client)

		c, _, err := client.Clusters.Get(rs.Primary.Attributes["group"], rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}

		*res = *c
		return nil
	}
}

func testAccCheckMongodbatlasClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ma.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" {
			continue
		}

		clusters, _, err := client.Clusters.List(rs.Primary.Attributes["group"])

		if err == nil {
			if len(clusters) != 0 {
				return fmt.Errorf("Cluster %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if err != nil {
			return fmt.Errorf("Error listing MongoDB Clusters: %s", err)
		}
	}

	return nil
}

func testAccMongodbatlasCluster(projectName, clusterName, size, diskSize string) string {
	return fmt.Sprintf(`resource "mongodbatlas_cluster" "test" {
  name = "%s"
  group = "${data.mongodbatlas_project.test.id}"
  mongodb_major_version = "3.6"
  provider_name = "AWS"
  region = "US_EAST_1"
  size = "%s"
  disk_size_gb = "%s"
  backup = false
  disk_gb_enabled = false
}

data "mongodbatlas_project" "test" {
  name = "%s"
}`, clusterName, size, diskSize, projectName)
}
