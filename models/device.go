package models

type Device struct {
	ID string `json:"id"` // UUID of the device

	Name  string `json:"name,omitempty"` // Displayed name of the device
	Token string `json:"-"`              // Bcrypt hashed token for the device, not returned in API responses

	PhoneNumberIDs []string `json:"phone_number_ids"` // List of phone number IDs associated with the device

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the device was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the device was last
}
