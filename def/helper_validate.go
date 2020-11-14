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
		return NewDefValidateError(
			key,
			fmt.Sprintf("must be one of: %s", strings.Join(slice, ", ")),
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
		return NewDefValidateError(
			key,
			fmt.Sprintf("must be one of: %s", outList),
		)
	}
	return nil
}
