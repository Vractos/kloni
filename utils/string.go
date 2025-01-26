package utils

// GetOrDefault returns the value of the given string pointer if it is not nil,
// otherwise it returns the provided default value.
func GetOrDefault(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return *s
}
