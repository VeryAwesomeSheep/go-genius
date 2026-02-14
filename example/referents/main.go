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

	// Request Referents data
	referents, _, err := client.Referents.Get(ctx, &genius.ReferentsOptions{WebPageID: webpage-id, CreatedByID: id, TextFormatOptions: genius.TextFormatOptions{TextFormat: genius.FormatPlain}})
	if err != nil {
		log.Fatal(err)
	}

	// Access individual pieces of information
	if len(referents) > 0 {
		fmt.Println(referents[0].Classification)
	}

	// Print entirety of requested data
	referentsData, _ := json.MarshalIndent(referents, "", "  ")
	fmt.Println(string(referentsData))
}
