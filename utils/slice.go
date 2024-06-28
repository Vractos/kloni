package utils

import (
	"fmt"
	"reflect"
)

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

// HashMap creates a map using the specified key field from the elements of the given slice.
// The key field must be a valid field name in the struct elements of the slice.
// The function returns the created map and an error, if any.
func HashMap[T any](slice *[]T, keyField string) (map[interface{}]T, error) {
	m := make(map[interface{}]T)

	if slice == nil || len(*slice) == 0 {
		return m, nil
	}

	v := reflect.ValueOf((*slice)[0])
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("slice elements must be structs")
	}

	field := v.FieldByName(keyField)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found in struct", keyField)
	}

	for _, e := range *slice {
		v := reflect.ValueOf(e)
		key := v.FieldByName(keyField).Interface()
		m[key] = e
	}

	return m, nil
}
