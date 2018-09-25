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

func TestAccMongodbatlasDatabaseUser_basic(t *testing.T) {
	var databaseUser ma.DatabaseUser
	projectName := "test"
	databaseUserName := fmt.Sprintf("test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	databaseUserPassword := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	roleName := "read"

	resourceName := "mongodbatlas_database_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasDatabaseUser(projectName, databaseUserName, databaseUserPassword, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttr(resourceName, "username", databaseUserName),
					resource.TestCheckResourceAttr(resourceName, "password", databaseUserPassword),
					resource.TestCheckResourceAttr(resourceName, "database", "admin"),
				),
			},
			{
				Config: testAccMongodbatlasDatabaseUser(projectName, databaseUserName, databaseUserPassword, "readWrite"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbatlasDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttrSet(resourceName, "group"),
					resource.TestCheckResourceAttr(resourceName, "username", databaseUserName),
					resource.TestCheckResourceAttr(resourceName, "password", databaseUserPassword),
					resource.TestCheckResourceAttr(resourceName, "database", "admin"),
				),
			},
		},
	})
}

func TestAccAWSEcsDatabaseUser_importBasic(t *testing.T) {
	projectName := "test"
	projectID := "5ba8c5c396e8211ae8272486"
	databaseUserName := fmt.Sprintf("test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	databaseUserPassword := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	roleName := "read"
	importStateID := fmt.Sprintf("%s-%s", projectID, databaseUserName)

	resourceName := "mongodbatlas_database_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbatlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbatlasDatabaseUser(projectName, databaseUserName, databaseUserPassword, roleName),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           importStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckMongodbatlasDatabaseUserExists(n string, res *ma.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No DatabaseUser ID is set")
		}

		if rs.Primary.Attributes["group"] == "" {
			return errors.New("No DatabaseUser group ID is set")
		}

		client := testAccProvider.Meta().(*ma.Client)

		c, _, err := client.DatabaseUsers.Get(rs.Primary.Attributes["group"], rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *c
		return nil
	}
}

func testAccCheckMongodbatlasDatabaseUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ma.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user" {
			continue
		}

		databaseUsers, _, err := client.DatabaseUsers.List(rs.Primary.Attributes["group"])

		if err == nil {
			if len(databaseUsers) != 0 {
				return fmt.Errorf("DatabaseUser %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if err != nil {
			return fmt.Errorf("Error listing MongoDB DatabaseUsers: %s", err)
		}
	}

	return nil
}

func testAccMongodbatlasDatabaseUser(projectName, databaseUserName, databaseUserPassword, roleName string) string {
	return fmt.Sprintf(`resource "mongodbatlas_database_user" "test" {
  username = "%s"
  password = "%s"
  group = "${data.mongodbatlas_project.test.id}"
  database = "admin"
  roles  = [
    {
      name = "%s"
      database = "admin"
    }
  ]
}

data "mongodbatlas_project" "test" {
  name = "%s"
}`, databaseUserName, databaseUserPassword, roleName, projectName)
}
