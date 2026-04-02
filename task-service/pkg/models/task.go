package models

import "time"

// Task represents a task entity
type Task struct {
	ID          string    `firestore:"-" json:"id"`
	Title       string    `firestore:"title" json:"title"`
	Description string    `firestore:"description" json:"description"`
	Completed   bool      `firestore:"completed" json:"completed"`
	ProjectID   string    `firestore:"projectId" json:"projectId"`
	AssignedTo  string    `firestore:"assignedTo" json:"assignedTo"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
}
