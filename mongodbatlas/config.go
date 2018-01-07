package mongodbatlas

import (
	"net/http"

	dac "github.com/akshaykarle/go-http-digest-auth-client"
	"github.com/akshaykarle/mongodb-atlas-go/mongodb"
)

type Config struct {
	AtlasUsername string
	AtlasAPIKey   string
}

func (c *Config) NewClient() *mongodb.Client {
	t := dac.NewTransport(c.AtlasUsername, c.AtlasAPIKey)
	httpClient := &http.Client{Transport: &t}
	client := mongodb.NewClient(httpClient)
	return client
}
