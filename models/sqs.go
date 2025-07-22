package models

type SMSRelayRequest struct {
	Device      Device      `json:"device"`       // Device details
	PhoneNumber PhoneNumber `json:"phone_number"` // Phone number details
	SMS         SMS         `json:"sms"`          // SMS message details
}
