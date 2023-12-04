package validation

import "github.com/beego/beego/v2/core/validation"

func init() {
	validation.SetDefaultMessage(map[string]string{
		"Required":     "Must be filled in",
		"Min":          "Minimum allowed value %d",
		"Max":          "Maximum allowed value %d",
		"Range":        "Must be between %d and %d",
		"MinSize":      "Minimum allowed length %d",
		"MaxSize":      "Maximum allowed length %d",
		"Length":       "Length must be %d",
		"Alpha":        "Must consist of letters",
		"Numeric":      "Must consist of numbers",
		"AlphaNumeric": "Must consist of letters or numbers",
		"Match":        "Must match %s",
		"NoMatch":      "Must not match %s",
		"AlphaDash":    "Must consist of letters, numbers or symbols (-_)",
		"Email":        "Must be in correct email format",
		"IP":           "Must be a valid IP address",
		"Base64":       "Must be in correct base64 format",
		"Mobile":       "Must be a valid mobile phone number",
		"Tel":          "Must be a valid phone number",
		"Phone":        "Must be a valid phone or mobile number",
		"ZipCode":      "Must be a valid zip code",
	})
}

type Example struct {
	Id     int
	Name   string `valid:"Required"`            // Name can't be empty or start with Bee
	Age    int    `valid:"Range(1, 140)"`       // 1 <= Age <= 140, only valid in this range
	Email  string `valid:"Email; MaxSize(100)"` // Need to be a valid Email address and no more than 100 characters.
	Mobile string `valid:"Mobile"`              // Must be a valid mobile number
	IP     string `valid:"IP"`                  // Must be a valid IPv4 address
}
