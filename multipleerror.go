package socketcan

import (
	"strings"
)

// MultipleError implements the error interface and may contain multiple errors
type MultipleError struct {
	errors []error
}

// Add adds the given error to the contained errors
func (e *MultipleError) Add(err error) {
	if e.errors == nil {
		e.errors = make([]error, 0)
	}
	e.errors = append(e.errors, err)
}

// Err returns, dependent upon the count of contained errors, nil, the one
// contained error or the MultipleError itself.
// This should be returned as error to the calling function.
func (e *MultipleError) Err() (err error) {
	switch len(e.errors) {
	case 0:
		err = nil
	case 1:
		err = e.errors[0]
	default:
		err = e
	}
	return
}

// Errors returns the contained errors
func (e *MultipleError) Errors() (errs []error) {
	errs = e.errors
	return
}

// Error implements the error interface
func (e *MultipleError) Error() (s string) {
	switch len(e.errors) {
	case 0:
		s = ""
	case 1:
		s = e.errors[0].Error()
	default:
		ss := make([]string, 0)
		for _, err := range e.errors {
			ss = append(ss, err.Error())
		}
		s = "multiple errors occurred: " + strings.Join(ss, ", ")
	}

	return
}
