package main

import (
	"log"

	se2 "github.com/suborbital/se2-go"
)

func main() {
	client := client()

	// This is a local reference to some Module. Nothing has run in Compute at this point.
	module := se2.NewModule("com.suborbital", "acmeco", "default", "rs-hello-world")

	// Request template source code for the above Module.
	template, _ := client.BuilderTemplate(module)

	// Log the default 'hello world' Rust template to stdout
	log.Println(template.Contents)

	// Run a remote build for the provided Module and the unmodified 'hello world'
	// template source code.
	build, _ := client.BuildFunctionString(module, template.Contents)

	if !build.Succeeded {
		// Log the builder output to see why the build failed
		log.Fatal(build.OutputLog)
	}

	// Deploy the function (the module's .Version field is adjusted here)
	client.PromoteDraft(module)

	// Execute the function
	result, _, _ := client.ExecString(module, "world!")

	// Log the execution output
	log.Println(string(result))
}
