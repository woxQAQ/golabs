package rbtree

// If is a single line if/else statement.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func Ternary[T any](condition bool, ifOutput, elseOutput T) T {
	if condition {
		return ifOutput
	}

	return elseOutput
}
