package slicex

// Column get column from slice
func Column[K, T any](slice []T, column func(item T) K) []K {
	r := make([]K, len(slice))
	for i := range slice {
		r[i] = column(slice[i])
	}
	return r
}
