package utils

import "reflect"

func Contains[T Ordered](slice *[]T, v T) bool {
	for _, e := range *slice {
		if e == v {
			return true
		}
	}
	return false
}

func Chunk[T any](slice []T, chunkSize int) [][]T {
	numChunks := (len(slice) + chunkSize - 1) / chunkSize
	chunks := make([][]T, numChunks)
	for i := range chunks {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks[i] = make([]T, end-start)
		copy(chunks[i], slice[start:end])
	}
	return chunks
}

// HashMap is a helper function to convert a slice of structs into a map
// where the key is the value of the field with the name keyField.
//
// The value of the map is the struct itself.
func HashMap[T any](slice *[]T, keyField string) map[interface{}]T {
	m := make(map[interface{}]T)
	v := reflect.ValueOf((*slice)[0])
	for _, e := range *slice {
		m[v.FieldByName(keyField).Interface()] = e
	}
	return m
}
