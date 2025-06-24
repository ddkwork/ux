package mylog

import (
	"errors"
	"io"
	"reflect"
	"slices"
)

func Check[T any](result T) (isEof bool) { // check these return type: bool,custom error and error.
	switch r := any(result).(type) {
	case bool:
		if !r {
			panic("return false detected")
		}
	case string:
		if !successfully(r) {
			panic(r)
		}
	case error:
		if errors.Is(r, io.EOF) { // should break in loop
			return true
		}
		if !successfully(r.Error()) {
			panic(r.Error())
		}
	}
	return
}

func successfully(err string) bool {
	success := []string{"The operation completed successfully.", "STATUS_SUCCESS"} // for windows api
	return slices.Contains(success, err)
}

func CheckNil(ptr any) {
	isNil(ptr)
	switch r := ptr.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if r == 0 {
		}
		if r == -1 {
			panic("return -1")
		}
	case string:
		if r == "" { // todo remove
		}
	case []byte:
		if r == nil {
			panic("detected nil []byte")
		}
		if len(r) == 0 { // todo remove
		}
	}
}

func isNil(i any) {
	if i == nil {
		return
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if reflect.ValueOf(i).IsNil() {
			panic("detected nil pointer " + reflect.TypeOf(i).String())
		}
	}
}

func Check2ForJsonNumberNodeType[T any](ret T, err error) bool {
	return err == nil
}

func Check2[T any](ret T, err error) T {
	Check(err)
	CheckNil(ret)
	return ret
}

func Check2Ignore[T any](ret T, err error) T {
	l.errorCall(err)
	return ret
}

func Check2IgnoreBool[T any](ret T, err error) (T, bool) {
	return ret, l.errorCall(err)
}

func CheckIgnore(err any) {
	l.errorCall(err)
}

func Check2Bool[T any](ret T, ok bool) T {
	Check(ok)
	CheckNil(ret)
	return ret
}

func Check3[T1 any, T2 any](ret1 T1, ret2 T2, err error) (T1, T2) {
	Check(err)
	CheckNil(ret1)
	CheckNil(ret2)
	return ret1, ret2
}

func Check3Bool[T1 any, T2 any](ret1 T1, ret2 T2, ok bool) (T1, T2) {
	Check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	return ret1, ret2
}

func Check4[T1 any, T2 any, T3 any](ret1 T1, ret2 T2, ret3 T3, err error) (T1, T2, T3) {
	Check(err)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	return ret1, ret2, ret3
}

func Check4Bool[T1 any, T2 any, T3 any](ret1 T1, ret2 T2, ret3 T3, ok bool) (T1, T2, T3) {
	Check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	return ret1, ret2, ret3
}

func Check5[T1 any, T2 any, T3 any, T4 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, err error) (T1, T2, T3, T4) {
	Check(err)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	return ret1, ret2, ret3, ret4
}

func Check5Bool[T1 any, T2 any, T3 any, T4 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ok bool) (T1, T2, T3, T4) {
	Check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	return ret1, ret2, ret3, ret4
}

func Check6[T1 any, T2 any, T3 any, T4 any, T5 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, err error) (T1, T2, T3, T4, T5) {
	Check(err)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	return ret1, ret2, ret3, ret4, ret5
}

func Check6Bool[T1 any, T2 any, T3 any, T4 any, T5 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ok bool) (T1, T2, T3, T4, T5) {
	Check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	return ret1, ret2, ret3, ret4, ret5
}

func Check7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ret6 T6, err error) (T1, T2, T3, T4, T5, T6) {
	Check(err)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	CheckNil(ret6)
	return ret1, ret2, ret3, ret4, ret5, ret6
}

func Check7Bool[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ret6 T6, ok bool) (T1, T2, T3, T4, T5, T6) {
	Check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	CheckNil(ret6)
	return ret1, ret2, ret3, ret4, ret5, ret6
}
