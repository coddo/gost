package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString returns a hashed string and an error
func HashString(data string) (string, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// HashBytes returns a hashed byte array and an error
func HashBytes(data []byte) ([]byte, error) {
	key, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// MatchBytes returns true if the hash matches the password
func MatchBytes(hash, data []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, data)
	if err == nil {
		return true
	}

	return false
}

// MatchString returns true if the hash matches the password
func MatchString(hash, data string) bool {
	return MatchBytes([]byte(hash), []byte(data))
}
