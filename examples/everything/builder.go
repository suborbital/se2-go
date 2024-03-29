package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/suborbital/se2-go"
)

func builder(client *se2.Client) {
	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	const (
		pluginName = "everythingbagel"
		namespace  = "everythingns"
	)

	m := ulid.Make()

	printHeader("creating tenant for session")
	sessionTenant, err := client.CreateTenant(ctx, m.String(), "description for session tenant")
	if err != nil {
		log.Fatalf("creating tenant for session failed with %s", err.Error())
	}

	printHeader(fmt.Sprintf("creating session for tenant '%s', namespace 'everythingns' and plugin 'everythingbagel'", sessionTenant.Name))
	s, err := client.CreateSession(ctx, sessionTenant.Name, namespace, pluginName)
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

	printHeader("Set the template to 'javascript'")
	_, err = client.CreatePluginDraft(ctx, "javascript", s)
	if err != nil {
		log.Fatalf("setting draft to javascript failed with %s", err.Error())
	}

	printHeader("getting plugin drafts saved on session")

	pd, err := client.GetPluginDraft(ctx, s)
	if err != nil {
		log.Fatalf("getting plugin drafts failed: %s", err.Error())
	}

	log.Printf("plugin draft is\n%#v\n", pd)

	printHeader("building a js plugin synchronously")

	newPlugin := []byte(`import { log } from "@suborbital/plugin";

export const run = (input) => {
    let message = "Hello, " + input;

    message = message.split("").reverse().join("")

    log.info(message);

    return message;
};`)

	// Synchronous builds take a while, we need a bigger timeout than what's on the previous context.
	buildCtx, buildCxl := context.WithTimeout(context.Background(), 5*time.Minute)
	defer buildCxl()

	built, err := client.BuildPlugin(buildCtx, newPlugin, s)
	if err != nil {
		log.Fatalf("building plugin failed with error:\n%s\n", err.Error())
	}

	fmt.Printf("this is the response to having built the plugin:\n\n%#v\n", built)

	printHeader("Grab the current draft again, it should be the new plugin code")

	draft, err := client.GetPluginDraft(buildCtx, s)
	if err != nil {
		log.Fatalf("getting plugin draft after a build failed with: %s", err.Error())
	}

	fmt.Printf("Current draft after a build is:\n\n%#v\n\n", draft)

	printHeader("running test for the draft")

	testResult, err := client.TestPluginDraft(buildCtx, []byte(`capt picard`), s)
	if err != nil {
		log.Fatalf("running the test on the plugin draft failed with %s", err.Error())
	}

	fmt.Printf("response from the test plugin draft is\n\n%#v\n", testResult)

	printHeader("promoting current draft to be live")

	promotionResult, err := client.PromotePluginDraft(buildCtx, s)
	if err != nil {
		log.Fatalf("promoting draft plugin failed with %s", err.Error())
	}

	fmt.Printf("this is the result of the promotion:\n%#v\n\n", promotionResult)

	printHeader("List the plugins for tenant")

	plugins, err := client.GetPlugins(buildCtx, sessionTenant.Name)
	if err != nil {
		log.Fatalf("listing plugins failed with %s", err.Error())
	}

	fmt.Printf("tenant %s has the following plugins:\n\n%#v\n\n", sessionTenant.Name, plugins)

	printHeader("sleep for 8 seconds to wait until module gets propagated")
	time.Sleep(8 * time.Second)

	printHeader("execution")
	exec, err := client.Exec(buildCtx, []byte(`uh hi`), sessionTenant.Name, namespace, pluginName)
	if err != nil {
		log.Fatalf("client.Exec: %s", err.Error())
	}

	fmt.Printf("exec result is\n%#v\nor in string: %s\n", exec, string(exec))

	printHeader("deleting tenant of the session")

	err = client.DeleteTenantByName(buildCtx, sessionTenant.Name)
	if err != nil {
		log.Fatalf("could not delete tenant by name %s: %s", sessionTenant.Name, err.Error())
	}
}
