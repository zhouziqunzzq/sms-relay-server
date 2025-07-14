package models

type PhoneNumber struct {
	ID string `json:"id"` // UUID of the phone number

	Number string `json:"number"`         // Full phone number in E.164 format
	Name   string `json:"name,omitempty"` // Displayed name of the phone number

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the phone number was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the phone number was last updated
}
