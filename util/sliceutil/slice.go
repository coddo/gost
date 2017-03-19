package sliceutil

// AreIntersected returns true if at least one item in the first slice is found in the second slice
// This method assumes the the slices contains unique items inside themselves
func AreIntersected(firstSlice []string, secondSlice []string) bool {
	var firstSliceLen, secondSliceLen = len(firstSlice), len(secondSlice)

	if firstSlice == nil || secondSlice == nil || firstSliceLen == 0 || secondSliceLen == 0 {
		return false
	}

	// Create a map of string items
	var itemsMap = make(map[string]int, firstSliceLen+secondSliceLen)

	// Add the items from the first slice into the map
	for i := 0; i < firstSliceLen; i++ {
		itemsMap[firstSlice[i]] = 1
	}

	// Check the items from the map with the items from the second slice
	for i := 0; i < secondSliceLen; i++ {
		var value, itemWasFound = itemsMap[secondSlice[i]]

		if itemWasFound && value == 1 {
			return true
		}
	}

	// No intersection was found between the slices
	return false
}
