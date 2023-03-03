package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/suborbital/se2-go"
)

func tenants(client *se2.Client2) {
	// Create a context for the requests.
	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	// List existing tenants. This should be empty.
	tenants, err := client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants L33: %s", err.Error())
	}

	printHeader("List tenants should be an empty list")
	fmt.Printf("%#v\n", tenants)

	// Create a new fake name for this run.
	printHeader("Creating a new tenant")
	sessionName := ulid.Make()
	log.Printf("Session name to start: %s\n", sessionName)

	// Create a new tenant by this name.
	newTenant, err := client.CreateTenant(ctx, sessionName.String(), "everythingDescription")
	if err != nil {
		log.Fatalf("create tenant by name errorred: %s", err.Error())
	}

	log.Printf("Created tenant is\n%#v\n", newTenant)

	printHeader("Listing tenants again to see new tenant being picked up")
	// List tenants again, this should have the new tenant in it.
	tenants, err = client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants number 2: %s", err.Error())
	}
	log.Printf("%#v\n", tenants)

	// Fetch the individual tenant by its name.
	printHeader("Fetching the same created tenant by its name")
	tenant, err := client.GetTenantByName(ctx, sessionName.String())
	if err != nil {
		log.Fatalf("get tenant by name: %s", err.Error())
	}
	log.Printf("%#v\n", tenant)

	// Update the tenant by its name to use a different name and description.
	printHeader("Updating the tenant's description")
	updated, err := client.UpdateTenantByName(ctx, sessionName.String(), "newDescription")
	if err != nil {
		log.Fatalf("update tenant errored: %s", err.Error())
	}
	log.Printf("%#v\n", updated)

	// Delete tenant by name
	printHeader("Deleting tenant by its name")
	err = client.DeleteTenantByName(ctx, sessionName.String())
	if err != nil {
		log.Fatalf("deleting tenant by its new name should not have failed. It did: %s", err.Error())
	}

	// List tenants should be empty after this
	printHeader("Listing tenants again to see that the tenant no longer exists")
	listTenants, err := client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants after deletion should have succeeded. It failed with %s", err.Error())
	}
	log.Printf("this is the list tenants (length: %d)\n%#v\n", len(listTenants.Tenants), listTenants)

	// Getting tenant by its new name should fail
	printHeader("Fetching tenant by its name to see that it fails")
	deleted, err := client.GetTenantByName(ctx, sessionName.String())
	if err == nil {
		log.Fatalf("grabbing tenant by its name should have failed. It did not")
	}
	log.Printf("deleted should be a zero value tenant: %#v\n", deleted)

	log.Println("All done with the tenants")
}
