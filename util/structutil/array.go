package structutil

func UniqArray[T comparable](array []T) []T {
	m := make(map[T]bool)
	uniq := []T{}
	for _, a := range array {
		if !m[a] {
			m[a] = true
			uniq = append(uniq, a)
		}
	}
	return uniq
}
