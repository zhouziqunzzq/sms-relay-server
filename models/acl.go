package models

// ACL (Access Control List) defines permissions for users, devices, and phone numbers:
// - If a user has an ACL entry for a device, they have access to SMS of all phone numbers
//   associated with that device.
// - If a user has an ACL entry for a phone number, they have access to SMS of that phone number.
type ACL struct {
	ID string `json:"id"` // UUID of the ACL entry

	UserID        string `json:"user_id"`         // ID of the user this ACL entry applies to
	DeviceID      string `json:"device_id"`       // ID of the device this ACL entry applies to
	PhoneNumberID string `json:"phone_number_id"` // ID of the phone number this ACL entry applies to

	CreatedAt string `json:"created_at,omitempty"` // Timestamp of when the ACL entry was created
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of when the ACL entry was last updated
}
