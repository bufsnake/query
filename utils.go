package query

import (
	"strings"
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func InArr(arr []string, data string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.ToLower(arr[i]) == data {
			return true
		}
	}
	return false
}
