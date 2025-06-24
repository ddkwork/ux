package deepcopy

import (
	"reflect"
	"time"
)

type Interface interface {
	Clone() any
}

func Iface(iface any) any {
	return Clone(iface)
}

// Clone tree node clone and tls ca clone
func Clone[T any](src T) T {
	if any(src) == nil {
		var zero T
		return zero
	}
	original := reflect.ValueOf(src)
	cpy := reflect.New(original.Type()).Elem()
	copyRecursive(original, cpy)
	return cpy.Interface().(T)
}

func copyRecursive(original, cpy reflect.Value) {
	if original.CanInterface() {
		if copier, ok := original.Interface().(Interface); ok {
			cpy.Set(reflect.ValueOf(copier.Clone()))
			return
		}
	}

	switch original.Kind() {
	case reflect.Ptr:

		originalValue := original.Elem()

		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, cpy.Elem())

	case reflect.Interface:

		if original.IsNil() {
			return
		}

		originalValue := original.Elem()

		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}

		for i := range original.NumField() {
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			copyRecursive(original.Field(i), cpy.Field(i))
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}

		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := range original.Len() {
			copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copyKey := Clone(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		cpy.Set(original)
	}
}
