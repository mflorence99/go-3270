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
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	// 👇  call and collect results
	outValues := v.Call(in)
	out := make([]any, len(outValues))
	for i, val := range outValues {
		out[i] = val.Interface()
	}
	return out
}
