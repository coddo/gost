package sliceutil

// AreIntersected returns true if at least one item in the first slice is found in the second slice.
// This method assumes the the slices contains unique items inside themselves.
// If the second slice is empty, then this func returns true.
func AreIntersected(source []string, target []string) bool {
	if source == nil || target == nil {
		return false
	}

	var sourceLen, targetLen = len(source), len(target)

	if targetLen == 0 {
		return true
	}

	if sourceLen == 0 {
		return false
	}

	// Create a map of string items
	var itemsMap = make(map[string]int, sourceLen+targetLen)

	// Add the items from the first slice into the map
	for i := 0; i < sourceLen; i++ {
		itemsMap[source[i]] = 1
	}

	// Check the items from the map with the items from the second slice
	for i := 0; i < targetLen; i++ {
		var value, itemWasFound = itemsMap[target[i]]

		if itemWasFound && value == 1 {
			return true
		}
	}

	// No intersection was found between the slices
	return false
}
