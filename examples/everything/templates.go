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
}
