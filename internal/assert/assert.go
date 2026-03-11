package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

func Null(t *testing.T, actual any) {
	t.Helper()

	fmt.Println("actual is: ", actual)
	if !DeepIsNil(actual) {
		t.Errorf("got: %v; expected: %v", actual, nil)
	}
}

func NotNull(t *testing.T, actual any) {
	t.Helper()

	if actual == nil {
		t.Errorf("expected value cannot be nil")
	}
}

func DeepIsNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Pointer,
		reflect.Map,
		reflect.Slice,
		reflect.Chan,
		reflect.Func:
		return v.IsNil()
	}

	return false
}
