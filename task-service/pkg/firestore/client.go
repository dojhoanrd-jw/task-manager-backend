package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/logger"
)

// NewClient creates a new Firestore client
func NewClient(projectID string) *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("Failed to create Firestore client: " + err.Error())
		os.Exit(1)
	}

	logger.Info("Firestore client connected successfully")
	return client
}
