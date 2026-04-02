package models

import "time"

// Project represents a project entity
type Project struct {
	ID          string    `firestore:"-" json:"id"`
	Name        string    `firestore:"name" json:"name"`
	Description string    `firestore:"description" json:"description"`
	OwnerID     string    `firestore:"ownerId" json:"ownerId"`
	Members     []string  `firestore:"members" json:"members"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
}
