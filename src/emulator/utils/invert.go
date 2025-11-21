package utils

// 🟦 Invert a map [k, v] -> [v, k]

func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	inv := make(map[V]K, len(m))
	for k, v := range m {
		inv[v] = k
	}
	return inv
}
