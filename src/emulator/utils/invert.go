package utils

func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	inv := make(map[V]K, len(m))
	for k, v := range m {
		inv[v] = k
	}
	return inv
}

func InvertMulti[K comparable, V comparable](m map[K]V) map[V][]K {
	inv := make(map[V][]K)
	for k, v := range m {
		inv[v] = append(inv[v], k)
	}
	return inv
}
