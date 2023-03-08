package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/suborbital/se2-go"
)

func templates(client *se2.Client2) {
	ctx, cxl := context.WithTimeout(context.Background(), 20*time.Second)
	defer cxl()
	templates, err := client.ListTemplates(ctx)
	if err != nil {
		log.Fatalf("client.ListTemplates: %s", err.Error())
	}

	fmt.Printf("templates are:\n\n%#v\n", templates)

	if len(templates.Templates) == 0 {
		log.Fatalf("got empty templates list, should have the defaults in it")
	}

	template, err := client.GetTemplate(ctx, templates.Templates[0].Name)
	if err != nil {
		log.Fatalf("client.GetTemplate: %s", err.Error())
	}

	fmt.Printf("template:\n%#v\n", template)

}
