package def

import "fmt"

// NewValidateError creates a new ValidateError.
func NewValidateError(key string, msg string) error {
	return &ValidateError{
		key: key,
		msg: msg,
	}
}

// ValidateError extends error by providing specific details about a def validation error.
type ValidateError struct {
	key string
	msg string
}

func (d *ValidateError) Error() string {
	if d.msg != "" {
		return fmt.Sprintf("validation error at '%s', %s", d.key, d.msg)
	}
	return fmt.Sprintf("validation error at '%s'", d.key)
}
