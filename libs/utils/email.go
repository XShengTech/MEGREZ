package utils

import (
	"regexp"
)

func EmailFormat(email string) bool {
	pattern := `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return match
}
