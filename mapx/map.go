package mapx

// Keys 返回 map 里面的所有的 key。
// 需要注意：这些 key 的顺序是随机。
func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
