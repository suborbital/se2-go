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

	// This is a local reference to some Plugin. Nothing has run in Compute at this point.
	plugin := se2.NewPlugin("com.suborbital", "acmeco", "default", "tinygo-hey")

	// Request template source code for the above Plugin.
	template, _ := client.BuilderTemplate(plugin, "tinygo")

	// Modify the default template
	modified := strings.Replace(template.Contents, "Hello", "Hey there", 1)
	log.Println(modified)

	// Run a remote build for the provided Plugin and the modified 'goodbye world'
	// template.
	build, err := client.BuildFunctionString(plugin, "tinygo", modified)

	if err != nil {
		log.Fatal(err)
	}

	if !build.Succeeded {
		// Log the builder output to see why the build failed
		log.Fatal(build.OutputLog)
	}

	// Deploy the function and get the new reference
	ref, _ := client.PromoteDraft(plugin)

	// Roll the credits
	time.Sleep(time.Second * 2)
	for _, name := range []string{"Connor", "Dan", "Dylan", "Flaki", "Gabor", "Jagger", "Jessica", "Laura", "Nyah", "Oscar", "Ramón", "Ryan", "Taryn"} {
		time.Sleep(time.Millisecond * 300)
		result, _, err := client.ExecRefString(ref.Version, name)
		if err != nil {
			log.Fatal(err)
		}

		// Log the execution output
		fmt.Println(string(result))
	}
}
