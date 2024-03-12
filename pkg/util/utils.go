package util

import "reflect"

func Unique(input []string) []string {
	result := make([]string, 0, len(input))
	values := make(map[string]bool)

	for _, val := range input {
		if _, ok := values[val]; !ok {
			values[val] = true
			result = append(result, val)
		}
	}
	return result
}

func NonZeroEntries[T any](s []T) int {
	var count int
	for _, v := range s {
		if reflect.ValueOf(v).IsValid() {
			count++
		}
	}

	return count
}
