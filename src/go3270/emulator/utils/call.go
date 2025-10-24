package utils

import (
	"reflect"
)

func Call(fn any, args ...any) []any {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic("Call: argument is not a function")
	}
	// 👇 convert args to []reflect.Value
	in := make([]reflect.Value, len(args))
	for ix, arg := range args {
		in[ix] = reflect.ValueOf(arg)
	}
	// 👇  call and collect results
	outValues := v.Call(in)
	out := make([]any, len(outValues))
	for ix, val := range outValues {
		out[ix] = val.Interface()
	}
	return out
}
