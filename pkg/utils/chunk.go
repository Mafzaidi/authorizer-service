package utils

func Chunk[T any](items []T, size int) [][]T {
	if size <= 0 {
		panic("chunk size must be greater than 0")
	}

	var result [][]T
	for size < len(items) {
		result = append(result, items[:size])
		items = items[size:]
	}
	result = append(result, items)
	return result
}
