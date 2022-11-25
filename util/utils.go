package util

import (
)

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
