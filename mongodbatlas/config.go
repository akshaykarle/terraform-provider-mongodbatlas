package mongodbatlas

import (
	"net/http"

	dac "github.com/akshaykarle/go-http-digest-auth-client"
	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
)

type Config struct {
	AtlasUsername string
	AtlasAPIKey   string
}

func (c *Config) NewClient() *ma.Client {
	t := dac.NewTransport(c.AtlasUsername, c.AtlasAPIKey)
	httpClient := &http.Client{Transport: &t}
	client := ma.NewClient(httpClient)
	return client
}
