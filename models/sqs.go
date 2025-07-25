package models

type SMSRelayRequest struct {
	Device      Device      `json:"device"`       // Device details
	DeviceName  string      `json:"device_name"`  // Name of the device
	PhoneNumber PhoneNumber `json:"phone_number"` // Phone number details
	SMS         SMS         `json:"sms"`          // SMS message details
}
