package main

import (
	"log"
	"os"

	"github.com/suborbital/compute-go"
)

// a basic example without much error handling
func main() {
	token, _ := os.LookupEnv("SCC_ENV_TOKEN")

	conf, _ := compute.DefaultConfig("http://localhost") // use your own base URL here
	client, _ := compute.NewClient(conf, token)

	// create a runnable that can be passed into compute.Client
	helloRunnable := compute.NewRunnable("com.suborbital", "acmeco", "default", "hello-world", "assemblyscript")

	// fetch an assemblyscript runnable template pre-filled with data from the above runnable
	template, _ := client.BuilderTemplate(helloRunnable)
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
