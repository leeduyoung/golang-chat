package array

// Contains 배열에서 아이템이 존재하는지 확인하는 함수
func Contains[T int | string](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// RemoveDuplcateItem 배열에서 중복된 아이템을 제거하는 함수
func RemoveDuplcateItem[T int | string](items []T) []T {
	result := []T{}

	m := make(map[T]bool)
	for _, item := range items {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}
