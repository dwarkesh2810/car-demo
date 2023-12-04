package models

type OTP struct {
	PhoneNumber string `json:"phone_number"`
}

type VerifyData struct {
	User string `json:"user,omitempty"`
	Code string `json:"code,omitempty"`
}