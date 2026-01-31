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

	// Request Song data
	song, _, err := client.Songs.Get(ctx, song-id)
	if err != nil {
		log.Fatal(err)
	}

	// Access individual pieces of information
	fmt.Println(song.Title)
	fmt.Println(song.PrimaryArtists[0].Name)
	fmt.Println(song.Stats.Hot)

	// Access pointer type variables
	if song.ReleaseDate != nil {
		fmt.Println(*song.ReleaseDate)
	}

	// Print entirety of requested data
	songData, _ := json.MarshalIndent(song, "", "  ")
	fmt.Println(string(songData))
}
