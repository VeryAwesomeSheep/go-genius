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

	// Request Annotation data
	annotation, referent, _, err := client.Annotations.Get(ctx, annotation-id, &genius.AnnotationsOptions{TextFormatOptions: genius.TextFormatOptions{TextFormat: genius.FormatPlain}})
	if err != nil {
		log.Fatal(err)
	}

	// Access individual pieces of information
	fmt.Println(annotation.Body.Plain)
	fmt.Println(referent.AnnotatorID)

	// Print entirety of requested data
	annotationData, _ := json.MarshalIndent(annotation, "", "  ")
	fmt.Println(string(annotationData))
}
