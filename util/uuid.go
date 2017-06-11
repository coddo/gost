package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
)

const uuidRegexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"

// A UUID representation compliant with specification in
// RFC 4122 document.
type UUID [16]byte

// IsValidUUID tells if a given UUID string is valid
func IsValidUUID(s string) bool {
	var re = regexp.MustCompile(uuidRegexPattern)

	md := re.FindStringSubmatch(s)
	if md == nil {
		return false
	}

	hash := md[2] + md[3] + md[4] + md[5] + md[6]

	_, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	return true
}

// GenerateUUID creates a new UUID string
func GenerateUUID() (string, error) {
	u := new(UUID)

	// Set all bits to randomly (or pseudo-randomly) chosen values.
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x40) & 0x7F

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
