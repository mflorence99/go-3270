package utils

import (
	"reflect"
)

// 🟦 Recursively flattens any slice elements into a single-level []any

func Flatten(input any, strconv func(string) string) []any {
	var out []any
	v := reflect.ValueOf(input)

	switch v.Kind() {

	case reflect.String:
		str := v.String()
		if strconv != nil {
			str = strconv(str)
		}
		out = append(out, ForceBytes2Any([]byte(str))...)

	case reflect.Slice,
		reflect.Array:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			out = append(out, Flatten(item, strconv)...)
		}

	default:
		out = append(out, input)

	}
	return out
}

// 🟦 ... or directly into a []byte slice

func Flatten2Bytes(input any, strconv func(string) string) []byte {
	return ForceAny2Bytes(Flatten(input, strconv))
}
