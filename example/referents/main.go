package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	genius "github.com/VeryAwesomeSheep/go-genius"
)

func main() {
	// Create new client
	config := &genius.Config{
		Token: "your-client-access-token",
	}

	client, err := genius.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Request Referents data
	webPageID := 1234
	createdByID := 5678
	referents, _, err := client.Referents.Get(ctx, &genius.ReferentsOptions{WebPageID: webPageID, CreatedByID: createdByID, TextFormatOptions: genius.TextFormatOptions{TextFormat: genius.FormatPlain}})
	if err != nil {
		log.Fatal(err)
	}

	if len(referents) > 0 {
		// Access individual pieces of information
		fmt.Println(referents[0].Classification)

		// Access pointer type variables
		if referents[0].SongID != nil {
			fmt.Println(*referents[0].SongID)
		}
	}

	// Print entirety of requested data
	referentsData, _ := json.MarshalIndent(referents, "", "  ")
	fmt.Println(string(referentsData))
}
