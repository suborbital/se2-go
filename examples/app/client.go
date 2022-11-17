package main

import (
	"os"

	"github.com/suborbital/se2-go"
)

func client() *se2.Client {
	token, _ := os.LookupEnv("SCC_ENV_TOKEN")
	client, _ := se2.NewClient(se2.LocalConfig(), token)

	return client
}
