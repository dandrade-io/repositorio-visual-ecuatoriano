package embedding

import (
	"net/http"
	"repository/config"
)

var client *http.Client

func init() {
	if config.CAPath == "" {
		client = http.DefaultClient
		return
	}

	c, err := config.ClientWithPEM(config.CAPath)
	if err != nil {
		panic(err)
	}

	client = c
}
