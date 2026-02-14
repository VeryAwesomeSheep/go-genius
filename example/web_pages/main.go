package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/VeryAwesomeSheep/go-genius"
)

func main() {
	token := "your-client-access-token"
	if token == "" {
		log.Fatal("No user token present")
	}

	// Create new client
	client, err := genius.NewClient(token)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Request WebPage data
	webPage, _, err := client.WebPages.Get(ctx, &genius.WebPagesOptions{RawAnnotatableURL: "your-url"})
	if err != nil {
		log.Fatal(err)
	}

	// Access individual pieces of information
	fmt.Println(webPage.Title)

	// Print entirety of requested data
	webPageData, _ := json.MarshalIndent(webPage, "", "  ")
	fmt.Println(string(webPageData))
}
