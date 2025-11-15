package utils

import (
	"reflect"
)

// 🟦 Forced conversions

func ForceAny2Bytes(items []any) []byte {
	b := make([]byte, len(items))
	for i, v := range items {
		rv := reflect.ValueOf(v)

		switch rv.Kind() {

		case reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Int:
			b[i] = byte(rv.Convert(reflect.TypeOf(uint8(0))).Uint())

		case reflect.Float64,
			reflect.Float32:
			b[i] = byte(rv.Float())

		default:
			b[i] = 0

		}
	}
	return b
}

func ForceBytes2Any(items []byte) []any {
	v := reflect.ValueOf(items)
	a := make([]any, v.Len())
	for ix := 0; ix < v.Len(); ix++ {
		a[ix] = v.Index(ix).Interface()
	}
	return a
}
