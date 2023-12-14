package validations

import (
	"car_demo/helper"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
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

func ValidErr(err []*validation.Error) []string {
	message := make([]string, 0, len(err))
	for i := range err {

		message = append(message, err[i].Message)
	}
	return message
}

func GetTag(err []*validation.Error) []string {
	tags := make([]string, 0, len(err))

	for i := range err {
		tags = append(tags, err[i].Name)
	}
	return tags
}

func ValidationErrorResponse(c beego.Controller, err []*validation.Error) []string {
	errs := make([]string, 0, len(err))
	Tags := GetTag(err)
	for i := range Tags {
		var errResponse string
		switch Tags[i] {
		case "Required", "Min", "Max", "Range", "MinSize", "MaxSize", "Length", "Match", "NotMatch":
			errResponse = fmt.Sprintf(helper.LanguageTranslate(c, "validation."+Tags[i]), err[i].LimitValue)
		default:
			errResponse = helper.LanguageTranslate(c, "validation."+Tags[i])
		}
		errs = append(errs, errResponse)
	}
	return errs
}
