package fail

import (
	"fmt"
	"strings"
)

// Error slice
type List []error

func (s *List) Append(errs ...error) *List {
	*s = append(*s, errs...)
	return s
}

func (s List) Error() string {
	msg := []string{}

	for _, err := range s {
		if IsError(err) {
			msg = append(msg, err.Error())
		}
	}

	return "[" + strings.Join(msg, `,`) + "]"
}

func (s List) ErrorView() interface{} {
	view := []interface{}{}

	for _, err := range s {
		if IsError(err) {
			view = append(view, View(err))
		}
	}

	return view
}

func (s List) IsError() bool {
	for _, err := range s {
		if IsError(err) {
			return true
		}
	}

	return false
}

// Like list but with keys preserved
type Collection map[int]error

func (col Collection) Error() string {
	msg := []string{}

	for key, err := range col {
		if IsError(err) {
			msg = append(msg, fmt.Sprintf("%d: %s", key, err))
		}
	}

	return "[" + strings.Join(msg, `,`) + "]"
}

func (col Collection) ErrorView() interface{} {
	view := map[int]interface{}{}

	for key, err := range col {
		if IsError(err) {
			view[key] = View(err)
		}
	}

	return view
}

func (col Collection) IsError() bool {
	for _, err := range col {
		if IsError(err) {
			return true
		}
	}

	return false
}

// Error map
type Map map[string]error

func (m Map) Error() string {
	msg := []string{}

	for key, err := range m {
		if IsError(err) {
			msg = append(msg, fmt.Sprintf("%s: %s", key, err))
		}
	}

	return "[" + strings.Join(msg, `,`) + "]"
}

func (m Map) ErrorView() interface{} {
	view := map[string]interface{}{}

	for key, err := range m {
		if IsError(err) {
			view[key] = View(err)
		}
	}

	return view
}

func (m Map) IsError() bool {
	for _, err := range m {
		if IsError(err) {
			return true
		}
	}

	return false
}

// Applied to compound errors (maps, collection) where there are no errors inside
type Emptiable interface {
	IsError() bool
}

// Tests if error is actually an error, if it's not return nil
func OrNil(err error) error {
	if !IsError(err) {
		return nil
	}

	return err
}

// Tests if error is not a nil, or is an error actually  (for compound types)
func IsError(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(Emptiable); ok {
		return e.IsError()
	}

	return true
}

// Expected to return error data that can be decoded (by json.Decoder for example)
type Viewable interface {
	ErrorView() interface{}
}

// Get error view data
func View(err error) interface{} {

	if !IsError(err) {
		return nil
	}

	if e, ok := err.(Viewable); ok {
		return e.ErrorView()
	}

	return err.Error()
}

//simple replacement for errors.errorString
type Message string

func (m Message) Error() string {
	return string(m)
}

func Errorf(format string, args ...interface{}) Message {
	return Message(fmt.Sprintf(format, args...))
}
