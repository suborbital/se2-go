package main

import (
	"log"

	"github.com/suborbital/compute-go"
)

func main() {
	client, err := compute.NewLocalClient()
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
