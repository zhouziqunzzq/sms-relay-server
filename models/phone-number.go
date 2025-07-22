package models

type PhoneNumber struct {
	ID string `json:"id"` // UUID of the phone number

	Number string `json:"number"`         // Full phone number in E.164 format
	Name   string `json:"name,omitempty"` // Displayed name of the phone number

	// ForwardDestinations contains the list of destinations to which SMS messages
	// sent to this phone number should be forwarded.
	ForwardDestinations ForwardDestinations `json:"forward_destinations"`

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the phone number was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the phone number was last updated
}
