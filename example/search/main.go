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

	// Request Search data
	hits, _, err := client.Search.Get(ctx, "your-query", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, hit := range hits {
		// Access individual pieces of information
		fmt.Println(hit.Title)

		// Access pointer type variables
		if hit.ReleaseDateComponents.Year != nil {
			fmt.Println(*hit.ReleaseDateComponents.Year)
		}

		// Access array type variables
		if len(hit.PrimaryArtists) > 0 {
			fmt.Println(hit.PrimaryArtists[0].Name)
		}
	}

	// Print entirety of requested data
	hitsData, _ := json.MarshalIndent(hits, "", "  ")
	fmt.Println(string(hitsData))
}
