package components

import (
	"regexp"
)

var FullnameRegexValidation = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

var EmailRegexValidation = regexp.MustCompile(`([a-z0-9\+_\-]+)(\.[a-z0-9\+_\-]+)*@([a-z0-9\-]+\.)+[a-z]{2,6}`).MatchString

var PasswordRegexValidation = regexp.MustCompile(`^[a-zA-Z-0-9]+$`).MatchString

func ValidateFullName(fullname string) bool {
	if FullnameRegexValidation(fullname) == true {
		return true
	} else {
		return false
	}
}

func ValidateEmail(email string) bool {
	if EmailRegexValidation(email) == true {
		return true
	} else {
		return false
	}
}

func ValidatePassword(password string) bool {
	if PasswordRegexValidation(password) == true {
		return true
	} else {
		return false
	}
}
