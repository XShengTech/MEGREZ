package utils

import (
	"regexp"
)

func EmailFormat(email string) bool {
	pattern := `\w[-\w.+]*@([-A-Za-z0-9]+\.)+[A-Za-z]{2,14}`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return match
}
