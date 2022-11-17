package main

import (
	"log"
	"os"

	"github.com/suborbital/se2-go"
)

func main() {
	token, exists := os.LookupEnv("SCC_ENV_TOKEN")
	if !exists {
		log.Fatal("SCC_ENV_TOKEN environment variable not set")
	}

	client, err := se2.NewLocalClient(token)
	if err != nil {
		log.Fatal(err)
	}

	// get a list of Modules
	modules, err := client.UserFunctions("userID", "namespace")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range modules {
		log.Println(r.FQMN)
	}
}
