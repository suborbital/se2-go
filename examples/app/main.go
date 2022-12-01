package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	se2 "github.com/suborbital/se2-go"
)

func main() {
	client := client()

	// This is a local reference to some plugin. Nothing has run in SE2 at this point.
	plugin := se2.NewPlugin("com.suborbital", "acmeco", "default", "tinygo-hey")

	// Request template source code for the above plugin.
	template, _ := client.BuilderTemplate(plugin, "tinygo")

	// Modify the default template
	modified := strings.Replace(template.Contents, "Hello", "Hey there", 1)
	log.Println(modified)

	// Run a remote build for the provided plugin and the modified 'goodbye world'
	// template.
	build, err := client.BuildPluginString(plugin, "tinygo", modified)

	if err != nil {
		log.Fatal(err)
	}

	if !build.Succeeded {
		// Log the builder output to see why the build failed
		log.Fatal(build.OutputLog)
	}

	// Deploy the plugin and get the new reference
	ref, _ := client.PromoteDraft(plugin)

	// Hello!
	time.Sleep(time.Second * 2)
	for _, name := range []string{"Europa", "Io", "Ganymede", "Callisto"} {
		time.Sleep(time.Millisecond * 300)
		result, _, err := client.ExecRefString(ref.Version, name)
		if err != nil {
			log.Fatal(err)
		}

		// Log the execution output
		fmt.Println(string(result))
	}
}
