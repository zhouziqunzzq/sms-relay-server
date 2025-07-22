package models

type Device struct {
	ID string `json:"id"` // UUID of the device

	PhoneNumberIDs []string `json:"phone_number_ids"` // List of phone number IDs associated with the device

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the device was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the device was last updated
}
