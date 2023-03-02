package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/suborbital/se2-go"
)

func main() {
	token, ok := os.LookupEnv("EVERYTHING_TOKEN")
	if !ok {
		log.Fatalf("There is no token set. Get an environment access key from the se2 admin area, and set it to the EVERYTHING_TOKEN env var.")
	}

	// Set up the client to point at staging with a valid access token.
	client, err := se2.NewClient2(se2.HostStaging, token)
	if err != nil {
		log.Fatalf("encountered new client error: %s", err.Error())
	}

	// Create a context for the requests.
	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	// List existing tenants. This should be empty.
	tenants, err := client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants: %s", err.Error())
	}

	fmt.Printf("list tenants is\n%#v\n\n", tenants)

	// Create a new fake name for this run.
	sessionName := ulid.Make()
	log.Printf("\nSession name to start: %s\n\n", sessionName)

	// Create a new tenant by this name.
	newTenant, err := client.CreateTenantByName(ctx, sessionName.String(), "everythingDescription")
	if err != nil {
		log.Fatalf("create tenant by name errorred: %s", err.Error())
	}

	log.Printf("created tenant is\n%#v\n\n", newTenant)

	// List tenants again, this should have the new tenant in it.
	tenants, err = client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants number 2: %s", err.Error())
	}
	log.Printf("list tenants again\n%#v\n\n", tenants)

	// Fetch the individual tenant by its name.
	tenant, err := client.GetTenantByName(ctx, sessionName.String())
	if err != nil {
		log.Fatalf("get tenant by name: %s", err.Error())
	}
	log.Printf("the fetched tenant is\n%#v\n\n", tenant)

	// Update the tenant by its name to use a different name and description.
	newName := ulid.Make()
	log.Printf("trying to change name from %s -> %s\n", sessionName, newName)
	updated, err := client.UpdateTenantByName(ctx, sessionName.String(), "newDescription")
	if err != nil {
		log.Fatalf("update tenant errored: %s", err.Error())
	}
	log.Printf("updated tenant is this:\n%#v\n\n", updated)

	// Try to get tenant by its new name. It should fail.
	tenantDidNotChange, err := client.GetTenantByName(ctx, newName.String())
	if err == nil {
		log.Fatalf("get tenant by new name after having patched the tenant returns nil error. It should not\n%#v\n", tenantDidNotChange)
	}
	log.Printf("getting tenant by its new name resulted in this error: %s\n"+
		"the new tenant should be a zero value: %#v\n", err.Error(), tenantDidNotChange)

	tenantStillExists, err := client.GetTenantByName(ctx, sessionName.String())
	if err != nil {
		log.Fatalf("getting tenant by name, after the update, failed: %s", err.Error())
	}
	log.Printf("we still have the tenant: %#v\n", tenantStillExists)

	// Delete tenant by name
	err = client.DeleteTenantByName(ctx, sessionName.String())
	if err != nil {
		log.Fatalf("deleting tenant by its new name should not have failed. It did: %s", err.Error())
	}

	// List tenants should be empty after this
	listTenants, err := client.ListTenants(ctx)
	if err != nil {
		log.Fatalf("list tenants after deletion should have succeeded. It failed with %s", err.Error())
	}
	log.Printf("this is the list tenants (length: %d)\n%#v\n", len(listTenants.Tenants), listTenants)

	// Getting tenant by its new name should fail
	deleted, err := client.GetTenantByName(ctx, newName.String())
	if err == nil {
		log.Fatalf("grabbing tenant by its name should have failed. It did not")
	}
	log.Printf("deleted should be a zero value tenant: %#v\n", deleted)

	log.Println("All done")
}
