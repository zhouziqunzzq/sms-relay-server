package models

type User struct {
	ID string `json:"id"` // UUID of the user

	Username string `json:"username"` // Unique username of the user
	Password string `json:"-"`        // Bcrypt hashed password of the user, not returned in API responses

	Name  string `json:"name,omitempty"`  // Displayed name of the user
	Email string `json:"email,omitempty"` // Email address of the user

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the user was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the user was last updated
}
