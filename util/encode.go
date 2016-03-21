package util

import "encoding/base64"

// Encode encodes data using the base64 package
func Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

// Decode decodes data using the base64 package
func Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
