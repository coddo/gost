package util

import "regexp"

const emailRegexPattern = "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$"

// IsValidEmail checks if an email is in the correct format
func IsValidEmail(email string) bool {
	var re = regexp.MustCompile(emailRegexPattern)

	return re.MatchString(email)
}
