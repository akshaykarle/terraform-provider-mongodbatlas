package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseUserCreate,
		Read:   resourceDatabaseUserRead,
		Update: resourceDatabaseUserUpdate,
		Delete: resourceDatabaseUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDatabaseUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},
						"collection": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.DatabaseUser{
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		DatabaseName: d.Get("database").(string),
	}

	params.Roles = readRolesFromSchema(d.Get("roles").([]interface{}))

	databaseUser, _, err := client.DatabaseUsers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB DatabaseUser: %s", err)
	}
	d.SetId(databaseUser.Username)
	log.Printf("[INFO] MongoDB DatabaseUser ID: %s", d.Id())

	return resourceDatabaseUserRead(d, meta)
}

func resourceDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	u, _, err := client.DatabaseUsers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB DatabaseUser %s (%s): %s", d.Id(), d.Get("group").(string), err)
	}

	d.Set("username", u.Username)
	d.Set("database", u.DatabaseName)
	rolesMap := make([]map[string]interface{}, len(u.Roles))
	for i, r := range u.Roles {
		rolesMap[i] = map[string]interface{}{
			"name":       r.RoleName,
			"database":   r.DatabaseName,
			"collection": r.CollectionName,
		}
	}
	d.Set("roles", rolesMap)

	return nil
}

func resourceDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	u, _, err := client.DatabaseUsers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB DatabaseUser %s: %s", d.Id(), err)
	}

	if d.HasChange("password") {
		u.Password = d.Get("password").(string)
		requestUpdate = true
	}
	if d.HasChange("roles") {
		u.Roles = readRolesFromSchema(d.Get("roles").([]interface{}))
		requestUpdate = true
	}

	if requestUpdate {
		_, _, err := client.DatabaseUsers.Update(d.Get("group").(string), d.Id(), u)
		if err != nil {
			return fmt.Errorf("Error updating MongoDB DatabaseUser %s: %s", d.Id(), err)
		}
		return resourceDatabaseUserRead(d, meta)
	}
	return nil
}

func resourceDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB DatabaseUser destroy: %v", d.Id())
	_, err := client.DatabaseUsers.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB DatabaseUser %s: %s", d.Id(), err)
	}

	return nil
}

func resourceDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ma.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a user, use the format {group id}-{username}")
	}
	gid := parts[0]
	username := parts[1]

	u, _, err := client.DatabaseUsers.Get(gid, username)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import user %s in group %s, error: %s", username, gid, err.Error())
	}

	d.SetId(u.Username)
	d.Set("group", u.GroupID)

	return []*schema.ResourceData{d}, nil
}

func readRolesFromSchema(rolesMap []interface{}) (roles []ma.Role) {
	roles = make([]ma.Role, len(rolesMap))
	for i, r := range rolesMap {
		roleMap := r.(map[string]interface{})

		roles[i] = ma.Role{
			RoleName:       roleMap["name"].(string),
			DatabaseName:   roleMap["database"].(string),
			CollectionName: roleMap["collection"].(string),
		}
	}
	return roles
}
