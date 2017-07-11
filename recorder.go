package errors

import (
	"errors"
	"fmt"
)

type Generator struct {
	Name string
	Fields
}

func NewGenerator(name string) *Generator {
	return &Generator{Name: name}
}

func (r *Generator) New(message string) *Error {
	return genError(errors.New(message), 0)
}

func (r *Generator) Errorf(format string, args ...interface{}) *Error {
	return genError(fmt.Errorf(format, args...), 0)
}

func (r *Generator) Wrap(err error) *Error {
	return genError(err, 0)
}
