package models

import (
	"encoding/json"
	"time"
)

// Task represents a task entity
type Task struct {
	ID          string    `firestore:"-" json:"id"`
	Title       string    `firestore:"title" json:"title"`
	Description string    `firestore:"description" json:"description"`
	Completed   bool      `firestore:"completed" json:"completed"`
	ProjectID   string    `firestore:"projectId" json:"projectId"`
	AssignedTo  string    `firestore:"assignedTo" json:"assignedTo"`
	CreatedAt   time.Time `firestore:"createdAt" json:"-"`
}

// MarshalJSON customizes the JSON output to include formatted date and time
func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task
	return json.Marshal(&struct {
		Alias
		CreatedDate string `json:"createdDate"`
		CreatedTime string `json:"createdTime"`
	}{
		Alias:       (Alias)(t),
		CreatedDate: t.CreatedAt.Format("02-Jan-06"),
		CreatedTime: t.CreatedAt.Format("15:04"),
	})
}
