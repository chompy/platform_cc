package def

import (
	"fmt"
	"strings"
)

func sliceContainsString(slice []string, s string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}

func validateMustContainOne(slice []string, s string, key string) error {
	if !sliceContainsString(slice, s) {
		return NewValidateError(
			key,
			fmt.Sprintf("must be one of: %s got: %s", strings.Join(slice, ", "), s),
		)
	}
	return nil
}

func sliceContainsInt(slice []int, n int) bool {
	for _, val := range slice {
		if val == n {
			return true
		}
	}
	return false
}

func validateMustContainOneInt(slice []int, n int, key string) error {
	if !sliceContainsInt(slice, n) {
		outList := ""
		for _, val := range slice {
			outList += fmt.Sprintf("%d, ", val)
		}
		return NewValidateError(
			key,
			fmt.Sprintf("must be one of: %s got: %d", strings.TrimRight(outList, ","), n),
		)
	}
	return nil
}
