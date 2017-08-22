package errors

import (
	"errors"
	"fmt"
	"io"
)

//Fields records additional fields into error object, inspired by logrus
type Fields map[string]interface{}

// Error is an error that has a message and a stack, but no caller.
type Error struct {
	err error
	*stack
	Fields
	Name      string
	stackSkip int
}

//DefaultStackSkip is to skip the topmost stack strace, as these traces are within this package
const DefaultStackSkip = 4

func genError(err error, stackSkip int) *Error {
	return &Error{
		err:   err,
		stack: callers(DefaultStackSkip + stackSkip),
	}
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) *Error {
	return genError(errors.New(message), 0)
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) *Error {
	return genError(fmt.Errorf(format, args...), 0)
}

//Wrap converts the golang standard error into Error and with stack trace is recorded
func Wrap(err error) *Error {
	switch err.(type) {
	case *Error:
		return err.(*Error)
	default:
		return genError(err, 0)
	}
}

func (e *Error) Error() string {
	if e.err == nil {
		return ""
	}
	if e.Name != "" {
		return e.Name + ": " + e.err.Error()
	}
	return e.err.Error()
}

//Format formats Error, so that you can use the %+v in fmt.Sprintf
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Error())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func mergeFields(forigin, fnew Fields) Fields {
	cap := 0
	if forigin != nil {
		cap += len(forigin)
	}
	if fnew != nil {
		cap += len(fnew)
	}

	result := make(map[string]interface{}, cap)
	if forigin != nil {
		for k, v := range forigin {
			result[k] = v
		}
	}
	if fnew != nil {
		for k, v := range fnew {
			result[k] = v
		}
	}
	return result
}

//WithFields attaches some key-value pairs additional information into the Error instance
func (e *Error) WithFields(fields Fields) *Error {
	return &Error{
		err:    e.err,
		stack:  e.stack,
		Fields: mergeFields(e.Fields, fields),
		Name:   e.Name,
	}
}

//WithField attaches a key-value pair additional information into the Error instance
func (e *Error) WithField(key string, val interface{}) *Error {
	return e.WithFields(Fields{key: val})
}

func stringify(val interface{}) string {
	if k, ok := val.(string); ok {
		return k
	}
	return fmt.Sprintf("%v", val)
}

//With attaches one or multiple key-value paris additional information into the Error instance
//For example:
// e.With(key1, val1, key2, val2, key3, val3)
func (e *Error) With(key string, val interface{}, extras ...interface{}) *Error {
	fields := Fields{}
	fields[key] = val

	n := len(extras)
	if n%2 != 0 {
		n--
	}
	for i := 0; i < n; i += 2 {
		fields[stringify(extras[i])] = extras[i+1]
	}

	//if forget to attach a value for the last key, then automatically adds a `nil` value
	if n < len(extras) {
		fields[stringify(extras[n])] = nil
	}
	return e.WithFields(fields)
}

//WithName sets a name for error
func (e *Error) WithName(name string) *Error {
	e.Name = name
	return e
}

//WithError attaches a key-value pair additional information into the Error instance
func (e *Error) WithError(key string, val interface{}) *Error {
	return e.WithFields(Fields{key: val})
}

//Stack records the current stack trace into the Error object
func (e *Error) Stack() string {
	return fmt.Sprintf("%+v", e.stack)
}

//Stack returns the current stack trace information
func Stack(skip ...int) string {
	theSkip := DefaultStackSkip
	if len(skip) > 0 {
		theSkip += skip[0]
	}
	return fmt.Sprintf("%+v", callers(theSkip))
}
