package models

type SMSRelayRequest struct {
	DeviceID string `json:"device_id"` // UUID of the device making the request
	SMS      SMS    `json:"sms"`       // SMS message details
}
