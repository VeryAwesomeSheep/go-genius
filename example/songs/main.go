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

	// Request Song data
	songID := 1234
	song, _, err := client.Songs.Get(ctx, songID)
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
