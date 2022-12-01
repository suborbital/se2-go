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

	// create a plugin that can be passed into se2.Client
	helloPlugin := se2.NewPlugin("com.suborbital", "acmeco", "default", "hello-world")

	// fetch an assemblyscript plugin template pre-filled with data from the above plugin
	template, err := client.BuilderTemplate(helloPlugin, tmpl)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("building with template:\n%s\n", template.Contents)

	// trigger a remote build of the plugin
	buildResult, _ := client.BuildPluginString(helloPlugin, tmpl, template.Contents)

	// if the build succeeds, run it
	if buildResult.Succeeded {
		log.Println("build succeeded!")

		client.PromoteDraft(helloPlugin)

		payload := "world!"

		log.Printf("Executing plugin with payload: '%s'\n", payload)
		output, _, _ := client.ExecString(helloPlugin, payload)

		log.Println(string(output))
	}
}
