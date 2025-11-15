package utils

// 🟦 We miss ternary expressions in Typescript etc

func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
