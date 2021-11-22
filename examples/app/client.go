package main

import (
	"os"

	"github.com/suborbital/compute-go"
)

func client() *compute.Client {
	token, _ := os.LookupEnv("SCC_ENV_TOKEN")
	client, _ := compute.NewClient(compute.LocalConfig(), token)

	return client
}
