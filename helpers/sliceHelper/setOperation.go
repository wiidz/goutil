package sliceHelper

// Difference 获取s1比s2差异的部分
func Difference(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var diff []interface{}

	// 将s2中的元素添加到map中
	for _, item := range s2 {
		m[item] = true
	}

	// 遍历s1，检查元素是否在s2的map中
	for _, item := range s1 {
		if _, found := m[item]; !found {
			// 如果元素不在map中，则添加到差集结果中
			diff = append(diff, item)
		}
	}

	return diff
}

// DifferenceUint64 获取s1比s2差异的部分
func DifferenceUint64(s1, s2 []uint64) []uint64 {
	m := make(map[interface{}]bool)
	var diff []uint64

	// 将s2中的元素添加到map中
	for _, item := range s2 {
		m[item] = true
	}

	// 遍历s1，检查元素是否在s2的map中
	for _, item := range s1 {
		if _, found := m[item]; !found {
			// 如果元素不在map中，则添加到差集结果中
			diff = append(diff, item)
		}
	}

	return diff
}

// Intersection 获取s1和s2的交集
func Intersection(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var intersection []interface{}

	// 遍历s1并将元素添加到map中
	for _, item := range s1 {
		m[item] = true
	}

	// 遍历s2，检查元素是否在s1的map中
	for _, item := range s2 {
		// 使用反射来比较元素
		if _, ok := m[item]; ok {
			// 如果元素在map中，则添加到交集结果中
			intersection = append(intersection, item)
			// 可以选择从map中删除该元素，以避免重复结果
			delete(m, item)
		}
	}

	return intersection
}

// Union 获取s1和s2的并集
func Union(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var union []interface{}

	// 将s1中的所有元素添加到map中
	for _, item := range s1 {
		m[item] = true
	}

	// 遍历s2，将不在map中的元素添加到map和union中
	for _, item := range s2 {
		if _, found := m[item]; !found {
			m[item] = true
			union = append(union, item)
		}
	}

	// 将s1中剩余的元素添加到union中
	for item := range m {
		union = append(union, item)
	}

	return union
}
