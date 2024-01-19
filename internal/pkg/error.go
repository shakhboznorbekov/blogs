package pkg

import (
	"errors"
	"fmt"
)

type Error struct {
	Err    error        `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
	Status int          `json:"status"`
}

type FieldError struct {
	Err   error  `json:"error"`
	Field string `json:"field"`
}

func WrapError(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s : %s", msg, err.Error()))
}
