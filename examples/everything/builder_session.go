package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/suborbital/se2-go"
)

func builderSession(client *se2.Client2) {
	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	m := ulid.Make()

	printHeader("creating tenant for session")
	sessionTenant, err := client.CreateTenant(ctx, m.String(), "description for session tenant")
	if err != nil {
		log.Fatalf("creating tenant for session failed with %s", err.Error())
	}

	printHeader(fmt.Sprintf("creating session for tenant '%s', namespace 'everythingns' and plugin 'everythingbagel'", sessionTenant.Name))
	s, err := client.CreateSession(ctx, sessionTenant.Name, "everythingns", "everythingbagel")
	if err != nil {
		log.Fatalf("creating a new session for 'everythingns.everythingbagel' failed with %s", err.Error())
	}

	log.Printf("session token is this\n%s\n", s.Token)

	printHeader("getting builder features")
	features, err := client.GetBuilderFeatures(ctx)
	if err != nil {
		log.Fatalf("get builder feautres failed: %s", err.Error())
	}

	log.Printf("%#v", features)

	printHeader("list templates")
	templates, err := client.ListTemplates(ctx)
	if err != nil {
		log.Fatalf("listing templates failed: %s", err.Error())
	}

	fmt.Printf("available templates:\n%#v\n", templates)

	if len(templates.Templates) == 0 {
		log.Fatal("there are no available templates")
	}

	for _, localTemplate := range templates.Templates {
		printHeader(fmt.Sprintf("start a draft for template named '%s'", localTemplate.Name))

		draft, err := client.CreatePluginDraft(ctx, localTemplate.Name, s)
		if err != nil {
			log.Fatalf("create plugin draft failed with error '%s': Name: %s", err.Error(), localTemplate.Name)
		}

		fmt.Printf("returned draft response is\n\n%#v\n\n", draft)
	}

	printHeader("getting plugin drafts saved on session")

	pd, err := client.GetPluginDraft(ctx, s)
	if err != nil {
		log.Fatalf("getting plugin drafts failed: %s", err.Error())
	}

	log.Printf("plugin draft is\n%#v\n", pd)

	err = client.DeleteTenantByName(ctx, sessionTenant.Name)
	if err != nil {
		log.Fatalf("could not delete tenant by name %s", sessionTenant.Name)
	}
}
