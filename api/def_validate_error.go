package api

import "fmt"

// NewDefValidateError - creates a new defValidateError
func NewDefValidateError(key string, msg string) error {
	return &defValidateError{
		key: key,
		msg: msg,
	}
}

// defValidateError - implementation of error that defines a def validation error
type defValidateError struct {
	key string
	msg string
}

func (d *defValidateError) Error() string {
	if d.msg != "" {
		return fmt.Sprintf("validation error at '%s', %s", d.key, d.msg)
	}
	return fmt.Sprintf("validation error at '%s'", d.key)
}
