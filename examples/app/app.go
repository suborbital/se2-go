package main

import (
	"log"

	"github.com/suborbital/compute-go"
)

func main() {
	client := client()

	// This is a local reference to some Runnable. Nothing has run in Compute at this point.
	runnable := compute.NewRunnable("com.suborbital", "acmeco", "default", "rs-hello-world", "rust")

	// Request template source code for the above Runnable.
	template, _ := client.BuilderTemplate(runnable)

	// Log the default 'hello world' Rust template to stdout
	log.Println(template.Contents)

	// Run a remote build for the provided Runnable and the unmodified 'hello world'
	// template source code.
	build, _ := client.BuildFunctionString(runnable, template.Contents)

	if !build.Succeeded {
		// Log the builder output to see why the build failed
		log.Fatal(build.OutputLog)
	}

	// Deploy the function (the runnable's .Version field is adjusted here)
	client.PromoteDraft(runnable)

	// Execute the function
	result, _ := client.ExecString(runnable, "world!")

	// Log the execution output
	log.Println(string(result))
}
