package utils

func Contains[T Ordered](slice *[]T, v T) bool {
	for _, e := range *slice {
		if e == v {
			return true
		}
	}
	return false
}
