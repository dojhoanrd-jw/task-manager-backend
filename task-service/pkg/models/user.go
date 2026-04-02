package models

import "time"

// Role defines user permission levels
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

// User represents a user entity
type User struct {
	ID        string    `firestore:"-" json:"id"`
	Name      string    `firestore:"name" json:"name"`
	Email     string    `firestore:"email" json:"email"`
	Password  string    `firestore:"password" json:"-"`
	Role      Role      `firestore:"role" json:"role"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}
