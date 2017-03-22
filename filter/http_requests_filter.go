package filter

// CheckNotNull verifies if the given []byte data is nil or represents the word "null".
// In the above stated cases, the method returns false
func CheckNotNull(data []byte) bool {
	if data == nil {
		return false
	}

	if string(data) == "null" {
		return false
	}

	return true
}
