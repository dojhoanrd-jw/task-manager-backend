package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// NewClient creates a new Firestore client
func NewClient(projectID string) *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	log.Println("Firestore client connected successfully")
	return client
}
