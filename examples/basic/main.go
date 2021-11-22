package main

import (
	"log"
	"os"

	"github.com/suborbital/compute-go"
)

// a basic example without much error handling
func main() {
	token, exists := os.LookupEnv("SCC_ENV_TOKEN")
	if !exists {
		log.Panic("SCC_ENV_TOKEN not set")
	}

	client, err := compute.NewClient(compute.LocalConfig(), token)
	if err != nil {
		log.Panic(err)
	}

	// create a runnable that can be passed into compute.Client
	helloRunnable := compute.NewRunnable("com.suborbital", "acmeco", "default", "hello-world", "assemblyscript")

	// fetch an assemblyscript runnable template pre-filled with data from the above runnable
	template, err := client.BuilderTemplate(helloRunnable)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("building with template:\n%s\n", template.Contents)

	// trigger a remote build of the runnable
	buildResult, _ := client.BuildFunctionString(helloRunnable, template.Contents)

	// if the build succeeds, run it
	if buildResult.Succeeded {
		log.Println("build succeeded!")

		client.PromoteDraft(helloRunnable)

		payload := "world!"

		log.Printf("Executing runnable with payload: '%s'\n", payload)
		output, _ := client.ExecString(helloRunnable, payload)

		log.Println(string(output))
	}
}
