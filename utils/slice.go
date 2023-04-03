package utils

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
