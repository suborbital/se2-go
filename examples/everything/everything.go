package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/suborbital/se2-go"
)

func main() {
	token, ok := os.LookupEnv("EVERYTHING_TOKEN")
	if !ok {
		log.Fatalf("There is no token set. Get an environment access key from the se2 admin area, and set it to the EVERYTHING_TOKEN env var.")
	}

	// Set up the client to point at staging (admin and builder) with a valid access token.
	client, err := se2.NewClient2(se2.ModeStaging, token)
	if err != nil {
		log.Fatalf("encountered new client error: %s", err.Error())
	}

	// Function
	// tenants(client)

	// templates(client)

	builderSession(client)

}

func printHeader(msg string) {
	sep := strings.Repeat("=", len(msg))

	fmt.Printf("\n"+
		"%s\n"+
		"%s\n"+
		"%s\n\n", sep, msg, sep)
}
