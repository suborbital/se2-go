package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/suborbital/se2-go"
)

func builderSession(c *se2.Client2) {
	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	m := ulid.Make()

	printHeader("creating tenant for session")
	sessionTenant, err := c.CreateTenant(ctx, m.String(), "description for session tenant")
	if err != nil {
		log.Fatalf("creating tenant for session failed with %s", err.Error())
	}

	printHeader(fmt.Sprintf("creating session for tenant '%s', namespace 'everythingns' and plugin 'everythingbagel'", sessionTenant.Name))
	s, err := c.CreateSession(ctx, sessionTenant.Name, "everythingns", "everythingbagel")
	if err != nil {
		log.Fatalf("creating a new session for 'everythingns.everythingbagel' failed with %s", err.Error())
	}

	log.Printf("session token is this\n%s\n", s.Token)

	printHeader("getting builder features")
	features, err := c.GetBuilderFeatures(ctx)
	if err != nil {
		log.Fatalf("get builder feautres failed: %s", err.Error())
	}

	log.Printf("%#v", features)

	printHeader("getting plugin drafts")

	pd, err := c.GetPluginDraft(ctx, s)
	if err != nil {
		log.Fatalf("getting plugin drafts failed: %s", err.Error())
	}

	log.Printf("plugin draft is\n%#v\n", pd)

	err = c.DeleteTenantByName(ctx, sessionTenant.Name)
	if err != nil {
		log.Fatalf("could not delete tenant by name %s", sessionTenant.Name)
	}
}
