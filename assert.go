// Assert is a tiny assertion library for Go, the kind that wouldn't be
// out of place in stdlib.
//
// All assertion methods in this package accept optional extra fmt.Sprintf-style
// arguments. If passed in, the first such argument must be a format string.
//
// All assertion methods return true when succeeded and false when failed
// to drive subsequent comparisons:
package assert

import (
	"errors"
	"fmt"
	"reflect"
)

// TB contains the parts of testing.TB that this package actually needs. Pass *testing.T or *testing.B for arguments of type TB.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
}

// OK asserts that the value is true.
func OK(t TB, a bool, messageAndArgs ...any) bool {
	if !a {
		t.Helper()
		t.Errorf("** %sgot false, wanted true", FormatPrefix(messageAndArgs))
		return false
	}
	return true
}

// False asserts that the value is false.
func False(t TB, a bool, messageAndArgs ...any) bool {
	if a {
		t.Helper()
		t.Errorf("** %sgot true, wanted false", FormatPrefix(messageAndArgs))
		return false
	}
	return true
}

// Eq asserts that two values are equal via == operator.
func Eq[T comparable](t TB, a, e T, messageAndArgs ...any) bool {
	if a != e {
		t.Helper()
		t.Errorf("** %sgot %v, wanted %v", FormatPrefix(messageAndArgs), a, e)
		return false
	}
	return true
}

// NotEq asserts that two values are not equal via != operator.
func NotEq[T comparable](t TB, a, e T, messageAndArgs ...any) bool {
	if a == e {
		t.Helper()
		t.Errorf("** %sgot %v, wanted anything else", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// DeepEqual asserts that two values are equal via reflect.DeepEqual.
func DeepEqual[T any](t TB, a, e T, messageAndArgs ...any) bool {
	if !reflect.DeepEqual(a, e) {
		t.Helper()
		t.Errorf("** %sgot %v, wanted %v", FormatPrefix(messageAndArgs), a, e)
		return false
	}
	return true
}

// NotDeepEqual asserts that two values are not equal via !reflect.DeepEqual.
//
// Note: I never had a use for this, but including it for completeness.
func NotDeepEqual[T any](t TB, a, e T, messageAndArgs ...any) bool {
	if reflect.DeepEqual(a, e) {
		t.Helper()
		t.Errorf("** %sgot %v, wanted anything else", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// MethodEqual asserts that two values are equal via their Equal method (like time.Time).
func MethodEqual[T interface{ Equal(T) bool }](t TB, a, e T, messageAndArgs ...any) bool {
	t.Helper()
	if !e.Equal(a) {
		t.Errorf("** %sgot %v, wanted %v", FormatPrefix(messageAndArgs), a, e)
		return false
	}
	return true
}

// NotMethodEqual asserts that two values are not equal via their Equal method (like time.Time).
func NotMethodEqual[T interface{ Equal(T) bool }](t TB, a, e T, messageAndArgs ...any) bool {
	t.Helper()
	if e.Equal(a) {
		t.Errorf("** %sgot %v, wanted anything else", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// Nil asserts that a pointer value is nil.
//
// Use Zero for interface values. Nil is
func Nil[T any, P ~*T](t TB, a P, messageAndArgs ...any) bool {
	if a != nil {
		t.Helper()
		t.Errorf("** %sgot &%v, wanted nil", FormatPrefix(messageAndArgs), *a)
		return false
	}
	return true
}

// NonNil asserts that a pointer value is anything but nil.
func NonNil[T any](t TB, a *T, messageAndArgs ...any) bool {
	if a == nil {
		t.Helper()
		t.Errorf("** %sgot nil %T, wanted non-nil", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// Zero asserts that the value is equal to a zero value for its type using == operator.
func Zero[T comparable](t TB, a T, messageAndArgs ...any) bool {
	var zero T
	if a != zero {
		t.Helper()
		t.Errorf("** %sgot %v, wanted zero value %v", FormatPrefix(messageAndArgs), a, zero)
		return false
	}
	return true
}

// NonZero asserts that the value is not equal to a zero value for its type using != operator.
func NonZero[T comparable](t TB, a T, messageAndArgs ...any) bool {
	var zero T
	if a == zero {
		t.Helper()
		t.Errorf("** %sgot zero value %v, wanted non-zero", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// EmptySlice asserts that the given slice is nil or empty.
func EmptySlice[T any, S ~[]T](t TB, a S, messageAndArgs ...any) bool {
	if len(a) > 0 {
		t.Helper()
		t.Errorf("** %sgot %v, wanted empty slice", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// NonEmptySlice asserts that the given slice has non-zero length.
func NonEmptySlice[T any, S ~[]T](t TB, a S, messageAndArgs ...any) bool {
	if len(a) == 0 {
		t.Helper()
		t.Errorf("** %sgot empty %T, wanted non-empty", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// EmptyMap asserts that the given map is nil or empty.
func EmptyMap[K comparable, V any, M ~map[K]V](t TB, a M, messageAndArgs ...any) bool {
	if len(a) > 0 {
		t.Helper()
		t.Errorf("** %sgot %v, wanted empty map", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// NonEmptyMap asserts that the given map has non-zero length.
func NonEmptyMap[K comparable, V any, M ~map[K]V](t TB, a M, messageAndArgs ...any) bool {
	if len(a) == 0 {
		t.Helper()
		t.Errorf("** %sgot empty %T, wanted non-empty", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// Success asserts that the error is nil.
func Success(t TB, a error, messageAndArgs ...any) bool {
	if a != nil {
		t.Helper()
		t.Errorf("** %sfailed: %v", FormatPrefix(messageAndArgs), a)
		return false
	}
	return true
}

// Error asserts that the actual error value is equivalent to the expected
// error value using errors.Is.
//
// If the expected error is nil, behaves exactly like Success.
func Error(t TB, a, e error, messageAndArgs ...any) bool {
	if e == nil {
		t.Helper()
		return Success(t, a, messageAndArgs...)
	} else if a == nil {
		t.Helper()
		t.Errorf("** %ssucceeded, wanted to fail with: %v", FormatPrefix(messageAndArgs), e)
		return false
	} else if !errors.Is(a, e) {
		t.Helper()
		t.Errorf("** %sfailed with: %v, wanted: %v", FormatPrefix(messageAndArgs), a, e)
		return false
	}
	return true
}

// Error asserts that the actual error message is equivalent to the expected one.
//
// If the expected error message is empty, behaves exactly like Success.
func ErrorMsg(t TB, a error, e string, messageAndArgs ...any) bool {
	if e == "" {
		return Success(t, a, messageAndArgs...)
	} else if a == nil {
		t.Helper()
		t.Errorf("** %ssucceeded, wanted to fail with: %v", FormatPrefix(messageAndArgs), e)
		return false
	} else if s := a.Error(); s != e {
		t.Helper()
		t.Errorf("** %sfailed with: %v, wanted: %v", FormatPrefix(messageAndArgs), s, e)
		return false
	}
	return true
}

// PanicMsg asserts that a function panics with the given message.
func PanicMsg(t TB, f func(), e string, messageAndArgs ...any) bool {
	actual := capturePanic(f)
	if actual == nil {
		t.Helper()
		t.Errorf("** %ssucceeded, wanted to panic with: %v", FormatPrefix(messageAndArgs), e)
		return false
	} else if a := fmt.Sprint(actual); a != e {
		t.Helper()
		t.Errorf("** %spaniced with: %v, wanted: %v", FormatPrefix(messageAndArgs), a, e)
		return false
	}
	return true
}

func capturePanic(f func()) (panicValue any) {
	defer func() {
		panicValue = recover()
	}()
	f()
	return
}

// FormatPrefix returns a prefix for assertion error messages based on messageAndArgs arguments.
// If messageAndArgs is empty, returns an empty string.
func FormatPrefix(messageAndArgs []any) string {
	if len(messageAndArgs) == 0 {
		return ""
	}
	if msg, ok := messageAndArgs[0].(string); ok {
		return fmt.Sprintf(msg, messageAndArgs[1:]...) + ": "
	} else {
		panic(fmt.Errorf("when passing messageAndArgs to assertion funcs, the first extra argument must be a format string, got %T %v", messageAndArgs[0], messageAndArgs[0]))
	}
}

// AddPrefix adds a prefix to assertion call's messageAndArgs argument.
func AddPrefix(messageAndArgs []any, prefix string, args ...any) []any {
	if len(messageAndArgs) > 0 {
		if msg, ok := messageAndArgs[0].(string); ok {
			prefix = prefix + ": " + msg
			messageAndArgs = messageAndArgs[1:]
		}
	}
	return append(append([]any{prefix}, args...), messageAndArgs...)
}
