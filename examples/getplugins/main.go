package main

import (
	"log"
	"os"

	"github.com/suborbital/se2-go"
)

func main() {
	token, exists := os.LookupEnv("SE2_ENV_TOKEN")
	if !exists {
		log.Fatal("SE2_ENV_TOKEN environment variable not set")
	}

	client, err := se2.NewLocalClient(token)
	if err != nil {
		log.Fatal(err)
	}

	// get a list of plugins
	plugins, err := client.UserPlugins("userID", "namespace")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range plugins {
		log.Println(r.FQMN)
	}
}
