package models

type ForwardDestinations struct {
	Email EmailForwardDestination `json:"email"` // Email destination for forwarding messages
}

type EmailForwardDestination struct {
	Email string `json:"email"` // Email address to forward messages to
}

func (efd *EmailForwardDestination) IsEmpty() bool {
	return efd.Email == ""
}
