package models

type SMS struct {
	ID string `json:"id"` // UUID of the SMS message

	From string `json:"from"` // Phone number ID of the sender, in E.164 format
	Body string `json:"body"` // Content of the SMS message, can be plaintext or encrypted

	PhoneNumberID string `json:"phone_number_id"` // ID of the phone number associated with this SMS

	ReceivedAt string `json:"received_at,omitempty"` // Timestamp of when the SMS was received by the device
	CreatedAt  string `json:"created_at,omitempty"`  // Timestamp of when the SMS entry was created in the database
}
