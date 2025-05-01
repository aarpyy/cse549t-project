package mergesort

const (
	basecase           = 4096
	maxDepth           = 4
	k                  = 8
	basecaseCacheAware = 8196
	basecaseMerge      = 8196
)

func Merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	// Compare elements from both halves and add the smaller one
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Add any remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func MergeC(left, right, C []int) {
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			C[k] = left[i]
			i++
		} else {
			C[k] = right[j]
			j++
		}
		k++
	}
	for i < len(left) {
		C[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		C[k] = right[j]
		j++
		k++
	}
}
