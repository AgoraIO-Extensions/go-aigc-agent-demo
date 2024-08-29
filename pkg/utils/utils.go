package utils

// MM 合并两个map
func MM[K comparable, V any](map1, map2 map[K]V) map[K]V {
	mergedMap := make(map[K]V)
	for k, v := range map1 {
		mergedMap[k] = v
	}
	for k, v := range map2 {
		mergedMap[k] = v
	}
	return mergedMap
}
