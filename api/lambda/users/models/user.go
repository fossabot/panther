package models

// Group is a struct for Panther Group containing employees.
type Group struct {
	Description *string `json:"description"`
	Name        *string `json:"name"`
}

// User is a struct describing a Panther User.
type User struct {
	CreatedAt   *int64  `json:"createdAt"`
	Email       *string `json:"email"`
	FamilyName  *string `json:"familyName"`
	GivenName   *string `json:"givenName"`
	ID          *string `json:"id"`
	PhoneNumber *string `json:"phoneNumber"`
	Role        *string `json:"role"` // Roles are group name
	Status      *string `json:"status"`
}
