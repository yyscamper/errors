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

func (r *Generator) WithFields(fields Fields) *Generator {
	if r.Fields == nil {
		r.Fields = make(map[string]interface{})
	}
	for k, v := range fields {
		r.Fields[k] = v
	}
	return r
}

func (r *Generator) WithField(key string, val interface{}) *Generator {
	return r.WithFields(Fields{key: val})
}

func (r *Generator) New(message string) *Error {
	return genError(errors.New(message), 0).WithFields(r.Fields).WithName(r.Name)
}

func (r *Generator) Errorf(format string, args ...interface{}) *Error {
	return genError(fmt.Errorf(format, args...), 0).WithFields(r.Fields).WithName(r.Name)
}

func (r *Generator) Wrap(err error) *Error {
	switch err.(type) {
	case *Error:
		return err.(*Error)
	default:
		return genError(err, 0).WithFields(r.Fields).WithName(r.Name)
	}
}
