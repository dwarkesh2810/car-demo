package validations

import (
	"car_demo/helper"
	"fmt"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

// func Init() {
// 	validation.SetDefaultMessage(map[string]string{
// 		"Required":     "Must be filled in",
// 		"Min":          "Minimum allowed value %d",
// 		"Max":          "Maximum allowed value %d",
// 		"Range":        "Must be between %d and %d",
// 		"MinSize":      "Minimum allowed length %d",
// 		"MaxSize":      "Maximum allowed length %d",
// 		"Length":       "Length must be %d",
// 		"Alpha":        "Must consist of letters",
// 		"Numeric":      "Must consist of numbers",
// 		"AlphaNumeric": "Must consist of letters or numbers",
// 		"Match":        "Must match %s",
// 		"NoMatch":      "Must not match %s",
// 		"AlphaDash":    "Must consist of letters, numbers or symbols (-_)",
// 		"Email":        "Must be in correct email format",
// 		"IP":           "Must be a valid IP address",
// 		"Base64":       "Must be in correct base64 format",
// 		"Mobile":       "Must be a valid mobile phone number",
// 		"Tel":          "Must be a valid phone number",
// 		"Phone":        "Must be a valid phone or mobile number",
// 		"ZipCode":      "Must be a valid zip code",
// 	})
// }

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
		case "Min", "Max", "Range", "MinSize", "MaxSize", "Length", "Match", "NotMatch":
			errResponse = fmt.Sprintf("%s :"+helper.LanguageTranslate(c, "validation."+Tags[i]), err[i].Field, err[i].LimitValue)

		case "Required", "Alpha", "Numeric", "AlphaNumeric", "Email", "IP", "AlphaDash":
			errResponse = fmt.Sprintf("%s :"+helper.LanguageTranslate(c, "validation."+Tags[i]), err[i].Field)
			
		default:
			fields := err[i].Key
			keys := strings.Split(fields, ".")
			errResponse = fmt.Sprintf("%s :"+helper.LanguageTranslate(c, "validation."+keys[1]), keys[0])
		}
		errs = append(errs, errResponse)
	}
	return errs
}

func IndianMobile(v *validation.Validation, obj interface{}, key string) {
	value, ok := obj.(string)
	if !ok {
		return
	}
	pattern := `^[6-9][0-9]{9}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(value) {
		v.SetError(key, "Please enter a valid Indian mobile number")
	}
}
