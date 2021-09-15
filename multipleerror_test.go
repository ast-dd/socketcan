package socketcan_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ast-dd/socketcan"
)

func TestMultipleError(t *testing.T) {
	err1 := fmt.Errorf("e1")
	err2 := fmt.Errorf("e2")
	err3 := fmt.Errorf("e3")

	tests := []struct {
		name    string
		errs    []error
		wantErr error
		wantS   string
	}{
		// TODO: Add test cases.
		{"empty", nil, nil, ""},
		{"1 error", []error{err1}, err1, "e1"},
		{"3 errors", []error{err1, err2, err3}, &socketcan.MultipleError{}, "multiple errors occurred: e1, e2, e3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &socketcan.MultipleError{}
			for _, err := range tt.errs {
				me.Add(err)
			}

			// check type returned by Err()
			if got, want := me.Err(), tt.wantErr; reflect.TypeOf(got) != reflect.TypeOf(want) {
				t.Errorf("MultipleError.Err() type got %T, want %T", got, want)
			}

			if got, want := me.Error(), tt.wantS; !reflect.DeepEqual(got, want) {
				t.Errorf("MultipleError.Error() got %+#v, want %+#v", got, want)
			}

			if got, want := me.Errors(), tt.errs; !reflect.DeepEqual(got, want) {
				t.Errorf("MultipleError.Errors() got %+#v, want %+#v", got, want)
			}
		})
	}
}
