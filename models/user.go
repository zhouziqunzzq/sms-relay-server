package models

const (
	UserTypeUser   = "USER"   // UserTypeUser represents a regular user account
	UserTypeDevice = "DEVICE" // UserTypeDevice represents a device account
)

type User struct {
	ID string `json:"id"` // UUID of the user

	Username string `json:"username"`       // Unique username of the user
	Password string `json:"-"`              // Bcrypt hashed password of the user, not returned in API responses
	UserType string `json:"user_type"`      // Type of user account (e.g., USER, DEVICE)
	Name     string `json:"name,omitempty"` // Displayed name of the user

	DeviceID string `json:"device_id,omitempty"` // ID of the device associated with this user, if applicable
	Email    string `json:"email,omitempty"`     // Email address of the user

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the user was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the user was last updated
}
