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
