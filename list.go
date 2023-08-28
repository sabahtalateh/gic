package gic

// List groups multiple Get function
func List[T any](cc ...T) []T {
	return cc
}
