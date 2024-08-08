package algext

// ZeroOf returns the zero value of the type T
func ZeroOf[T any]() T {
	var zero T
	return zero
}
