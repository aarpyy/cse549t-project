package mergesort

func CheckSorted(original, sorted []int) bool {
	// Check length
	if len(original) != len(sorted) {
		return false
	}

	// Check if sorted
	for i := 1; i < len(sorted); i++ {
		if sorted[i-1] > sorted[i] {
			return false
		}
	}

	return true
}
