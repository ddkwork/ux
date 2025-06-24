package assert

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/ddkwork/golibrary/std/assert/common"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/google/go-cmp/cmp"
)

//
// func Int(t common.T, x, y int) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Int8(t common.T, x, y int8) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Int16(t common.T, x, y int16) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Int32(t common.T, x, y int32) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Int64(t common.T, x, y int64) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Uint(t common.T, x, y uint) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Uint8(t common.T, x, y uint8) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Uint16(t common.T, x, y uint16) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Uint32(t common.T, x, y uint32) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Uint64[T uint64](t common.T, x, y T) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }

// func Uintptr(t common.T, x, y uintptr) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %d, got %d", x, y))
// }
//
// func Bool(t common.T, x, y bool) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %t, got %t", x, y))
// }
//
// func Float32(t common.T, x, y float32) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %f, got %f", x, y))
// }
//
// func Float64(t common.T, x, y float64) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %f, got %f", x, y))
// }
//
// func String(t common.T, x, y string) {
//	t.Helper()
//	if x == y {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %s, got %s", x, y))
// }
//
// func Slice(t common.T, x, y []any) {
//	t.Helper()
//	if len(x) != len(y) {
//		t.Error(fmt.Sprintf("expected slices of equal length, got %d and %d", len(x), len(y)))
//		return
//	}
//	for i := range x {
//		if x[i] != y[i] {
//			t.Error(fmt.Sprintf("expected slices to be equal, but differ at index %d", i))
//			return
//		}
//	}
// }
//
// func Map(t common.T, x, y map[any]any) {
//	t.Helper()
//	if len(x) != len(y) {
//		t.Error(fmt.Sprintf("expected maps of equal length, got %d and %d", len(x), len(y)))
//		return
//	}
//	for k, v := range x {
//		if y[k] != v {
//			t.Error(fmt.Sprintf("expected maps to be equal, but differ at key %v", k))
//			return
//		}
//	}
// }

// func Signed[T Signed](t common.T, want, got T, opts ...cmp.Option) {
//	t.Helper()
//	if want == got {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }

func UnsignedInteger[T Unsigned](t common.T, want, got T, opts ...cmp.Option) {
	t.Helper()
	if want == got {
		return
	}
	t.Error(fmt.Sprintf("expected %#X, got %#X", want, got))
}

// func Integer[T Integer](t common.T, want, got T, opts ...cmp.Option) {
//	t.Helper()
//	if want == got {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }
//
// func Float[T Float](t common.T, want, got T, opts ...cmp.Option) {
//	t.Helper()
//	if want == got {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }
//
// func Complex[T Complex](t common.T, want, got T, opts ...cmp.Option) {
//	t.Helper()
//	if want == got {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }
//
// func Ordered[T Ordered](t common.T, want, got T, opts ...cmp.Option) {
//	t.Helper()
//	if want == got {
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }
//
// func Bytes[T []byte](t common.T, want T, got T, opts ...cmp.Option) {
//	t.Helper()
//	if bytes.Equal(want, got) { // todo 对发包和注册算法进行大量单元测试，diff着色，以便于快速定位diff的位置
//		return
//	}
//	t.Error(fmt.Sprintf("expected %v, got %v", want, got))
// }

// ////////////////////////////////////

func True(t common.T, x bool) {
	t.Helper()
	if x {
		return
	}
	t.Error("expected true")
}

func False(t common.T, x bool) {
	t.Helper()
	if !x {
		return
	}
	t.Error("expected false")
}

func Equal[T any](t common.T, want, got T, opts ...cmp.Option) {
	t.Helper()
	if cmp.Equal(want, got, opts...) {
		return
	}
	t.Error(fmt.Sprintf("expected want == got\n--- want\n+++ got\n%s", cmp.Diff(want, got, opts...)))
}

func NotEqual[T any](t common.T, want, got T, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(want, got, opts...) {
		return
	}
	t.Error(fmt.Sprintf("expected want != got\nwant: %+v\n got: %+v", want, got))
}

func LessThan[T Ordered](t common.T, small, big T) {
	t.Helper()
	if small < big {
		return
	}
	t.Error(fmt.Sprintf("expected %v < %v", small, big))
}

func LessThanOrEqual[T Ordered](t common.T, small, big T) {
	t.Helper()
	if small <= big {
		return
	}
	t.Error(fmt.Sprintf("expected %v <= %v", small, big))
}

func GreaterThan[T Ordered](t common.T, big, small T) {
	t.Helper()
	if big > small {
		return
	}
	t.Error(fmt.Sprintf("expected %v > %v", big, small))
}

func GreaterThanOrEqual[T Ordered](t common.T, big, small T) {
	t.Helper()
	if big >= small {
		return
	}
	t.Error(fmt.Sprintf("expected %v >= %v", big, small))
}

func Error(t common.T, e error) {
	t.Helper()
	if e != nil {
		return
	}
	t.Error("expected error, received <nil>")
	mylog.Check(e)

	// t.Error("expected error, received <nil>")
}

func NoError(t common.T, e error) {
	t.Helper()
	if e == nil {
		return
	}
	mylog.Check(e)

	// if e == nil {
	//	return
	// }
	// t.Error(fmt.Sprintf("expected no error, received %v", e))
}

func Nil(t common.T, v any) {
	t.Helper()
	if isNil(v) {
		return
	}
	t.Error(fmt.Sprintf("expected <nil> error, received %v", v))
}

func NotNil(t common.T, v any) {
	t.Helper()
	if !isNil(v) {
		return
	}
	t.Error(fmt.Sprintf("expected <nil> error, received %v", v))
}

func isNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

func In[T any](t common.T, element T, slice []T, opts ...cmp.Option) {
	t.Helper()
	for _, value := range slice {
		if cmp.Equal(element, value, opts...) {
			return
		}
	}
	t.Error(fmt.Sprintf("expected slice to contain element:\nelement: %+v\n", element))
}

func NotIn[T any](t common.T, element T, slice []T, opts ...cmp.Option) {
	t.Helper()
	for _, value := range slice {
		if cmp.Equal(element, value, opts...) {
			t.Error(fmt.Sprintf("expected slice to not contain element\nelement: %+v\n  found: %+v", element, value))
			return
		}
	}
}

func ContainsString(t common.T, s, v string) {
	t.Helper()
	if !strings.Contains(s, v) {
		t.Error("expected to contain string, but does not")
	}
}

func Contains[S ~[]E, E comparable](t common.T, s S, v E) {
	t.Helper()
	if !slices.Contains(s, v) {
		t.Error("expected to contain element, but does not\n")
	}
}

func NotContainsString(t common.T, s, v string) { // todo remove ? use true api ?
	t.Helper()
	if strings.Contains(s, v) {
		t.Error("expected to not contain string, but does")
	}
}

func NotContains[S ~[]E, E comparable](t common.T, s S, v E) {
	t.Helper()
	if slices.Contains(s, v) {
		t.Error("expected to not contain element, but does")
	}
}

func isEmpty(object any) bool {
	// get nil case out of the way
	if object == nil {
		return true
	}
	objValue := reflect.ValueOf(object)
	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func Empty(t common.T, object any) {
	t.Helper()
	if !isEmpty(object) {
		t.Error(fmt.Sprintf("Should be empty, but was %v", object))
	}
}

func NotEmpty(t common.T, object any) {
	t.Helper()
	if isEmpty(object) {
		t.Error(fmt.Sprintf("Should be empty, but was %v", object))
	}
}

func getLen(x any) (length int, ok bool) {
	v := reflect.ValueOf(x)
	defer func() {
		ok = recover() == nil
	}()
	return v.Len(), true
}

func Len(t common.T, object any, length int) {
	t.Helper()
	l, ok := getLen(object)
	if !ok {
		t.Error(fmt.Sprintf("\"%v\" could not be applied builtin len()", object))
	}
	if l != length {
		t.Error(fmt.Sprintf("\"%v\" should have %d item(s), but has %d", object, length, l))
	}
}

func isZero(data any) bool {
	v := reflect.Indirect(reflect.ValueOf(data))
	if v.IsValid() {
		return true
	}
	return v.IsZero()
}

func Zero(t common.T, v any) {
	t.Helper()
	if !isZero(v) {
		t.Error(fmt.Sprintf("Should be zero, but was %v", v))
	}
}

func NotZero(t common.T, v any) {
	t.Helper()
	if isZero(v) {
		t.Error(fmt.Sprintf("Should not be zero, but was %v", v))
	}
}

/*

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (didPanic bool, message any, stack string) {
	didPanic = true

	defer func() {
		message = recover()
		if didPanic {
			stack = string(debug.Stack())
		}
	}()

	// call the target function
	f()
	didPanic = false

	return
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//	assert.Panics(t, func(){ GoCrazy() })
func Panics(t TestingT, f PanicTestFunc, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if funcDidPanic, panicValue, _ := didPanic(f); !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}

	return true
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//	assert.PanicsWithValue(t, "crazy error", func(){ GoCrazy() })
func PanicsWithValue(t TestingT, expected any, f PanicTestFunc, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	funcDidPanic, panicValue, panickedStack := didPanic(f)
	if !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}
	if panicValue != expected {
		return Fail(t, fmt.Sprintf("func %#v should panic with value:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s", f, expected, panicValue, panickedStack), msgAndArgs...)
	}

	return true
}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//	assert.PanicsWithError(t, "crazy error", func(){ GoCrazy() })
func PanicsWithError(t TestingT, errString string, f PanicTestFunc, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	funcDidPanic, panicValue, panickedStack := didPanic(f)
	if !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}
	panicErr, ok := panicValue.(error)
	if !ok || panicErr.Error() != errString {
		return Fail(t, fmt.Sprintf("func %#v should panic with error message:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s", f, errString, panicValue, panickedStack), msgAndArgs...)
	}

	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//	assert.NotPanics(t, func(){ RemainCalm() })
func NotPanics(t TestingT, f PanicTestFunc, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if funcDidPanic, panicValue, panickedStack := didPanic(f); funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should not panic\n\tPanic value:\t%v\n\tPanic stack:\t%s", f, panicValue, panickedStack), msgAndArgs...)
	}

	return true
}
*/

func Panics(t common.T, f func()) {
	t.Helper()
	mylog.Check(f)
}

func NotPanics(t common.T, expected any) {
	t.Helper()
	mylog.Check(expected)
}

/*
func Equalf(t TestingT, expected any, actual any, msg string, args ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Equal(t, expected, actual, append([]any{msg}, args...)...)
}
*/

func Equalf(t common.T, expected any, actual any, msg string, args ...any) {
	t.Helper()
	if cmp.Equal(expected, actual) {
		return
	}
	t.Error(fmt.Sprintf(msg, args...))
}
