package main

import (
	"log"
	"os"

	"github.com/suborbital/compute-go"
)

func main() {
	token, exists := os.LookupEnv("SCC_ENV_TOKEN")
	if !exists {
		log.Fatal("SCC_ENV_TOKEN environment variable not set")
	}

	client, err := compute.NewLocalClient(token)
	if err != nil {
		log.Fatal(err)
	}

	// get a list of Runnables
	runnables, err := client.UserFunctions("userID", "namespace")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range runnables {
		log.Println(r.FQFN)
	}
}
