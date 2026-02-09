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

	// Request Artist data
	artist, _, err := client.Artists.Get(ctx, artist-id)
	if err != nil {
		log.Fatal(err)
	}

	// Request Artist's Songs data
	artistSongs, _, err := client.Artists.GetSongs(ctx, artist-id,
		&genius.ArtistSongsOptions{Sort: genius.SortPopularity, PagingOptions: genius.PagingOptions{PerPage: 10, Page: 1}})

	// Access individual pieces of information
	fmt.Println(artist.Name)
	fmt.Println(artist.FollowersCount)

	// Access pointer type variables
	if artist.SocialLinks != nil && artist.SocialLinks.Instagram != nil {
		fmt.Println(*artist.SocialLinks.Instagram)
	}

	// Access array type variables
	if len(artist.AlternateNames) > 0 {
		fmt.Println(artist.AlternateNames[0])
	}

	if len(artistSongs) > 0 {
		for _, song := range artistSongs {
			fmt.Println(song.Title)
		}
	}

	// Print entirety of requested data
	artistData, _ := json.MarshalIndent(artist, "", "  ")
	fmt.Println(string(artistData))
}
