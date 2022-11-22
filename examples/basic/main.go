package main

import (
	"log"
	"os"

	"github.com/suborbital/se2-go"
)

var tmpl = "assemblyscript"

// a basic example without much error handling
func main() {
	token, exists := os.LookupEnv("SE2_ENV_TOKEN")
	if !exists {
		log.Panic("SE2_ENV_TOKEN not set")
	}

	client, err := se2.NewClient(se2.LocalConfig(), token)
	if err != nil {
		log.Panic(err)
	}

	// create a module that can be passed into se2.Client
	helloModule := se2.NewModule("com.suborbital", "acmeco", "default", "hello-world")

	// fetch an assemblyscript module template pre-filled with data from the above module
	template, err := client.BuilderTemplate(helloModule, tmpl)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("building with template:\n%s\n", template.Contents)

	// trigger a remote build of the module
	buildResult, _ := client.BuildFunctionString(helloModule, tmpl, template.Contents)

	// if the build succeeds, run it
	if buildResult.Succeeded {
		log.Println("build succeeded!")

		client.PromoteDraft(helloModule)

		payload := "world!"

		log.Printf("Executing module with payload: '%s'\n", payload)
		output, _, _ := client.ExecString(helloModule, payload)

		log.Println(string(output))
	}
}
