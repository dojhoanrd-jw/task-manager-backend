package users

// UpdateRoleRequest represents the request body for changing a user's role
type UpdateRoleRequest struct {
	Role string `json:"role"`
}
