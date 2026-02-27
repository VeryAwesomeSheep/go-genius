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

	// Request User data
	userID := 1234
	user, _, err := client.Users.Get(ctx, userID, &genius.UserOptions{TextFormatOptions: genius.TextFormatOptions{TextFormat: genius.FormatPlain}})
	if err != nil {
		log.Fatal(err)
	}

	// Access individual pieces of information
	fmt.Println(user.Name)
	fmt.Println(user.FollowersCount)
	fmt.Println(user.RoleForDisplay)

	// Access pointer type variables
	if user.AboutMe.Plain != nil {
		fmt.Println(*user.AboutMe.Plain)
	}

	// Print entirety of requested data
	userData, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(userData))
}
