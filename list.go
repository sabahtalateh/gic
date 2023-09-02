package gic

// List groups multiple Get function
// Useful to store array of interfaces
func List[T any](cc ...T) []T {
	return cc
}
